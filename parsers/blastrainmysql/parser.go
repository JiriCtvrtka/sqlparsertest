package blastrainmysql

import (
	"fmt"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

func Parse(q string) (string, error) {
	stmt, err := sqlparser.Parse(q)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("stmt = %+v\n", stmt), nil
}
