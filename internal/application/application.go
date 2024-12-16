package application

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/MaxGolubev19/GoCalculator/pkg/calc"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", a.CalcHandler)
	return http.ListenAndServe(":"+a.config.Addr, nil)
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (a *Application) CalcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request Request
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		error := ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(error)
		return
	}

	expr := strings.TrimSpace(request.Expression)
	res, err := calc.Calc(expr)
	if errors.Is(err, calc.ErrorIncorrectExpression) {
		log.Println(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		error := ErrorResponse{"Expression is not valid"}
		json.NewEncoder(w).Encode(error)
		return
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		error := ErrorResponse{"Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}

	response := Response{Result: res}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		error := ErrorResponse{"Internal server error"}
		json.NewEncoder(w).Encode(error)
		return
	}
}
