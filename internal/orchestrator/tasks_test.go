package orchestrator_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/internal/orchestrator"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func TestTaskHandlerGet(t *testing.T) {
	orch := orchestrator.New()
	action := &schemas.Action{
		Operation: schemas.AddOperation,
		Left:      &schemas.Action{Value: 1, IsCalculated: true},
		Right:     &schemas.Action{Value: 2, IsCalculated: true},
	}
	orch.AddTask(action)

	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response schemas.Task
	json.Unmarshal(rr.Body.Bytes(), &response)
	if response.Operation != schemas.AddOperation {
		t.Errorf("Expected operation %v, got %v", schemas.AddOperation, response.Operation)
	}
}

func TestTaskHandlerGetEmpty(t *testing.T) {
	orch := orchestrator.New()

	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestTaskHandlerPostSuccess(t *testing.T) {
	orch := orchestrator.New()
	action := &schemas.Action{
		Operation: schemas.AddOperation,
		Left:      &schemas.Action{Value: 1, IsCalculated: true},
		Right:     &schemas.Action{Value: 2, IsCalculated: true},
	}
	taskId := orch.AddTask(action)

	body, _ := json.Marshal(schemas.TaskRequest{
		Id:         taskId,
		Result:     3,
		StatusCode: 200,
	})
	req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	resultAction := orch.GetAction(taskId)
	if resultAction == nil || !resultAction.IsCalculated || resultAction.Value != 3 {
		t.Errorf("Expected result 3, got %f", resultAction.Value)
	}
}

func TestTaskHandlerPostInvalidData(t *testing.T) {
	orch := orchestrator.New()

	body := []byte(`{"id": "invalid", "result": 3}`)
	req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", rr.Code)
	}
}

func TestTaskHandlerPostNotFound(t *testing.T) {
	orch := orchestrator.New()

	body := []byte(`{"id": 999, "result": 3, "statusCode": 200}`)
	req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", rr.Code)
	}
}

func TestTaskHandlerPostErrorStatus(t *testing.T) {
	orch := orchestrator.New()
	action := &schemas.Action{
		Operation: schemas.AddOperation,
		Left:      &schemas.Action{Value: 1, IsCalculated: true},
		Right:     &schemas.Action{Value: 2, IsCalculated: true},
	}
	taskId := orch.AddTask(action)

	body := []byte(`{"id": ` + strconv.Itoa(taskId) + `, "result": 3, "statusCode": 500}`)
	req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(orch.TaskHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	if !orch.GetAction(taskId).IsError {
		t.Errorf("Expected IsError to be true, but got false")
	}
}
