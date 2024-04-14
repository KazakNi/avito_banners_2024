package http

import (
	bannerEntity "banners/internal/banner/domain/entity"
	"banners/internal/banner/infrastructure/dto"
	"banners/internal/banner/infrastructure/repository"
	slogger "banners/pkg/logger"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	BannerRe       = regexp.MustCompile(`^/banner/*$`)
	BannerReWithID = regexp.MustCompile(`^/banner/\d+$`)
	UserBannerRe   = regexp.MustCompile(`^/user_banner/?`)
)

type BannerHandler struct {
	Store repository.BannerRepository
}

func (b *BannerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && BannerRe.MatchString(r.URL.Path):
		AuthRequiredCheck(IsAdminCheck(http.HandlerFunc(b.CreateBanner))).ServeHTTP(w, r)
		return

	case r.Method == http.MethodGet && BannerRe.MatchString(r.URL.Path):
		AuthRequiredCheck(IsAdminCheck(http.HandlerFunc(b.ListBanner))).ServeHTTP(w, r)
		return

	case r.Method == http.MethodGet && UserBannerRe.MatchString(r.URL.Path):
		AuthRequiredCheck(http.HandlerFunc(b.GetBanner)).ServeHTTP(w, r)
		return

	case r.Method == http.MethodPatch && BannerReWithID.MatchString(r.URL.Path):
		AuthRequiredCheck(IsAdminCheck(http.HandlerFunc(b.UpdateBanner))).ServeHTTP(w, r)
		return

	case r.Method == http.MethodDelete && BannerReWithID.MatchString(r.URL.Path):
		AuthRequiredCheck(IsAdminCheck(http.HandlerFunc(b.DeleteBanner))).ServeHTTP(w, r)
		return

	default:
		NotFoundHandler(w, r)
		return
	}
}

func (b *BannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	banner := &dto.CreateBanner{}
	if err := json.NewDecoder(r.Body).Decode(banner); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	if err := banner.Validate(); err != nil {
		BadRequestHandler(w, r)
		slogger.Logger.Info("error while banner validation: %s", err)
		return
	}

	id, err := b.Store.CreateBanner(*banner)
	if err != nil {
		slogger.Logger.Info("error while banner creation", "err", err)
		InternalServerErrorHandler(w, r)
		return
	}

	StatusCreatedBannerHandler(w, r, id)
}

func (b *BannerHandler) ListBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := r.URL.Query()

	banners, err := b.Store.GetBannersList(params)

	if err == repository.ErrInvalidParamType || err == repository.ErrNegativeId {
		BadRequestHandler(w, r)
		slogger.Logger.Info("error params validation while GetBannersList", "err", err)
		return
	}

	if err != nil {
		log.Println(err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	StatusOkListBanner(w, r, banners)

}

func (b *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	banner := &dto.UpdateBanner{}
	if err := json.NewDecoder(r.Body).Decode(banner); err != nil {
		slogger.Logger.Info("error while UpdateBanner validation: %s", err)
		InternalServerErrorHandler(w, r)
		return
	}
	if err := banner.Validate(); err != nil {
		BadRequestHandler(w, r)
		return
	}
	idParam := strings.TrimPrefix(r.URL.Path, "/banner/")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		BadRequestHandler(w, r)
		return
	}

	_, err = b.Store.GetBannerById(int64(id))

	if err == sql.ErrNoRows {
		NotFoundHandler(w, r)
		return
	}

	if err := b.Store.UpdateBanner(int64(id), *banner); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (b *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idParam := strings.TrimPrefix(r.URL.Path, "/banner/")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		BadRequestHandler(w, r)
		return
	}

	_, err = b.Store.GetBannerById(int64(id))

	if err == sql.ErrNoRows {
		NotFoundHandler(w, r)
		return
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := b.Store.DeleteBanner(int64(id)); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}

func (b *BannerHandler) GetBanner(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	params := r.URL.Query()

	content, active, err := b.Store.GetUserBanner(params)

	if !active {
		if requestUserIsAdmin := CheckIsAdmin(w, r); !requestUserIsAdmin {
			StatusForbidden(w, r)
			return
		}
		StatusOkContent(w, r, content)
		return
	}

	if err == repository.ErrInvalidParamType {
		BadRequestHandler(w, r)
		slogger.Logger.Info("error params validation while GetBanner", "err", err)
		return
	}

	if err == sql.ErrNoRows {
		NotFoundHandler(w, r)
		return
	}

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	StatusOkContent(w, r, content)
}

func NewBannerHandler(b repository.BannerRepository) *BannerHandler {
	return &BannerHandler{
		Store: b,
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "500 Internal Server Error"})
	w.Write([]byte(b))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "404 Not Found"})
	w.Write([]byte(b))
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "400 Bad request"})
	w.Write([]byte(b))
}

func StatusCreatedBannerHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, _ := json.Marshal(dto.CreatedBannerId{Banner_id: strconv.Itoa(int(id))})
	w.Write([]byte(b))
}

func StatusOkListBanner(w http.ResponseWriter, r *http.Request, banners []dto.GetBannerById) {
	data, _ := json.Marshal(banners)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func StatusOkContent(w http.ResponseWriter, r *http.Request, content bannerEntity.Content) {
	data, _ := json.Marshal(content)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func StatusForbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "403 Forbidden"})
	w.Write(b)
}
