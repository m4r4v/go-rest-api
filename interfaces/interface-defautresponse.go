package interfaces

type IDefaultResponse struct {
	Status  int    `json:"status-code"`
	Message string `json:"message"`
}
