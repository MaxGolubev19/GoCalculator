package schemas

type Operation rune

const (
	AddOperation Operation = '+'
	SubOperation Operation = '-'
	MulOperation Operation = '*'
	DivOperation Operation = '/'
)

type Action struct {
	Value        float64
	Operation    Operation
	Left         *Action
	Right        *Action
	IsCalculated bool
	IsError      bool
}
