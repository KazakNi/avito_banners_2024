package http

import (
	"banners/config"
	"banners/internal/auth/domain/entity"
	slogger "banners/pkg/logger"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func AuthRequiredCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")

		if len(AuthHeader) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Auth required!"))
			return
		}

		token := strings.Fields(AuthHeader)[1]

		err := ParseToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf("%s", err)))
			log.Println(err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ParseToken(mytoken string) error {

	token, err := jwt.ParseWithClaims(mytoken, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Token.Secret + config.Cfg.Token.Salt), nil
	})
	if err != nil {
		log.Println("Error during token parsing", err)
		return errors.New("invalid token")
	}
	if claims, ok := token.Claims.(*entity.CustomClaims); ok && token.Valid {
		if claims.RegisteredClaims.ExpiresAt.Unix() < time.Now().Unix() {
			return errors.New("token is expired")
		}
	} else {
		return errors.New("invalid token")
	}
	return nil

}

func IsAdminCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")

		header_token := strings.Fields(AuthHeader)[1]
		token, _ := jwt.ParseWithClaims(header_token, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.Token.Secret + config.Cfg.Token.Salt), nil
		})

		claims, _ := token.Claims.(*entity.CustomClaims)

		if !claims.AdminRights {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			log.Printf("User ID %s is not admin\n", claims.ID)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CheckIsAdmin(w http.ResponseWriter, r *http.Request) bool {
	AuthHeader := r.Header.Get("Authorization")

	header_token := strings.Fields(AuthHeader)[1]
	token, _ := jwt.ParseWithClaims(header_token, &entity.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.Token.Secret + config.Cfg.Token.Salt), nil
	})

	claims, _ := token.Claims.(*entity.CustomClaims)
	return claims.AdminRights
}

type ResponseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

func (rww ResponseWriterWrapper) String() string {
	var buf bytes.Buffer

	buf.WriteString("Response:")

	buf.WriteString("Headers:")
	for k, v := range (*rww.w).Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	buf.WriteString(fmt.Sprintf(" Status Code: %d", *(rww.statusCode)))

	buf.WriteString("Body")
	buf.WriteString(rww.body.String())
	return buf.String()
}
func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()

}

func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}
func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode int = 200
	return ResponseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slogger.Logger.Info("income request", "endpoint", r.URL, "method", r.Method, "rBody", r.Body)
		defer func() {
			rww := NewResponseWriterWrapper(w)
			slogger.Logger.Info("Response data", "Request", r, "Response", rww.String())
		}()
		next.ServeHTTP(w, r)

	})
}
