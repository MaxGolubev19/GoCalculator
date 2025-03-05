package schemas

type Status string

const (
	IN_PROGRESS Status = "IN PROGRESS"
	DONE        Status = "DONE"
	ERROR       Status = "ERROR"
)

type ExpressionRequest struct {
	Id int `json:"id"`
}

type Expression struct {
	Id     int     `json:"id"`
	Status Status  `json:"status"`
	Result float64 `json:"result"`
}

type ExpressionResponse struct {
	Expression Expression `json:"expression"`
}

type ExpressionsResponse struct {
	Expressions []Expression `json:"expressions"`
}
