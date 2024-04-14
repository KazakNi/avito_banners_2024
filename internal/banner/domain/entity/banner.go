package entity

import (
	"encoding/json"
	"time"
)

type Banner struct {
	Banner_id  int64     `json:"banner_id" db:"id"`
	Tags_ids   []int64   `json:"tags_ids" db:"tags"`
	Feature_id int64     `json:"feature_id" db:"feature_id"`
	Content    Content   `json:"content" db:"content"`
	Is_active  bool      `json:"is_active" db:"is_active"`
	Created_at time.Time `json:"created_at" db:"created_at"`
	Updated_at time.Time `json:"updated_at" db:"updated_at"`
}

type Content struct {
	Title string `json:"title" db:"title"`
	Text  string `json:"text" db:"text"`
	Url   string `json:"url" db:"url"`
}

func (c *Content) Scan(src interface{}) error {
	err := json.Unmarshal(src.([]byte), &c)
	return err
}
