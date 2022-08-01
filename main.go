package main

import (
	"fmt"
	"os"

	"github.com/percona/sqlparsertest/parsers/pgquery"
	"github.com/percona/sqlparsertest/parsers/vitessmysql"
)

func main() {
	if len(os.Args) == 1 {
		help()
		return
	}

	var errors, success, queries []string
	var err error

	cfg := proceedFlags(os.Args)
	switch cfg.Comment {
	case "":
		queries, err = getQueries(cfg.InputFile)
	default:
		queries, err = getQueriesWithComment(cfg.InputFile, cfg.Comment)
	}
	if err != nil {
		fmt.Printf("Fail: %s", err)
		return
	}

	var parse func(q string) (string, error)
	switch cfg.Parser {
	case "vitessmysql":
		parse = vitessmysql.Parse
	case "pgquery":
		parse = pgquery.Parse
	default:
		fmt.Println(cfg.Parser)
		os.Exit(1)
	}

	for _, q := range queries {
		res, err := parse(q)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Original query: %s, error: %s\n", q, err.Error()))
			continue
		}

		success = append(success, res)
	}

	title("final results")
	fmt.Printf("Total queries: %d\n", len(queries))
	fmt.Printf("Queries with error: %d\n", len(errors))
	fmt.Printf("OK queries: %d\n", len(success))
	devider()

	header := fmt.Sprintf("Collected with values: Input file: %s, Parser: %s, Comment: %s\n", cfg.InputFile, cfg.Parser, cfg.Comment)
	success = append([]string{header}, success...)
	errors = append([]string{header}, errors...)
	saveToFile("success.txt", success)
	saveToFile("errors.txt", errors)
}
