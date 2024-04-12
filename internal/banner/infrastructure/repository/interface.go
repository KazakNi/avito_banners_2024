package repository

import (
	"banners/internal/banner/domain/entity"
	"banners/internal/banner/infrastructure/dto"
	"net/url"
)

type BannerRepository interface {
	CreateBanner(banner dto.CreateBanner) error
	GetBannersList(params url.Values) (dto.GetListBanners, error)
	GetUserBanner(params url.Values) (entity.Content, error)
	UpdateBanner(id int64, banner dto.UpdateBanner) error
	DeleteBanner(id int64) error
}
