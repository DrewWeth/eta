package models

type AppDataReqParams struct {
	Comment  string `json:"comment"`
	Username string `json:"username"`
	PostID   string `json:"post_id" schema:"post_id"`
	ParentID string `json:"parent_id" schema:"parent_id"`
	APIToken string `json:"api_token" schema:"api_token"`
}

type GenericResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
