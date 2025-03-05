package schemas

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Id int `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
