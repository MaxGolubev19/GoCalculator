package calc_test

import (
	"errors"
	"testing"

	"github.com/MaxGolubev19/GoCalculator/pkg/calc"
)

func TestCalcSuccess(t *testing.T) {
	successTests := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "uno",
			expression:     "42",
			expectedResult: 42,
		},
		{
			name:           "sum",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
	}

	for _, test := range successTests {
		t.Run(test.name, func(t *testing.T) {
			res, err := calc.Calc(test.expression)
			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
			if res != test.expectedResult {
				t.Fatalf("Expected result %f, but got %f", test.expectedResult, res)
			}
		})
	}
}

func TestCalcFail(t *testing.T) {
	failTests := []struct {
		name          string
		expression    string
		expectedError error
	}{
		{
			name:          "empty",
			expression:    "",
			expectedError: calc.ErrorIncorrectExpression,
		},
		{
			name:          "letters",
			expression:    "1+a",
			expectedError: calc.ErrorIncorrectExpression,
		},
		{
			name:          "operant at the end",
			expression:    "1+1*",
			expectedError: calc.ErrorIncorrectExpression,
		},
		{
			name:          "double operation",
			expression:    "2+2**2",
			expectedError: calc.ErrorIncorrectExpression,
		},
		{
			name:          "incorrect priority",
			expression:    "((2+2-*(2",
			expectedError: calc.ErrorIncorrectExpression,
		},
		{
			name:          "division by zero",
			expression:    "42/0",
			expectedError: calc.ErrorDivisionByZero,
		},
	}

	for _, test := range failTests {
		t.Run(test.name, func(t *testing.T) {
			val, err := calc.Calc(test.expression)
			if err == nil {
				t.Fatalf("Expected error \"%f\", but got result %f", test.expectedError, val)
			}
			if !errors.Is(test.expectedError, err) {
				t.Fatalf("Expected error \"%f\", but got \"%f\"", test.expectedError, err)
			}
		})
	}
}
