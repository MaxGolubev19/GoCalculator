package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

type Config struct {
	Port           string
	ComputingPower int
}

func ConfigFromEnv() *Config {
	config := new(Config)

	config.Port = os.Getenv("PORT")
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
	for i := 0; i < a.config.ComputingPower; i++ {
		go worker("http://orchestrator:" + a.config.Port + "/internal/task")
	}

	select {}
}

func worker(url string) {
	for {
		time.Sleep(100 * time.Millisecond)

		task, err := get(url)
		if err != nil {
			continue
		}

		result, err := calc(task)
		if err != nil {
			continue
		}

		err = post(url, task.Id, result)
		if err != nil {
			continue
		}
	}
}

func get(url string) (*schemas.Task, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("404: задач нет")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("500: ошибка сервера")
	}

	var tr schemas.Task
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	return &tr, nil
}

func calc(t *schemas.Task) (float64, error) {
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
			return 0, errors.New("division by zero")
		}
		return t.Arg1 / t.Arg2, nil
	default:
		return 0, errors.New("unknown operation")
	}
}

func post(url string, id int, result float64) error {
	tr := schemas.TaskRequest{
		Id:     id,
		Result: result,
	}

	trJson, err := json.Marshal(tr)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(trJson))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("404: not found")
	}
	if resp.StatusCode == http.StatusUnprocessableEntity {
		return errors.New("422: inbavid data")
	}
	return errors.New("unknown error")
}
