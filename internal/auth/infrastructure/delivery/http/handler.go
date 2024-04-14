package auth

import (
	"banners/config"
	"banners/internal/auth/domain/entity"

	"banners/internal/auth/infrastructure/dto"
	authrepository "banners/internal/auth/infrastructure/repository"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	SignIn = regexp.MustCompile(`^/user/sign_in/*$`)
	SignUp = regexp.MustCompile(`^/user/sign_up/*$`)
)

type AuthHandler struct {
	Store authrepository.AuthRepository
}

func (a *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && SignUp.MatchString(r.URL.Path):
		a.SignUp(w, r)
		return

	case r.Method == http.MethodPost && SignIn.MatchString(r.URL.Path):
		a.SignIn(w, r)
		return
	}
}

func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	user := &dto.User{}
	err := json.NewDecoder(r.Body).Decode(user)

	if err != nil {
		BadRequestHandler(w, r)
		return
	}

	userExists, err := a.Store.IsUserExists(*user)

	if err != nil {
		log.Printf("Error while querying user: %s", err)
		InternalServerErrorHandler(w, r)
		return
	}

	if userExists {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("409 - User already exists"))
		return
	}

	user.HashPassword(user.Password)
	userId, err := a.Store.CreateUser(*user)

	if err != nil {
		log.Printf("Error while creating an user: %s", err)
		InternalServerErrorHandler(w, r)
		return
	}

	StatusIdUserHandler(w, r, userId)
}

func (a *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	user := &dto.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Printf("Error while %s endpoint response body parsing: %s", r.URL, err)
		InternalServerErrorHandler(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	db_user, err := a.Store.GetUserByEmail(*user)

	if err != nil {
		log.Printf("Error while getting user by email")
		InternalServerErrorHandler(w, r)
		return
	}

	if db_user.Email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid credentials"))
		log.Println("Wrong email")
		return
	}

	if user.CheckPassword(user.Password, db_user.Password) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid credentials"))
		log.Println("Wrong password")
		return
	}
	userID := strconv.Itoa(db_user.Id)

	claims := entity.CustomClaims{db_user.AdminRights, jwt.RegisteredClaims{ID: userID, ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour))}}

	SetToken(w, r, claims)

}

func SetToken(w http.ResponseWriter, r *http.Request, claims entity.CustomClaims) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_string, err := token.SignedString([]byte(config.Cfg.Token.Secret + config.Cfg.Token.Salt))

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid token"))
	}

	w.Header().Set("Content-Type", "application/json")

	t := make(map[string]string)
	t["token"] = token_string

	res, _ := json.Marshal(t)

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "400 Bad request"})
	w.Write([]byte(b))
}

func StatusIdUserHandler(w http.ResponseWriter, r *http.Request, id int64) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	b, _ := json.Marshal(dto.CreatedUserId{User_id: strconv.Itoa(int(id))})
	w.Write([]byte(b))
}

func UserTokenHandler(w http.ResponseWriter, r *http.Request, token string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(dto.Token{BearerToken: token})
	w.Write([]byte(b))
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	b, _ := json.Marshal(dto.ErrorResponse{Error: "500 Internal Server Error"})
	w.Write([]byte(b))
}
