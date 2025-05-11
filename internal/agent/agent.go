package agent

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	pb "github.com/MaxGolubev19/GoCalculator/pkg/proto"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Port           string
	ComputingPower int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.Port = os.Getenv("GRPC_PORT")
	if config.Port == "" {
		config.Port = "8080"
	}

	power, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		config.ComputingPower = 1
	} else {
		config.ComputingPower = power
	}

	return config
}

type Agent struct {
	config *Config
}

func New() *Agent {
	return &Agent{
		config: ConfigFromEnv(),
	}
}

func (a *Agent) Run() error {
	addr := fmt.Sprintf("%s:%s", "orchestrator", a.config.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}

	defer conn.Close()

	grpcClient := pb.NewTaskServiceClient(conn)

	for i := 0; i < a.config.ComputingPower; i++ {
		go worker(grpcClient)
	}

	select {}
}

func worker(client pb.TaskServiceClient) {
	for {
		time.Sleep(100 * time.Millisecond)

		task, err := client.GetTask(context.TODO(), nil)
		if err != nil {
			continue
		}

		code := 200

		result, err := Calc(task.Task)
		if err != nil {
			code = 500
		}

		client.SubmitTask(context.TODO(), &pb.TaskRequest{
			Id:         task.Task.Id,
			Result:     result,
			StatusCode: int32(code),
		})
	}
}

func Calc(t *pb.Task) (float64, error) {
	time.Sleep(time.Duration(t.OperationTime) * time.Millisecond)

	switch t.Operation {
	case '+':
		return t.Arg1 + t.Arg2, nil
	case '-':
		return t.Arg1 - t.Arg2, nil
	case '*':
		return t.Arg1 * t.Arg2, nil
	case '/':
		if t.Arg2 == 0 {
			return 0, schemas.ErrorDivisionByZero
		}
		return t.Arg1 / t.Arg2, nil
	default:
		return 0, schemas.ErrorDivisionByZero
	}
}
