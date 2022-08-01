package main

import (
	"fmt"
	"os"

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
	default:
		os.Exit(1)
	}

	header := fmt.Sprintf("Collected with values: Input file: %s, Parser: %s, Comment: %s\n", cfg.InputFile, cfg.Parser, cfg.Comment)

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

	success = append([]string{header}, success...)
	errors = append([]string{header}, errors...)
	saveToFile("errors.txt", errors)
	saveToFile("success.txt", success)
}
