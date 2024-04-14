package authrepository

import (
	"banners/internal/auth/domain/entity"
	"banners/internal/auth/infrastructure/dto"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(user dto.User) (int64, error)
	IsUserExists(user dto.User) (bool, error)
	GetUserByEmail(user dto.User) (entity.User, error)
}

type AuthRepo struct {
	db *sqlx.DB
}

func (a *AuthRepo) CreateUser(user dto.User) (int64, error) {
	var id int64
	row := a.db.QueryRow("INSERT INTO users (username, email, password, is_admin) VALUES ($1, $2, $3, $4) RETURNING id", user.Username, user.Email, user.Password, false)
	err := row.Scan(&id)
	if err != nil {
		log.Printf("error while insert values: %s", err)
		return id, err
	}
	return id, nil
}

func (a *AuthRepo) IsUserExists(user dto.User) (bool, error) {
	var id int
	err := a.db.Get(&id, "SELECT id FROM users WHERE email = $1", user.Email)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		log.Printf("error while querying: %s", err)
		return false, err
	} else {
		return true, nil
	}
}

func (a *AuthRepo) GetUserByEmail(user dto.User) (entity.User, error) {
	var db_user entity.User

	err := a.db.Get(&db_user, "SELECT username, email, password, is_admin FROM users WHERE email = $1", user.Email)
	if err == sql.ErrNoRows {
		log.Println("No email in db!")
		return db_user, nil

	} else if err != nil {
		log.Printf("error while querying user: %s", err)
		return db_user, err
	} else {
		return db_user, nil
	}
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &AuthRepo{db: db}
}
