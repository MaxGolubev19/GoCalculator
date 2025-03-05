package parse

import (
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

type Parsing struct {
	expression string
	actions    []*schemas.Action
	length     int
}

func New(expression string) (*[]*schemas.Action, error) {
	parsing := Parsing{
		expression: expression,
		length:     len(expression),
	}

	err := parsing.Parse()
	if err != nil {
		return nil, err
	}

	return &parsing.actions, nil
}

func (p *Parsing) NewNumber(n float64) *schemas.Action {
	return &schemas.Action{
		Value:        n,
		IsCalculated: true,
	}
}

func (p *Parsing) NewExpression(operation schemas.Operation, left *schemas.Action, right *schemas.Action) *schemas.Action {
	return &schemas.Action{
		Operation:    operation,
		Left:         left,
		Right:        right,
		IsCalculated: false,
	}
}

func (p *Parsing) Parse() error {
	a, index, err := p.get(0, p.length)
	if err != nil {
		return err
	}

	var b *schemas.Action

	for index < p.length {
		if p.expression[index] == ' ' {
			index++
			continue
		}

		if p.expression[index] == '+' {
			b, index, err = p.get(index+1, p.length)
			if err != nil {
				return err
			}
			a = p.NewExpression(schemas.AddOperation, a, b)
			p.actions = append(p.actions, a)
			continue
		}

		if p.expression[index] == '-' {
			b, index, err = p.get(index+1, p.length)
			if err != nil {
				return err
			}
			a = p.NewExpression(schemas.SubOperation, a, b)
			p.actions = append(p.actions, a)
			continue
		}

		return ErrorIncorrectExpression
	}

	return nil
}

func (p *Parsing) get(index int, length int) (*schemas.Action, int, error) {
	var n float64
	var a *schemas.Action
	var isNumber, isNegative, isEndNumber bool
	var isOpen int

	for {
		if index == length {
			if isNumber {
				a = p.NewNumber(n)
				p.actions = append(p.actions, a)
				isNumber = false
			}
			if a != nil && isOpen == 0 {
				return a, index, nil
			} else {
				return nil, 0, ErrorIncorrectExpression
			}
		}

		if p.expression[index] == ' ' {
			isEndNumber = true
			index++
			continue
		}

		if p.expression[index] == '(' {
			isOpen++
			index++
			continue
		}

		if p.expression[index] == ')' {
			if isOpen > 0 {
				isOpen--
				index++
			} else {
				if isNumber {
					a = p.NewNumber(n)
					p.actions = append(p.actions, a)
					isNumber = false
				}
				if a != nil {
					return a, index, nil
				}
				return nil, 0, ErrorIncorrectExpression
			}
			continue
		}

		if p.expression[index] == '+' {
			if isNumber {
				a = p.NewNumber(n)
				p.actions = append(p.actions, a)
				isNumber = false
			}

			if a != nil && isOpen == 0 {
				return a, index, nil
			} else if a != nil {
				b, new_index, err := p.get(index+1, length)
				if err != nil {
					return nil, 0, err
				}
				a = p.NewExpression(schemas.AddOperation, a, b)
				p.actions = append(p.actions, a)
				index = new_index
			} else {
				index++
			}
			continue
		}

		if p.expression[index] == '-' {
			if isNumber {
				a = p.NewNumber(n)
				p.actions = append(p.actions, a)
				isNumber = false
			}

			if a != nil && isOpen == 0 {
				return a, index, nil
			} else if a != nil {
				b, new_index, err := p.get(index+1, length)
				if err != nil {
					return nil, 0, err
				}
				a = p.NewExpression(schemas.SubOperation, a, b)
				p.actions = append(p.actions, a)
				index = new_index
			} else {
				isNegative = true
				index++
			}
			continue
		}

		if p.expression[index] == '*' {
			if isNumber {
				a = p.NewNumber(n)
				p.actions = append(p.actions, a)
				isNumber = false
			}

			if a != nil {
				b, new_index, err := p.get(index+1, length)
				if err != nil {
					return nil, 0, err
				}
				a = p.NewExpression(schemas.MulOperation, a, b)
				p.actions = append(p.actions, a)
				index = new_index
			} else {
				return nil, 0, ErrorIncorrectExpression
			}
			continue
		}

		if p.expression[index] == '/' {
			if isNumber {
				a = p.NewNumber(n)
				p.actions = append(p.actions, a)
				isNumber = false
			}

			if a != nil {
				b, new_index, err := p.get(index+1, length)
				if err != nil {
					return nil, 0, err
				}
				a = p.NewExpression(schemas.DivOperation, a, b)
				p.actions = append(p.actions, a)
				index = new_index
			} else {
				return nil, 0, ErrorIncorrectExpression
			}
			continue
		}

		if p.expression[index]-'0' < 10 {
			if isNumber && isEndNumber {
				return nil, 0, ErrorIncorrectExpression
			} else {
				n = n*10 + float64((p.expression)[index]-'0')
				if !isNumber {
					if isNegative {
						n *= -1
						isNegative = false
					}
					isNumber = true
					isEndNumber = false
				}
				index++
			}
			continue
		}

		return nil, 0, ErrorIncorrectExpression
	}
}
