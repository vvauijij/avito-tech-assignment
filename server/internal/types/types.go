package types

type (
	BannerFilter struct {
		FeatureID *int64 `bson:"feature_id,omitempty"`
		TagID     *int64 `bson:"tag_ids,omitempty"`
	}

	Banner struct {
		BannerID int64 `json:"banner_id" bson:"banner_id"`

		FeatureID *int64  `json:"feature_id,omitempty" bson:"feature_id,omitempty"`
		TagIDs    []int64 `json:"tag_ids,omitempty" bson:"tag_ids,omitempty"`

		Content   map[string]any `json:"content,omitempty" bson:"content,omitempty"`
		IsActive  *bool          `json:"is_active,omitempty" bson:"is_active,omitempty"`
		CreatedAt *string        `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt *string        `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}
)
