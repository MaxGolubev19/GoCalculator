package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func (o *Orchestrator) TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		response := o.GetTask()
		if response == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			error := schemas.ErrorResponse{Error: "Internal server error"}
			json.NewEncoder(w).Encode(error)
			return
		}

		return
	}

	if r.Method == "POST" {
		var request schemas.TaskRequest
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			error := schemas.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(error)
			return
		}

		if o.taskId <= request.Id {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		o.actions[request.Id].Value = request.Result
		o.actions[request.Id].IsCalculated = true

		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (o *Orchestrator) AddTask(action *schemas.Action) int {
	taskId := o.taskId
	o.taskId++

	o.actions[taskId] = action

	var operationTime int
	switch action.Operation {
	case schemas.AddOperation:
		operationTime = o.config.TimeAdditionMS
	case schemas.SubOperation:
		operationTime = o.config.TimeSubstractionMS
	case schemas.MulOperation:
		operationTime = o.config.TimeMultiplicationMS
	case schemas.DivOperation:
		operationTime = o.config.TimeDivisionsMS
	}

	o.tasks = append(o.tasks, schemas.Task{
		Id:            taskId,
		Arg1:          action.Left.Value,
		Arg2:          action.Right.Value,
		Operation:     action.Operation,
		OperationTime: operationTime,
	})

	return taskId
}

func (o *Orchestrator) GetTask() *schemas.Task {
	if len(o.tasks) == 0 {
		return nil
	}

	task := &o.tasks[0]
	o.tasks = o.tasks[1:]
	return task
}
