package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/internal/application"
)

func post(t *testing.T, url, expression string) *http.Response {
	request := application.Request{Expression: expression}
	requestJSON, _ := json.Marshal(request)
	resp, err := http.Post(url+"/api/v1/calculate", "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	return resp
}

func TestApplicationSuccess(t *testing.T) {
	app := application.New()
	server := httptest.NewServer(http.HandlerFunc(app.CalcHandler))
	defer server.Close()

	successTests := []struct {
		name           string
		expression     string
		expectedCode   int
		expectedResult float64
	}{
		{
			name:           "uno",
			expression:     "42",
			expectedCode:   http.StatusOK,
			expectedResult: 42,
		},
		{
			name:           "sum",
			expression:     "1+1",
			expectedCode:   http.StatusOK,
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedCode:   http.StatusOK,
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedCode:   http.StatusOK,
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedCode:   http.StatusOK,
			expectedResult: 0.5,
		},
	}

	for _, test := range successTests {
		t.Run(test.name, func(t *testing.T) {
			resp := post(t, server.URL, test.expression)
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Fatalf("Expected status %d, but got %d", test.expectedCode, resp.StatusCode)
			}

			var response application.Response
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %s", err.Error())
			}

			if response.Result != test.expectedResult {
				t.Fatalf("Expected result %f, but got %f", test.expectedResult, response.Result)
			}
		})
	}
}

func TestApplicationFail(t *testing.T) {
	app := application.New()
	server := httptest.NewServer(http.HandlerFunc(app.CalcHandler))
	defer server.Close()

	failTests := []struct {
		name          string
		expression    string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "empty",
			expression:    "",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Expression is not valid",
		},
		{
			name:          "letters",
			expression:    "1+a",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Expression is not valid",
		},
		{
			name:          "operant at the end",
			expression:    "1+1*",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Expression is not valid",
		},
		{
			name:          "double operation",
			expression:    "2+2**2",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Expression is not valid",
		},
		{
			name:          "incorrect priority",
			expression:    "((2+2-*(2",
			expectedCode:  http.StatusUnprocessableEntity,
			expectedError: "Expression is not valid",
		},
		{
			name:          "division by zero",
			expression:    "42/0",
			expectedCode:  http.StatusInternalServerError,
			expectedError: "Internal server error",
		},
	}

	for _, test := range failTests {
		t.Run(test.name, func(t *testing.T) {
			resp := post(t, server.URL, test.expression)
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Fatalf("Expected status %d, but got %d", test.expectedCode, resp.StatusCode)
			}

			var response application.ErrorResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %s", err.Error())
			}

			if response.Error != test.expectedError {
				t.Fatalf("Expected error \"%s\", but got \"%s\"", test.expectedError, response.Error)
			}
		})
	}
}
