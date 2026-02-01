package dto

type BadReqErrorResponse struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"bad request from client"`
}

type NotFoundErrorResponse struct {
	Code  int    `json:"code" example:"404"`
	Error string `json:"error" example:"task not found"`
}

type InternalErrorResponse struct {
	Code  int    `json:"code" example:"500"`
	Error string `json:"error" example:"internal server error"`
}
