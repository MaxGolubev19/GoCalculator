package calc

func get(expression *string, index int, length int) (float64, int, error) {
	var a float64
	var isNumber, isNegative, isEndNumber, isOpen bool

	for {
		if index == length {
			if isNumber && !isOpen {
				return a, index, nil
			} else {
				return 0, 0, ErrorIncorrectExpression
			}
		} else if (*expression)[index] == ' ' {
			isEndNumber = true
			index++
		} else if (*expression)[index] == '(' {
			if !isOpen {
				isOpen = true
				index++
			} else {
				return 0, 0, ErrorIncorrectExpression
			}
		} else if (*expression)[index] == ')' {
			if isOpen {
				isOpen = false
				index++
			} else {
				return a, index, nil
			}
		} else if (*expression)[index] == '+' {
			if isNumber && !isOpen {
				return a, index, nil
			} else if isNumber {
				b, new_index, err := get(expression, index+1, length)
				if err != nil {
					return 0, 0, err
				}
				a += b
				index = new_index
			} else {
				index++
			}
		} else if (*expression)[index] == '-' {
			if isNumber && !isOpen {
				return a, index, nil
			} else if isNumber {
				b, new_index, err := get(expression, index+1, length)
				if err != nil {
					return 0, 0, err
				}
				a -= b
				index = new_index
			} else {
				isNegative = true
				index++
			}
		} else if (*expression)[index] == '*' {
			if isNumber {
				b, new_index, err := get(expression, index+1, length)
				if err != nil {
					return 0, 0, err
				}
				a *= b
				index = new_index
			} else {
				return 0, 0, ErrorIncorrectExpression
			}
		} else if (*expression)[index] == '/' {
			if isNumber {
				b, new_index, err := get(expression, index+1, length)
				if err != nil {
					return 0, 0, err
				} else if b == 0 {
					return 0, 0, ErrorDivisionByZero
				}
				a /= b
				index = new_index
			}
		} else if (*expression)[index]-'0' < 10 {
			if isNumber && isEndNumber {
				return 0, 0, ErrorIncorrectExpression
			} else {
				a = a*10 + float64((*expression)[index]-'0')
				if !isNumber {
					if isNegative {
						a *= -1
						isNegative = false
					}
					isNumber = true
					isEndNumber = false
				}
				index++
			}
		} else {
			return 0, 0, ErrorIncorrectExpression
		}
	}
}

func Calc(expression string) (float64, error) {
	length := len(expression)
	a, index, err := get(&expression, 0, length)
	var b float64
	if err != nil {
		return 0, err
	}

	for index < length {
		if expression[index] == ' ' {
			index++
		} else if expression[index] == '+' {
			b, index, err = get(&expression, index+1, length)
			if err != nil {
				return 0, err
			}
			a += b
		} else if expression[index] == '-' {
			b, index, err = get(&expression, index+1, length)
			if err != nil {
				return 0, err
			}
			a -= b
		} else {
			return 0, ErrorIncorrectExpression
		}
	}

	return a, nil
}
