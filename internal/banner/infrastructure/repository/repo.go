package repository

import (
	"banners/internal/banner/domain/entity"
	"banners/internal/banner/infrastructure/dto"
	"banners/pkg/db"
	slogger "banners/pkg/logger"
	"encoding/json"
	"net/url"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type BannerRepository interface {
	CreateBanner(banner dto.CreateBanner) (int64, error)
	GetBannersList(params url.Values) ([]dto.GetBannerById, error)
	GetUserBanner(params url.Values) (entity.Content, bool, error)
	UpdateBanner(id int64, banner dto.UpdateBanner) error
	DeleteBanner(id int64) error
	GetBannerById(id int64) (dto.GetBannerById, error)
}

const (
	None          int64 = -1
	LimitDefault  int64 = 5
	OffsetDefault int64 = 0
)

type BannerRepo struct {
	db    *sqlx.DB
	cache *cache.Cache
}

func NewBannerRepository(db *sqlx.DB, cache *cache.Cache) BannerRepository {
	return &BannerRepo{db: db, cache: cache}
}

func (b *BannerRepo) CreateBanner(banner dto.CreateBanner) (int64, error) {
	tx, err := b.db.Begin()
	if err != nil {
		return -1, err
	}

	db_banner := banner.DtoToEntity()

	row := tx.QueryRow("INSERT INTO banners (feature_id, title, text, url, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", db_banner.Feature_id,
		db_banner.Content.Title, db_banner.Content.Text, db_banner.Content.Url, db_banner.Is_active, db_banner.Created_at, db_banner.Updated_at)
	err = row.Scan(&db_banner.Banner_id)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	for tagID := range db_banner.Tags_ids {
		_, err = tx.Exec("INSERT INTO bannertags (banner_id, tag_id) VALUES ($1, $2)", db_banner.Banner_id, tagID)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return db_banner.Banner_id, nil
}

func (b *BannerRepo) GetUserBanner(params url.Values) (entity.Content, bool, error) {
	var content entity.Content
	var banner dto.GetUserBanner

	tag_id, feature_id := params.Get("tag_id"), params.Get("feature_id")

	tagID, fID, err := convertParamsToInt(tag_id, feature_id)

	if err != nil {
		return content, false, err
	}
	use_last_revision := params.Get("use_last_revision")
	pairFeatureAndTag := tag_id + feature_id

	if use_last_revision == "false" || len(use_last_revision) == 0 {

		data, ok := b.cache.Get(pairFeatureAndTag)

		if !ok {

			if err := b.db.Get(&banner, db.GetBannerByFeatureAndTag, fID, tagID); err != nil {
				return content, false, err

			}

			b.cache.Set(pairFeatureAndTag, banner, cache.DefaultExpiration)
			return banner.Content, banner.Is_active, nil
		}

		bdata, _ := json.Marshal(data)
		json.Unmarshal(bdata, &banner)
		return banner.Content, banner.Is_active, nil

	} else if use_last_revision == "true" {

		if err := b.db.Get(&banner, db.GetBannerByFeatureAndTag, fID, tagID); err != nil {
			return content, false, err
		}

		b.cache.Set(pairFeatureAndTag, banner, cache.DefaultExpiration)
		return banner.Content, banner.Is_active, nil

	} else {
		return content, false, ErrInvalidParamType
	}

}

func (b *BannerRepo) GetBannersList(params url.Values) ([]dto.GetBannerById, error) {
	var banners []dto.GetBannerById

	paramArr, err := validateParams(params)

	if err != nil {
		return nil, err
	}

	feature_id, tag_id, offset, limit := paramArr[0], paramArr[1], paramArr[2], paramArr[3]

	validateOffsetLimit(&limit, &offset)

	if feature_id != None && tag_id != None {
		err = b.db.Select(&banners, db.GetBannersByFeatureAndTag, tag_id, feature_id, limit, offset)
		if err != nil {
			return nil, err
		}
	} else if feature_id != None {
		err = b.db.Select(&banners, db.GetBannersByFeature, feature_id, limit, offset)
		if err != nil {
			return nil, err
		}
	} else if tag_id != None {
		err = b.db.Select(&banners, db.GetBannersByTag, tag_id, limit, offset)
		if err != nil {
			return nil, err
		}
	} else {
		err = b.db.Select(&banners, db.GetBanners, limit, offset)
		if err != nil {
			return nil, err
		}
	}

	return banners, nil
}

func (b *BannerRepo) UpdateBanner(id int64, banner dto.UpdateBanner) error {

	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`UPDATE banners SET feature_id = COALESCE($1, feature_id), title = COALESCE($2, title), text = COALESCE($3, text), url = COALESCE($4, url),
					 is_active = COALESCE($5, is_active), updated_at = $6`,
		banner.Feature_id, banner.Content.Title, banner.Content.Text, banner.Content.Url, banner.Is_active, db.GetCurrentTime())

	if err != nil {
		slogger.Logger.Info("error while update banners: ", "msg", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM bannertags WHERE banner_id = $1", id)
	if err != nil {
		slogger.Logger.Info("error while deletion from bannertags: ", "msg", err)
		return err
	}

	for _, tagID := range banner.Tags_ids {
		_, err = tx.Exec(`INSERT INTO bannertags (banner_id, tag_id) VALUES ($1, $2)`, id, tagID)
		if err != nil {
			slogger.Logger.Info("error while UpdateBanner insertion: %v", err)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()

	if err != nil {
		slogger.Logger.Debug("error while commit: ", "msg", err)
		return err
	}
	return nil
}

func (b *BannerRepo) DeleteBanner(id int64) error {
	_, err := b.db.Exec("DELETE FROM banners WHERE id = $1", id)
	if err != nil {
		slogger.Logger.Info("error while banner deletion: %s", err)
		return err
	}
	return nil
}

func (b *BannerRepo) GetBannerById(id int64) (dto.GetBannerById, error) {
	var banner dto.GetBannerById

	err := b.db.Get(&banner, db.GetBannerById, id)

	if err != nil {
		slogger.Logger.Info("error while querying banner: %v", err)
		return banner, err
	}
	return banner, nil
}

func validateParams(params url.Values) ([]int64, error) {

	paramsArr := []string{params.Get("feature_id"), params.Get("tag_id"), params.Get("offset"), params.Get("limit")}

	var numerals []int64

	for _, val := range paramsArr {

		if len(val) == 0 {
			numerals = append(numerals, None)
		} else {
			num, err := strconv.Atoi(val)
			if err != nil {
				return nil, ErrInvalidParamType
			}
			if num < 0 {
				return nil, ErrNegativeId
			}
			numerals = append(numerals, int64(num))
		}
	}
	return numerals, nil
}

func convertParamsToInt(t string, f string) (int64, int64, error) {
	var t_id, f_id int
	t_id, err := strconv.Atoi(t)
	if err != nil {
		return 0, 0, ErrInvalidParamType
	}
	if t_id < 0 {
		return 0, 0, ErrNegativeId
	}
	f_id, err = strconv.Atoi(f)
	if err != nil {
		return 0, 0, ErrInvalidParamType
	}
	if f_id < 0 {
		return 0, 0, ErrNegativeId
	}
	return int64(t_id), int64(f_id), nil

}

var ErrInvalidParamType = errors.New("params should be integer")
var ErrNegativeId = errors.New("id could not be negative")

func validateOffsetLimit(limit, offset *int64) {
	if *offset == None && *limit == None {
		*offset, *limit = OffsetDefault, LimitDefault
	} else if *offset == None {
		*offset = OffsetDefault
	} else if *limit == None {
		*limit = LimitDefault
	}
}
