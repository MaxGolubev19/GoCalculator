package parse_test

import (
	"testing"

	"github.com/MaxGolubev19/GoCalculator/pkg/parse"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func TestParseTreeStructure(t *testing.T) {
	tests := []struct {
		expression string
		expected   *schemas.Action
	}{
		{
			expression: "1 + 2",
			expected: &schemas.Action{
				Operation: schemas.AddOperation,
				Left:      &schemas.Action{Value: 1, IsCalculated: true},
				Right:     &schemas.Action{Value: 2, IsCalculated: true},
			},
		},
		{
			expression: "5 - 3",
			expected: &schemas.Action{
				Operation: schemas.SubOperation,
				Left:      &schemas.Action{Value: 5, IsCalculated: true},
				Right:     &schemas.Action{Value: 3, IsCalculated: true},
			},
		},
		{
			expression: "2 * 3 + 4",
			expected: &schemas.Action{
				Operation: schemas.AddOperation,
				Left: &schemas.Action{
					Operation: schemas.MulOperation,
					Left:      &schemas.Action{Value: 2, IsCalculated: true},
					Right:     &schemas.Action{Value: 3, IsCalculated: true},
				},
				Right: &schemas.Action{Value: 4, IsCalculated: true},
			},
		},
		{
			expression: "(1 + 2) * 3",
			expected: &schemas.Action{
				Operation: schemas.MulOperation,
				Left: &schemas.Action{
					Operation: schemas.AddOperation,
					Left:      &schemas.Action{Value: 1, IsCalculated: true},
					Right:     &schemas.Action{Value: 2, IsCalculated: true},
				},
				Right: &schemas.Action{Value: 3, IsCalculated: true},
			},
		},
		{
			expression: "3 + 4 * (2 - 1)",
			expected: &schemas.Action{
				Operation: schemas.AddOperation,
				Left:      &schemas.Action{Value: 3, IsCalculated: true},
				Right: &schemas.Action{
					Operation: schemas.MulOperation,
					Left:      &schemas.Action{Value: 4, IsCalculated: true},
					Right: &schemas.Action{
						Operation: schemas.SubOperation,
						Left:      &schemas.Action{Value: 2, IsCalculated: true},
						Right:     &schemas.Action{Value: 1, IsCalculated: true},
					},
				},
			},
		},
	}

	for _, test := range tests {
		actions, err := parse.New(test.expression)
		if err != nil {
			t.Errorf("Parsing failed for expression %q: %v", test.expression, err)
			continue
		}

		lastAction := (*actions)[len(*actions)-1]
		if !compareActions(lastAction, test.expected) {
			t.Errorf("Incorrect parse tree for %q\nExpected: %+v\nGot: %+v", test.expression, test.expected, lastAction)
		}
	}
}

func TestParseInvalidExpressions(t *testing.T) {
	tests := []string{
		"5 * (3 -",
		"7 * * 3",
		"9 + (3 /)",
		"(1 + 2)) * 3",
		"((3 + 4) * 2",
	}

	for _, expr := range tests {
		_, err := parse.New(expr)
		if err == nil {
			t.Errorf("Expected an error while parsing %q, but parsing succeeded", expr)
		}
	}
}

func compareActions(a, b *schemas.Action) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.IsCalculated != b.IsCalculated || a.Value != b.Value || a.Operation != b.Operation {
		return false
	}
	return compareActions(a.Left, b.Left) && compareActions(a.Right, b.Right)
}
