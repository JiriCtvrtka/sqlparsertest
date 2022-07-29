package main

import (
	"encoding/json"
	"fmt"

	"vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

func main() {
	var errors []string
	var success []string

	queries, err := getQueries()
	if err != nil {
		fmt.Printf("fail %s", err)
		return
	}

	for _, q := range queries {
		normalizedQuery, bindVars, err := sqlparser.Parse2(q)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		bv := make(map[string]*query.BindVariable)
		err = sqlparser.Normalize(normalizedQuery, sqlparser.NewReservedVars("", bindVars), bv)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}
		parsedQuery := sqlparser.NewParsedQuery(normalizedQuery)

		normalizedQueryJSON, err := json.Marshal(normalizedQuery)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		parsedQueryJSON, err := json.Marshal(parsedQuery)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		bindVarsJSON, err := json.Marshal(sqlparser.GetBindvars(normalizedQuery))
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		literalsJSON, err := json.Marshal(GetLiterals(normalizedQuery))
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		tablesJSON, err := json.Marshal(GetTables(normalizedQuery))
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		comment, _ := sqlparser.ExtractMysqlComment(q)
		commentsJSON, err := json.Marshal(comment)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		o := output{
			Query:           fmt.Sprintf("%+s\n", q),
			NormalizedQuery: fmt.Sprintf("%+s\n", string(normalizedQueryJSON)),
			BindVars:        fmt.Sprintf("%+s\n", string(bindVarsJSON)),
			Literals:        fmt.Sprintf("%+s\n", string(literalsJSON)),
			Tables:          fmt.Sprintf("%+s\n", string(tablesJSON)),
			Comments:        fmt.Sprintf("%+s\n", string(commentsJSON)),
			Parsed:          fmt.Sprintf("%+s\n", string(parsedQueryJSON)),
		}
		success = append(success, fmt.Sprintf("%+v\n", o))
	}

	title("final results")
	fmt.Printf("Total queries: %d\n", len(queries))
	fmt.Printf("Queries with error: %d\n", len(errors))
	fmt.Printf("OK queries: %d\n", len(success))
	devider()

	saveToFile("errors.txt", errors)
	saveToFile("success.txt", success)
}

// GetLiterals returns a map of the bind vars referenced in the statement.
func GetLiterals(stmt sqlparser.Statement) []string {
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

// GetTables returns a map of the bind vars referenced in the statement.
func GetTables(stmt sqlparser.Statement) []string {
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
