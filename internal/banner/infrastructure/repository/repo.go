package repository

import (
	"banners/internal/banner/infrastructure/dto"

	"github.com/jmoiron/sqlx"
)

type bannerRepo struct {
	db    *sqlx.DB
	cache interface{} // implement cache
}

func NewBannerRepository(db *sqlx.DB) BannerRepository {
	return &bannerRepo{db: db}
}

/*
CreateBanner(banner dto.CreateBanner) error
	GetBannersList(params url.Values) (dto.GetListBanners, error)
	GetUserBanner(params url.Values) (entity.Content, error)
	UpdateBanner(id int64, banner dto.UpdateBanner) error
	DeleteBanner(id int64) error
*/

func (b *bannerRepo) CreateBanner(banner dto.CreateBanner) error {

}
