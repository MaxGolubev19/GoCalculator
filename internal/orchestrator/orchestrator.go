package orchestrator

import (
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

type Config struct {
	Port string

	TimeAdditionMS       int
	TimeSubstractionMS   int
	TimeMultiplicationMS int
	TimeDivisionsMS      int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.Port = os.Getenv("PORT")
	if config.Port == "" {
		config.Port = "8080"
	}

	time, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil {
		config.TimeAdditionMS = 100
	} else {
		config.TimeAdditionMS = time
	}

	time, err = strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil {
		config.TimeSubstractionMS = 100
	} else {
		config.TimeSubstractionMS = time
	}

	time, err = strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err != nil {
		config.TimeMultiplicationMS = 100
	} else {
		config.TimeMultiplicationMS = time
	}

	time, err = strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err != nil {
		config.TimeDivisionsMS = 100
	} else {
		config.TimeDivisionsMS = time
	}

	return config
}

type Orchestrator struct {
	config *Config

	expressions  []schemas.Expression
	expressionId int
	muExpression sync.Mutex

	actions map[int]*schemas.Action
	tasks   []schemas.Task
	taskId  int
	muTask  sync.Mutex
}

func New() *Orchestrator {
	return &Orchestrator{
		config:      ConfigFromEnv(),
		expressions: make([]schemas.Expression, 0),
		actions:     make(map[int]*schemas.Action, 0),
		tasks:       make([]schemas.Task, 0),
	}
}

func (o *Orchestrator) Run() error {
	http.HandleFunc("/api/v1/calculate", o.CalculateHandler)
	http.HandleFunc("/api/v1/expressions/", o.ExpressonByIdHandler)
	http.HandleFunc("/api/v1/expressions", o.ExpressonsHandler)
	http.HandleFunc("/internal/task", o.TaskHandler)
	return http.ListenAndServe(":"+o.config.Port, nil)
}

func (o *Orchestrator) worker(id int, actions *[]*schemas.Action) {
	index := 0

	for {
		if index == len(*actions) {
			break
		}

		if (*actions)[index].IsCalculated {
			index++
			continue
		}

		if (*actions)[index].Left.IsError || (*actions)[index].Right.IsError {
			o.expressions[id].Status = schemas.ERROR
			return
		}

		if (*actions)[index].Left.IsCalculated && (*actions)[index].Right.IsCalculated {
			o.AddTask((*actions)[index])
			index++
			continue
		}

		time.Sleep(100 * time.Millisecond)
	}

	for !(*actions)[index-1].IsCalculated && !(*actions)[index-1].IsError {
		time.Sleep(100 * time.Millisecond)
	}

	if (*actions)[index-1].IsError {
		o.expressions[id].Status = schemas.ERROR
		return
	}

	o.expressions[id].Status = schemas.DONE
	o.expressions[id].Result = (*actions)[index-1].Value
}

// For tests
func (o *Orchestrator) GetAction(id int) *schemas.Action {
	if action, exists := o.actions[id]; exists {
		return action
	}
	return nil
}

func (o *Orchestrator) GetExpression(id int) *schemas.Expression {
	if id < 0 || id >= len(o.expressions) {
		return nil
	}
	return &o.expressions[id]
}
