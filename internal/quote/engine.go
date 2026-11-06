package quote

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/expr-lang/expr"
	"golang.org/x/sync/singleflight"
)

type Table struct {
	ID          string `json:"id"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Formula     string `json:"formula"`
}

type Request struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	WeightKg    float64 `json:"weight_kg"`
}

type Option struct {
	TableID string  `json:"table_id"`
	Amount  float64 `json:"amount"`
	Formula string  `json:"formula"`
}

type Response struct {
	Origin      string   `json:"origin"`
	Destination string   `json:"destination"`
	WeightKg    float64  `json:"weight_kg"`
	Currency    string   `json:"currency"`
	Quotes      []Option `json:"quotes"`
}

type Engine struct {
	path string
	ttl  time.Duration

	mu    sync.RWMutex
	tables []Table
	expiry time.Time
	sf     singleflight.Group
}

func New(path string, ttl time.Duration) *Engine {
	return &Engine{path: path, ttl: ttl}
}

func (e *Engine) Quote(req Request) (Response, error) {
	if req.Origin == "" || req.Destination == "" {
		return Response{}, fmt.Errorf("origin and destination are required")
	}
	if req.WeightKg <= 0 {
		return Response{}, fmt.Errorf("weight_kg must be greater than zero")
	}

	tables, err := e.load()
	if err != nil {
		return Response{}, err
	}

	out := Response{
		Origin:      req.Origin,
		Destination: req.Destination,
		WeightKg:    req.WeightKg,
		Currency:    "BRL",
		Quotes:      []Option{},
	}

	env := map[string]any{"weight": req.WeightKg}
	for _, t := range tables {
		if t.Origin != req.Origin || t.Destination != req.Destination {
			continue
		}
		amount, err := eval(t.Formula, env)
		if err != nil {
			return Response{}, fmt.Errorf("table %s: %w", t.ID, err)
		}
		out.Quotes = append(out.Quotes, Option{
			TableID: t.ID,
			Amount:  amount,
			Formula: t.Formula,
		})
	}

	if len(out.Quotes) == 0 {
		return Response{}, fmt.Errorf("no rate table for %s → %s", req.Origin, req.Destination)
	}
	return out, nil
}

func (e *Engine) load() ([]Table, error) {
	e.mu.RLock()
	if e.tables != nil && time.Now().Before(e.expiry) {
		tables := e.tables
		e.mu.RUnlock()
		return tables, nil
	}
	e.mu.RUnlock()

	v, err, _ := e.sf.Do("tables", func() (any, error) {
		data, err := os.ReadFile(e.path)
		if err != nil {
			return nil, err
		}
		var tables []Table
		if err := json.Unmarshal(data, &tables); err != nil {
			return nil, err
		}
		e.mu.Lock()
		e.tables = tables
		e.expiry = time.Now().Add(e.ttl)
		e.mu.Unlock()
		return tables, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Table), nil
}

func eval(formula string, env map[string]any) (float64, error) {
	prog, err := expr.Compile(formula, expr.Env(env), expr.AsFloat64())
	if err != nil {
		return 0, err
	}
	v, err := expr.Run(prog, env)
	if err != nil {
		return 0, err
	}
	n, ok := v.(float64)
	if !ok {
		return 0, fmt.Errorf("formula did not return a number")
	}
	return n, nil
}
