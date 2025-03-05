package agent_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/internal/agent"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func TestGetTaskSuccess(t *testing.T) {
	task := schemas.Task{
		Id:            1,
		Arg1:          10,
		Arg2:          5,
		Operation:     '+',
		OperationTime: 100,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}))
	defer server.Close()

	gotTask, err := agent.Get(server.URL)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if gotTask.Id != task.Id || gotTask.Arg1 != task.Arg1 || gotTask.Arg2 != task.Arg2 {
		t.Errorf("Task mismatch: expected %+v, got %+v", task, gotTask)
	}
}

func TestGetTaskNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	_, err := agent.Get(server.URL)
	if err == nil || err.Error() != "404: задач нет" {
		t.Errorf("Expected '404: задач нет' error, got %v", err)
	}
}

func TestCalcAddition(t *testing.T) {
	task := schemas.Task{
		Arg1:      3,
		Arg2:      7,
		Operation: '+',
	}

	result, err := agent.Calc(&task)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != 10 {
		t.Errorf("Expected 10, got %f", result)
	}
}

func TestCalcDivisionByZero(t *testing.T) {
	task := schemas.Task{
		Arg1:      3,
		Arg2:      0,
		Operation: '/',
	}

	_, err := agent.Calc(&task)
	if err == nil || err != schemas.ErrorDivisionByZero {
		t.Errorf("Expected division by zero error, got %v", err)
	}
}

func TestPostSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var received schemas.TaskRequest
		json.NewDecoder(r.Body).Decode(&received)

		if received.Id != 1 || received.Result != 5 || received.StatusCode != 200 {
			t.Errorf("Unexpected request: %+v", received)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	err := agent.Post(server.URL, 1, 5, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestPostFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err := agent.Post(server.URL, 1, 5, errors.New("some error"))
	if err == nil || err.Error() != "unknown error" {
		t.Errorf("Expected 'unknown error', got %v", err)
	}
}
