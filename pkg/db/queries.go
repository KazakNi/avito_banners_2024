package db

var GetBannersByFeatureAndTag = `SELECT banners.id, array_agg(bt.tag_id) as tags_ids, feature_id, 
json_build_object('title', title, 'text', text, 'url', url) as content, created_at, updated_at
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id
WHERE bt.tag_id = $1 AND feature_id = $2
GROUP BY banners.id, feature_id, created_at, updated_at, title, text, url
LIMIT $3
OFFSET $4`

var GetBannersByFeature = `SELECT banners.id, array_agg(bt.tag_id) as tags_ids, feature_id, 
json_build_object('title', title, 'text', text, 'url', url) as content, created_at, updated_at
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id
WHERE feature_id = $1
GROUP BY banners.id, feature_id, created_at, updated_at, title, text, url
LIMIT $2
OFFSET $3`

var GetBannersByTag = `SELECT banners.id, array_agg(bt.tag_id) as tags_ids, feature_id, 
json_build_object('title', title, 'text', text, 'url', url) as content, created_at, updated_at
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id
WHERE bt.tag_id = $1
GROUP BY banners.id, feature_id, created_at, updated_at, title, text, url
LIMIT $2
OFFSET $3`

var GetBanners = `SELECT banners.id, array_agg(bt.tag_id) as tags_ids, feature_id, 
json_build_object('title', title, 'text', text, 'url', url) as content, created_at, updated_at
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id
GROUP BY banners.id, feature_id, created_at, updated_at, title, text, url
LIMIT $1
OFFSET $2
`

var GetBannerById = `SELECT banners.id, array_agg(bt.tag_id) as tags_ids, feature_id, 
json_build_object('title', title, 'text', text, 'url', url) as content, created_at, updated_at
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id
WHERE banners.id = $1
GROUP BY banners.id, feature_id, created_at, updated_at, title, text, url
`

var GetBannerByFeatureAndTag = `SELECT is_active,
json_build_object('title', title, 'text', text, 'url', url) as content
FROM banners 
JOIN bannertags bt
ON bt.banner_id = banners.id AND bt.tag_id = $2
WHERE banners.feature_id = $1
GROUP BY title, text, url, is_active
`
