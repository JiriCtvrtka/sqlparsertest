package vitessmysql

import (
	"encoding/json"
	"fmt"

	"vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

type output struct {
	Query           string
	NormalizedQuery string
	BindVars        string
	Literals        string
	Tables          string
	Comments        string
	Parsed          string
	Values          string
}

func Parse(q string) (string, error) {
	normalizedQuery, bindVars, err := sqlparser.Parse2(q)
	if err != nil {
		return "", err
	}

	bv := make(map[string]*query.BindVariable)
	err = sqlparser.Normalize(normalizedQuery, sqlparser.NewReservedVars("", bindVars), bv)
	if err != nil {
		return "", err
	}
	parsedQuery := sqlparser.NewParsedQuery(normalizedQuery)

	normalizedQueryJSON, err := json.Marshal(normalizedQuery)
	if err != nil {
		return "", err
	}

	parsedQueryJSON, err := json.Marshal(parsedQuery)
	if err != nil {
		return "", err
	}

	bindVarsJSON, err := json.Marshal(sqlparser.GetBindvars(normalizedQuery))
	if err != nil {
		return "", err
	}

	literalsJSON, err := json.Marshal(getLiterals(normalizedQuery))
	if err != nil {
		return "", err
	}

	tablesJSON, err := json.Marshal(getTables(normalizedQuery))
	if err != nil {
		return "", err
	}

	valuesJSON, err := json.Marshal(getAll(normalizedQuery))
	if err != nil {
		return "", err
	}

	comment, _ := sqlparser.ExtractMysqlComment(q)
	commentsJSON, err := json.Marshal(comment)
	if err != nil {
		return "", err
	}

	o := output{
		Query:           q,
		NormalizedQuery: string(normalizedQueryJSON),
		BindVars:        string(bindVarsJSON),
		Literals:        string(literalsJSON),
		Tables:          string(tablesJSON),
		Comments:        string(commentsJSON),
		Parsed:          string(parsedQueryJSON),
		Values:          string(valuesJSON),
	}

	res, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

// getLiterals returns a map of the bind vars referenced in the statement.
func getLiterals(stmt sqlparser.Statement) []string {
	var bindvars []string
	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node := node.(type) {
		case *sqlparser.ColName, sqlparser.TableName:
			// Common node types that never contain expressions but create a lot of object
			// allocations.
			return false, nil
		case *sqlparser.Literal:
			bindvars = append(bindvars, node.Val)
		}
		return true, nil
	}, stmt)

	return bindvars
}

// getTables returns a map of the bind vars referenced in the statement.
func getTables(stmt sqlparser.Statement) []string {
	var results []string
	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch node := node.(type) {
		case *sqlparser.ColName:
			// Common node types that never contain expressions but create a lot of object
			// allocations.
			return false, nil
		case sqlparser.TableName:
			results = append(results, node.Name.String())
			return false, nil
		}
		return true, nil
	}, stmt)

	return results
}

func getAll(stmt sqlparser.Statement) map[string]string {
	results := make(map[string]string)
	_ = sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {

		switch n := node.(type) {
		default:
			results[fmt.Sprintf("%T", n)] = fmt.Sprintf("%+v", n)
		}
		return true, nil
	}, stmt)

	return results
}
