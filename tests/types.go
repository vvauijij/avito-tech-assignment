package tests

type (
	Banner struct {
		BannerID *int64 `json:"banner_id,omitempty"`

		FeatureID *int64  `json:"feature_id,omitempty"`
		TagIDs    []int64 `json:"tag_ids,omitempty"`

		Content  Content `json:"content,omitempty"`
		IsActive *bool   `json:"is_active,omitempty"`
	}

	Content map[string]any
	Filter  map[string]any
)
