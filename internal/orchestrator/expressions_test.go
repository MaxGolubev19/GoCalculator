package orchestrator_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/internal/orchestrator"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func TestExpressonsHandler(t *testing.T) {
	orch := orchestrator.New()
	orch.AddExpression(&[]*schemas.Action{
		{Value: 42, IsCalculated: true},
	})

	req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.ExpressonsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response schemas.ExpressionsResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response.Expressions) == 0 {
		t.Errorf("Expected at least one expression, but got empty list")
	}
}

func TestExpressonByIdHandlerValidID(t *testing.T) {
	orch := orchestrator.New()
	exprID := orch.AddExpression(&[]*schemas.Action{
		{Value: 42, IsCalculated: true},
	})

	req, _ := http.NewRequest("GET", "/api/v1/expressions/"+strconv.Itoa(exprID), nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.ExpressonByIdHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response schemas.Expression
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Id != exprID {
		t.Errorf("Expected expression ID %d, got %d", exprID, response.Id)
	}
}

func TestExpressonByIdHandlerInvalidID(t *testing.T) {
	orch := orchestrator.New()

	req, _ := http.NewRequest("GET", "/api/v1/expressions/9999", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.ExpressonByIdHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestExpressonByIdHandlerMalformedID(t *testing.T) {
	orch := orchestrator.New()

	req, _ := http.NewRequest("GET", "/api/v1/expressions/abc", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.ExpressonByIdHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", rr.Code)
	}
}
