package model

type ListRequest struct {
	ID      string `param:"id" query:"id"`
	Name    string `json:"name,omitempty"`
	Order   int32  `json:"order,omitempty"`
	BoardID string `param:"board_id" query:"board_id"`
}
