package orchestrator

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/MaxGolubev19/GoCalculator/pkg/parse"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

type Response struct {
	Result float64 `json:"result"`
}

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request schemas.CalculateRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		error := schemas.ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	expr := strings.TrimSpace(request.Expression)
	actions, err := parse.New(expr)
	if errors.Is(err, schemas.ErrorIncorrectExpression) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		error := schemas.ErrorResponse{Error: "Expression is not valid"}
		json.NewEncoder(w).Encode(error)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := schemas.ErrorResponse{Error: "Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}

	id := o.AddExpression(actions)

	response := schemas.CalculateResponse{Id: id}
	data, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		error := schemas.ErrorResponse{Error: "Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
