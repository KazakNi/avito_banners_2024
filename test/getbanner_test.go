package test

import (
	"banners/config"
	"banners/internal/auth/domain/entity"
	auth "banners/internal/auth/infrastructure/delivery/http"
	"banners/internal/auth/infrastructure/dto"
	authrepository "banners/internal/auth/infrastructure/repository"
	banentity "banners/internal/banner/domain/entity"
	banhttp "banners/internal/banner/infrastructure/delivery/http"
	banners "banners/internal/banner/infrastructure/delivery/http"
	bandto "banners/internal/banner/infrastructure/dto"
	"banners/internal/banner/infrastructure/repository"
	"banners/pkg/cache"
	"banners/pkg/db"
	"banners/pkg/db/migrations"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/assert.v1"
)

func TearUpDB() (auth.AuthHandler, banners.BannerHandler, *sqlx.DB) {

	config.LoadConfig()
	migrations.LoadTestMigrations()
	dbConnection, _ := db.NewDBTestConnection()

	authRepo := authrepository.NewAuthRepository(dbConnection)
	authHandler := auth.AuthHandler{Store: authRepo}

	cache.Cache = *cache.LoadCache()

	bannerRepo := repository.NewBannerRepository(dbConnection, &cache.Cache)
	bannerHadler := banners.NewBannerHandler(bannerRepo)

	return authHandler, *bannerHadler, dbConnection

}

func CreateAdminUser() entity.User {
	var user entity.User
	user.Username = "tes9t"
	user.Email = "te8st@test.ru"
	user.Password = "lolkek"
	user.AdminRights = true
	return user
}

func CreateBanner() bandto.CreateBanner {
	return bandto.CreateBanner{
		Tags_ids:   []int64{1, 2, 3},
		Feature_id: 1,
		Content: banentity.Content{
			Text:  "test",
			Title: "test",
			Url:   "www.avito.ru",
		},
		Is_active: true,
	}
}

func TestGetBanner(t *testing.T) {
	authHandler, bannerHandler, db := TearUpDB()
	user := CreateAdminUser()
	b, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/user/sign_up", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	authHandler.SignUp(w, req)

	res := w.Result()

	assert.Equal(t, res.StatusCode, 201)

	req = httptest.NewRequest(http.MethodPost, "/user/sign_in", bytes.NewBuffer(b))
	w = httptest.NewRecorder()

	authHandler.SignIn(w, req)

	res = w.Result()

	db.Exec("UPDATE users SET is_admin = true WHERE id = 0")

	var token dto.Token
	json.NewDecoder(res.Body).Decode(&token)
	assert.Equal(t, res.StatusCode, 200)

	banner := CreateBanner()
	b, _ = json.Marshal(banner)

	req = httptest.NewRequest(http.MethodPost, "/banner", bytes.NewBuffer(b))
	w = httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer "+token.BearerToken)

	banhttp.AuthRequiredCheck(banhttp.IsAdminCheck(http.HandlerFunc(bannerHandler.CreateBanner))).ServeHTTP(w, req)
	res = w.Result()
	assert.Equal(t, res.StatusCode, 201)

	req = httptest.NewRequest(http.MethodGet, "/user_banner?tag_id=2&feature_id=1", bytes.NewBuffer([]byte("")))
	w = httptest.NewRecorder()

	bannerHandler.GetBanner(w, req)

	res = w.Result()

	var content banentity.Content
	json.NewDecoder(res.Body).Decode(&content)
	assert.Equal(t, content.Text, "test")
	migrations.TestMigrationsDown()
}
