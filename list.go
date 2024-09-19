package driplimit

import "github.com/go-playground/validator/v10"

// ListMetadata represents the metadata of a list.
type ListMetadata struct {
	Page     int `json:"page"`
	Limit    int `json:"limit"`
	LastPage int `json:"last_page"`
}

// NewListMetadata creates a new list metadata. It calculates the last page based on the total count.
func NewListMetadata(payload ListPayload, totalCount int) ListMetadata {
	addPage := 0
	if totalCount%payload.Limit > 0 || totalCount == 0 {
		addPage = 1
	}
	return ListMetadata{
		Page:     payload.Page,
		Limit:    payload.Limit,
		LastPage: totalCount/payload.Limit + addPage,
	}
}

// ListPayload represents the payload for listing items.
type ListPayload struct {
	Page  int `query:"page" json:"page" validate:"gte=1" description:"The page number"`
	Limit int `query:"limit" json:"limit" validate:"gte=1,lte=100" description:"The number of items per page"`
}

// Validate validates the list payload.
func (lp *ListPayload) Validate(validator *validator.Validate) error {
	if lp.Limit == 0 {
		lp.Limit = 10
	}
	if lp.Page == 0 {
		lp.Page = 1
	}
	return validator.Struct(lp)
}

// Offset returns the offset based on the page and limit.
func (lp *ListPayload) Offset() int {
	return (lp.Page - 1) * lp.Limit
}
