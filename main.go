package main

import (
	"fmt"

	"github.com/kr/pretty"
	"vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

func main() {
	queries, err := getQueries()
	if err != nil {
		fmt.Printf("fail %s", err)
		return
	}

	printDevider()
	for _, q := range queries {
		normalizedQuery, bindVars, err := sqlparser.Parse2(q)
		if err != nil {
			fmt.Printf("%s", err)
			continue
		}
		fmt.Println("Original query:")
		fmt.Println(q)

		fmt.Println("Normalized query:")
		pretty.Println(sqlparser.String(normalizedQuery))
		printDevider()
		fmt.Println("Bind vars:")
		pretty.Println(bindVars)
		printDevider()
		fmt.Println("Literals:")
		pretty.Println(GetLiterals(normalizedQuery))
		printDevider()
		bv := make(map[string]*query.BindVariable)
		err = sqlparser.Normalize(normalizedQuery, sqlparser.NewReservedVars("", bindVars), bv)
		if err != nil {
			fmt.Printf("%s", err)
		}
		printDevider()
		fmt.Println("Bind vars:")
		pretty.Println(sqlparser.GetBindvars(normalizedQuery))
		printDevider()
		fmt.Println("Tables:")
		pretty.Println(GetTables(normalizedQuery))
		printDevider()
		parsedQuery := sqlparser.NewParsedQuery(normalizedQuery)
		pretty.Println(parsedQuery.GenerateQuery(bv, map[string]sqlparser.Encodable{}))
		pretty.Println(sqlparser.ExtractMysqlComment(q))
		printDevider()
	}
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
