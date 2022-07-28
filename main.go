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
		fmt.Errorf("fail %s", err)
		return
	}

	for _, q := range queries {
		normalizedQuery, bindVars, err := sqlparser.Parse2(q)
		if err != nil {
			fmt.Printf("%q", err)
			continue
		}
		fmt.Println("Original query:")
		fmt.Println(q)
		fmt.Println("====================================")
		fmt.Println("Normalized query")
		fmt.Printf("%+v", normalizedQuery.Comments)
		fmt.Println("====================================")
		fmt.Println("Bind vars")
		pretty.Println(bindVars)
		fmt.Println("====================================")
		fmt.Println("Literals")
		pretty.Println(GetLiterals(normalizedQuery))
		fmt.Println("====================================")
		bv := make(map[string]*query.BindVariable)
		err = sqlparser.Normalize(normalizedQuery, sqlparser.NewReservedVars("", bindVars), bv)
		if err != nil {
			fmt.Printf("%q", err)
		}
		//pretty.Println(normalizedQuery)
		fmt.Println("====================================")
		pretty.Println(sqlparser.GetBindvars(normalizedQuery))
		fmt.Println("====================================")
		pretty.Println(GetTables(normalizedQuery))
		fmt.Println("====================================")
		pretty.Println(sqlparser.String(normalizedQuery))
		fmt.Println("====================================")
		q := sqlparser.NewParsedQuery(normalizedQuery)
		pretty.Println(q.GenerateQuery(bv, map[string]sqlparser.Encodable{}))
		//pretty.Println(sqlparser.ExtractMysqlComment(sql))
		fmt.Println("====================================")

		break
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
