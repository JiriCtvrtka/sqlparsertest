package pgquery

import (
	"encoding/json"

	pg_query "github.com/pganalyze/pg_query_go"
)

type output struct {
	Query  string
	Output pg_query.ParsetreeList
}

func Parse(q string) (string, error) {
	o, err := pg_query.Parse(q)
	if err != nil {
		return "", err
	}

	res, err := json.MarshalIndent(output{Query: q, Output: o}, "", "  ")
	if err != nil {
		return "", err
	}

	return string(res), nil
}
