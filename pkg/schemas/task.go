package schemas

type Task struct {
	Id            int       `json:"id"`
	Arg1          float64   `json:"arg1"`
	Arg2          float64   `json:"arg2"`
	Operation     Operation `json:"operation"`
	OperationTime int       `json:"operation_time"`
}

type TaskResponse struct {
	Task *Task `json:"task"`
}

type TaskRequest struct {
	Id         int     `json:"id"`
	Result     float64 `json:"result"`
	StatusCode int     `json:"status_code"`
}
