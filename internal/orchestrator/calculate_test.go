package orchestrator_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/internal/orchestrator"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func TestCalculateHandlerSuccess(t *testing.T) {
	orch := orchestrator.New()

	body, _ := json.Marshal(schemas.CalculateRequest{
		Expression: "2 + 2",
	})
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.CalculateHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var response schemas.CalculateResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Id < 0 {
		t.Errorf("Expected valid expression ID, got %d", response.Id)
	}
}

func TestCalculateHandlerInvalidExpression(t *testing.T) {
	orch := orchestrator.New()

	body, _ := json.Marshal(schemas.CalculateRequest{
		Expression: "2 +",
	})
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.CalculateHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
}

func TestCalculateHandlerInvalidJSON(t *testing.T) {
	orch := orchestrator.New()

	body := []byte(`{"invalid_json"}`)
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.CalculateHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}
}
