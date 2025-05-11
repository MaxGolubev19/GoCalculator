package orchestrator

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MaxGolubev19/GoCalculator/pkg/parse"
	"github.com/MaxGolubev19/GoCalculator/pkg/schemas"
)

func (o *Orchestrator) ExpressonsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	login, ok := r.Context().Value(schemas.UserContextKey).(string)
	if !ok {
		http.Error(w, "User not authorized", http.StatusUnauthorized)
		return
	}

	expressions, err := o.GetExpressions(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := schemas.ExpressionsResponse{
		Expressions: expressions,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (o *Orchestrator) GetExpressions(login string) ([]schemas.Expression, error) {
	const query = `
	SELECT id, expression, status, result
	FROM expressions
	WHERE user_login = ?;
	`

	rows, err := o.db.Query(query, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []schemas.Expression
	for rows.Next() {
		var expr schemas.Expression
		if err := rows.Scan(&expr.Id, &expr.Expression, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expressions, nil
}

func (o *Orchestrator) ExpressonByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/api/v1/expressions/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	login, ok := r.Context().Value(schemas.UserContextKey).(string)
	if !ok {
		http.Error(w, "User not authorized", http.StatusUnauthorized)
		return
	}

	expression, err := o.GetExpressionById(id, login)
	if err == schemas.ErrorNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(expression); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (o *Orchestrator) GetExpressionById(id int, login string) (*schemas.Expression, error) {
	const query = `
	SELECT id, expression, status, result
	FROM expressions
	WHERE user_login = ? AND id = ?;
	`

	var expr schemas.Expression
	if err := o.db.QueryRow(query, login, id).Scan(&expr.Id, &expr.Expression, &expr.Status, &expr.Result); err != nil {
		if err == sql.ErrNoRows {
			return nil, schemas.ErrorNotFound
		}
		return nil, err
	}

	return &expr, nil
}

func (o *Orchestrator) AddExpression(expression string, login string) (int, error) {
	const query = `
	INSERT INTO expressions (expression, user_login, status, result)
	VALUES (?, ?, ?, ?);
	`
	resultExec, err := o.db.Exec(query, expression, login, schemas.IN_PROGRESS, 0)
	if err != nil {
		return 0, err
	}

	id, err := resultExec.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = o.ParseExpression(int(id), expression)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (o *Orchestrator) ParseExpression(id int, expression string) error {
	actions, err := parse.New(expression)
	if err != nil {
		return err
	}

	go o.worker(int(id), actions)
	return nil
}

func (o *Orchestrator) SetExpressionError(id int) error {
	const query = `
	UPDATE expressions
	SET status = ?
	WHERE id = ?;
	`
	_, err := o.db.Exec(query, schemas.ERROR, id)
	if err != nil {
		return err
	}

	return nil
}

func (o *Orchestrator) SetExpressionDone(id int, result float64) error {
	const query = `
	UPDATE expressions
	SET result = ?, status = ?
	WHERE id = ?;
	`
	_, err := o.db.Exec(query, result, schemas.DONE, id)
	if err != nil {
		return err
	}

	return nil
}

func (o *Orchestrator) GetExpressionsInProgress() ([]schemas.Expression, error) {
	const query = `
	SELECT id, expression, user_login, status, result
	FROM expressions
	WHERE status = ?;
	`

	rows, err := o.db.Query(query, schemas.IN_PROGRESS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expressions []schemas.Expression
	for rows.Next() {
		var expr schemas.Expression
		if err := rows.Scan(&expr.Id, &expr.Expression, &expr.Status, &expr.Result); err != nil {
			return nil, err
		}
		expressions = append(expressions, expr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expressions, nil
}
