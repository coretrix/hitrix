package graphqlparser

import (
	"bytes"
	"encoding/json"
)

type operationType uint8

const (
	queryOperation operationType = iota
	mutationOperation
)

type QueryParser struct{}

func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

func (c *QueryParser) ParseQuery(q interface{}, variables map[string]interface{}) (bytes.Buffer, error) {
	return c.parse(queryOperation, q, variables)
}

func (c *QueryParser) ParseMutation(m interface{}, variables map[string]interface{}) (bytes.Buffer, error) {
	return c.parse(mutationOperation, m, variables)
}

func (c *QueryParser) parse(op operationType, v interface{}, variables map[string]interface{}) (bytes.Buffer, error) {
	var query string

	switch op {
	case queryOperation:
		query = constructQuery(v, variables)
	case mutationOperation:
		query = constructMutation(v, variables)
	}

	in := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables,omitempty"`
	}{
		Query:     query,
		Variables: variables,
	}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(in); err != nil {
		return buff, err
	}

	return buff, nil
}
