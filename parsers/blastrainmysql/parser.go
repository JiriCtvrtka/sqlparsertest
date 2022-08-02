package blastrainmysql

import (
	"encoding/json"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
	"github.com/davecgh/go-spew/spew"
)

type output struct {
	Query  string
	Output string
}

func Parse(q string) (string, error) {
	stmt, err := sqlparser.Parse(q)
	if err != nil {
		return "", err
	}

	o := output{
		Query:  q,
		Output: spew.Sdump(stmt),
	}

	res, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", err
	}

	return string(res), nil
}
