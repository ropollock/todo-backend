package model

type BoardRequest struct {
	ID      string `param:"id" query:"id"`
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
}
