package orchestrator

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	pb "github.com/MaxGolubev19/GoCalculator/pkg/proto"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

type Config struct {
	PublicPort string
	GrpcPort   string

	SecretKey string

	TimeAdditionMS       int
	TimeSubstractionMS   int
	TimeMultiplicationMS int
	TimeDivisionsMS      int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.PublicPort = os.Getenv("PUBLIC_PORT")
	if config.PublicPort == "" {
		config.PublicPort = "8080"
	}

	config.GrpcPort = os.Getenv("GRPC_PORT")
	if config.GrpcPort == "" {
		config.GrpcPort = "50051"
	}

	config.SecretKey = os.Getenv("SECRET_KEY")
	if config.SecretKey == "" {
		config.SecretKey = "super secret key"
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

	db *sql.DB

	actions map[int]*schemas.Action

	pb.TaskServiceServer
	tasks   []pb.Task
	taskId  int
	muTasks sync.Mutex
}

func New() *Orchestrator {
	return &Orchestrator{
		config:  ConfigFromEnv(),
		actions: make(map[int]*schemas.Action, 0),
		tasks:   make([]pb.Task, 0),
	}
}

func (o *Orchestrator) Run() error {
	go o.RunTaskServer()

	db, err := InitDB()
	if err != nil {
		return err
	}
	o.db = db
	defer o.db.Close()

	if expressions, err := o.GetExpressionsInProgress(); err != nil {
		return err
	} else {
		for _, expr := range expressions {
			o.ParseExpression(expr.Id, expr.Expression)
		}
	}

	http.HandleFunc("/api/v1/register", o.RegisterHandler)
	http.HandleFunc("/api/v1/login", o.LoginHandler)

	http.Handle("/api/v1/calculate", o.CheckJWT(http.HandlerFunc(o.CalculateHandler)))
	http.Handle("/api/v1/expressions/", o.CheckJWT(http.HandlerFunc(o.ExpressonByIdHandler)))
	http.Handle("/api/v1/expressions", o.CheckJWT(http.HandlerFunc(o.ExpressonsHandler)))

	return http.ListenAndServe(":"+o.config.PublicPort, nil)
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
			o.SetExpressionError(id)
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
		o.SetExpressionError(id)
		return
	}

	o.SetExpressionDone(id, (*actions)[index-1].Value)
}

// For tests
func (o *Orchestrator) GetAction(id int) *schemas.Action {
	if action, exists := o.actions[id]; exists {
		return action
	}
	return nil
}
