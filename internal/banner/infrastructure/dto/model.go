package dto

import (
	"banners/internal/banner/domain/entity"
	"banners/pkg/db"
	"time"

	"github.com/go-playground/validator"
	"github.com/lib/pq"
)

type UserBanner struct {
	Title string `json:"title" db:"title"`
	Text  string `json:"text" db:"text"`
	Url   string `json:"url" db:"url"`
}

type CreateBanner struct {
	Tags_ids   []int64        `json:"tag_ids" db:"tags" validate:"required,dive,min=0"`
	Feature_id int64          `json:"feature_id" db:"feature_id" validate:"required,min=0"`
	Content    entity.Content `json:"content" db:"content" validate:"required,dive,min=0"`
	Is_active  bool           `json:"is_active" db:"is_active" validate:"required"`
}

func (c *CreateBanner) Validate() error {
	validate := validator.New()
	err := validate.Struct(c)

	if err != nil {
		return err
	}
	return nil
}

func (c *CreateBanner) DtoToEntity() entity.Banner {

	return entity.Banner{
		Tags_ids:   c.Tags_ids,
		Feature_id: c.Feature_id,
		Content:    c.Content,
		Is_active:  c.Is_active,
		Created_at: db.GetCurrentTime(),
		Updated_at: db.GetCurrentTime(),
	}

}

type GetUserBanner struct {
	Content   entity.Content `json:"content" db:"content"`
	Is_active bool           `json:"is_active" db:"is_active"`
}

type GetBannerById struct {
	Banner_id  int64          `json:"banner_id" db:"id"`
	Tags_ids   pq.Int64Array  `json:"tags_ids" db:"tags_ids"`
	Feature_id int64          `json:"feature_id" db:"feature_id"`
	Content    entity.Content `json:"content" db:"content"`
	Is_active  bool           `json:"is_active" db:"is_active"`
	Created_at time.Time      `json:"created_at" db:"created_at"`
	Updated_at time.Time      `json:"updated_at" db:"updated_at"`
}

type UpdateBanner struct {
	Tags_ids   []int64        `json:"tag_ids" db:"tags" validate:"omitempty,dive,numeric,min=0"`
	Feature_id int64          `json:"feature_id" db:"feature_id" validate:"omitempty,numeric,min=0"`
	Content    entity.Content `json:"content" db:"content" validate:"omitempty,dive,min=0"`
	Is_active  bool           `json:"is_active" db:"is_active" validate:"omitempty"`
}

func (u *UpdateBanner) Validate() error {
	validate := validator.New()
	err := validate.Struct(u)

	if err != nil {
		return err
	}
	return nil
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreatedBannerId struct {
	Banner_id string `json:"banner_id"`
}
