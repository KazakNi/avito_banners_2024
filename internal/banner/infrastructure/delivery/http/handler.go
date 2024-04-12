package http

import (
	"banners/internal/banner/infrastructure/repository"
	"net/http"
	"regexp"
)

var (
	BannerRe       = regexp.MustCompile(`^/banner/*$`)
	BannerReWithID = regexp.MustCompile(`^/banner/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
	UserBannerRe   = regexp.MustCompile(`^/user_bannner/?`)
)

type BannerHandler struct {
	store repository.BannerRepository
}

// прописать методы хэндлера

func (h *BannerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && RecipeRe.MatchString(r.URL.Path):
		h.CreateRecipe(w, r)
		return
	case r.Method == http.MethodGet && RecipeRe.MatchString(r.URL.Path):
		h.ListRecipes(w, r)
		return
	case r.Method == http.MethodGet && RecipeReWithID.MatchString(r.URL.Path):
		h.GetRecipe(w, r)
		return
	case r.Method == http.MethodPut && RecipeReWithID.MatchString(r.URL.Path):
		h.UpdateRecipe(w, r)
		return
	case r.Method == http.MethodDelete && RecipeReWithID.MatchString(r.URL.Path):
		h.DeleteRecipe(w, r)
		return
	default:
		NotFoundHandler(w, r)
		return
	}
}

func NewRecipesHandler(b repository.BannerRepository) *BannerHandler {
	return &BannerHandler{
		store: b,
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("404 Not Found"))
}
