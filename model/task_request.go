package model

type TaskRequest struct {
	ID      string `param:"id" query:"id"`
	Name    string `json:"name,omitempty"`
	Content string `json:"content,omitempty"`
	Order   int32  `json:"order,omitempty"`
	ListID  string `param:"list_id" query:"list_id"`
	BoardID string `param:"board_id" query:"board_id"`
}
