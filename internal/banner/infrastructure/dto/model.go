package dto

import (
	"banners/internal/banner/domain/entity"
	"fmt"
	"time"
)

type UserBanner struct {
	Title string `json:"title" db:"title"`
	Text  string `json:"text" db:"text"`
	Url   string `json:"url" db:"url"`
}

type InvalidDataError struct {
	Line int
	Col  int
}

func (e *InvalidDataError) Error() string {
	return fmt.Sprintf("%d:%d: invalid data error", e.Line, e.Col)
}

type CreateBanner struct {
	Tags_ids   []int          `json:"tags_ids" db:"tags" validate:"required,dive,min=0"`
	Feature_id int            `json:"feature_id" db:"feature_id" validate:"required,min=0"`
	Content    entity.Content `json:"content" db:"content" validate:"required,dive,min=0"`
	Is_active  bool           `json:"is_active" db:"is_active" validate:"required"`
	Created_at time.Time      `json:"created_at" db:"created_at" validate:"required"`
	Updated_at time.Time      `json:"updated_at" db:"updated_at" validate:"required"`
}

type GetBannerById struct {
	Banner_id  int            `json:"banner_id" db:"id"`
	Tags_ids   []int          `json:"tags_ids" db:"tags"`
	Feature_id int            `json:"feature_id" db:"feature_id"`
	Content    entity.Content `json:"content" db:"content"`
	Is_active  bool           `json:"is_active" db:"is_active"`
	Created_at time.Time      `json:"created_at" db:"created_at"`
	Updated_at time.Time      `json:"updated_at" db:"updated_at"`
}

type UpdateBanner struct {
	Tags_ids   []int          `json:"tags_ids" db:"tags" validate:"omitempty,dive,numeric,min=0"`
	Feature_id int            `json:"feature_id" db:"feature_id" validate:"omitempty,numeric,min=0"`
	Content    entity.Content `json:"content" db:"content" validate:"omitempty,dive,min=0"`
	Is_active  bool           `json:"is_active" db:"is_active" validate:"omitempty"`
	Created_at time.Time      `json:"created_at" db:"created_at" validate:"omitempty"`
	Updated_at time.Time      `json:"updated_at" db:"updated_at" validate:"omitempty"`
}

type ListBanners []GetBannerById

type GetListBanners struct {
	ListBanners
}
