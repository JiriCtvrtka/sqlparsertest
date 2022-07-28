package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/kr/pretty"
	"vitess.io/vitess/go/vt/proto/query"
	"vitess.io/vitess/go/vt/sqlparser"
)

func getQueries() ([]string, error) {
	var res []string

	file, err := os.Open("queries.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line[0:2] == "//" {
			continue
		}
		res = append(res, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func addCommentToRandomWord(str, comment string) string {
	var res []string

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 1
	words := strings.Fields(str)
	wordsCount := len(words)
	if wordsCount >= 2 {
		max = wordsCount - 2
	}

	random := rand.Intn(max-min+1) + min
	for index, word := range words {
		if random == index {
			if wordsCount == 1 {
				res = append(res, comment, word)
				continue
			}

			res = append(res, word, comment)
			continue
		}
		res = append(res, word)
	}

	return strings.Join(res, " ")
}

func getQueriesWithComment(comment string) ([]string, error) {
	queries, err := getQueries()
	if err != nil {
		return nil, err
	}
	for index, q := range queries {
		queries[index] = addCommentToRandomWord(q, fmt.Sprintf("/* %s */", comment))
	}

	return queries, nil
}

func getQueriesWithSimpleComment(comment string) ([]string, error) {
	queries, err := getQueries()
	if err != nil {
		return nil, err
	}
	for index, q := range queries {
		queries[index] = addCommentToRandomWord(q, fmt.Sprintf("--%s", comment))
	}

	return queries, nil
}

func main() {
	queries, err := getQueries()
	if err != nil {
		fmt.Errorf("fail %s", err)
		return
	}

	for _, q := range queries {
		normalizeQuery, bindVars, err := sqlparser.Parse2(q)
		if err != nil {
			fmt.Printf("%q", err)
			continue
		}
		//pretty.Println(normalizeQuery)
		pretty.Println(bindVars)
		pretty.Println(GetLiterals(normalizeQuery))
		bv := make(map[string]*query.BindVariable)
		err = sqlparser.Normalize(normalizeQuery, sqlparser.NewReservedVars("", bindVars), bv)
		if err != nil {
			fmt.Printf("%q", err)
		}
		//pretty.Println(normalizeQuery)
		pretty.Println(sqlparser.GetBindvars(normalizeQuery))
		pretty.Println(GetTables(normalizeQuery))
		pretty.Println(sqlparser.String(normalizeQuery))
		q := sqlparser.NewParsedQuery(normalizeQuery)
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
