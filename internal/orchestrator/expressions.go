package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func (o *Orchestrator) ExpressonsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := schemas.ExpressionsResponse{
		Expressions: o.expressions,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		error := schemas.ErrorResponse{Error: "Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
}

func (o *Orchestrator) ExpressonByIdHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/api/v1/expressions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println(id)

	if o.expressionId <= id {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := o.expressions[id]

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		error := schemas.ErrorResponse{Error: "Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
}

func (o *Orchestrator) AddExpression(actions *[]*schemas.Action) int {
	o.muExpression.Lock()

	expressionId := o.expressionId
	o.expressionId++

	o.expressions = append(o.expressions, schemas.Expression{
		Id:     expressionId,
		Status: schemas.IN_PROGRESS,
	})

	o.muExpression.Unlock()

	go o.worker(expressionId, actions)

	return expressionId
}
