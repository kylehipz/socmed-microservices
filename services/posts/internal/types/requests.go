package types

type CreateOrUpdatePostRequest struct {
	Content string `json:"content" validate:"required"`
}
