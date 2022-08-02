package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/percona/sqlparsertest/parsers/blastrainmysql"
	"github.com/percona/sqlparsertest/parsers/pgquery"
	"github.com/percona/sqlparsertest/parsers/vitessmysql"
)

const resultsDir = "results"

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
	case "blastrainmysql":
		parse = blastrainmysql.Parse
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

	totalCount := len(queries)
	successCount := len(success)
	errorsCount := len(errors)
	title("final results")
	fmt.Printf("Total queries: %d\n", totalCount)
	fmt.Printf("OK queries: %d\n", successCount)
	fmt.Printf("Queries with error: %d\n", errorsCount)
	devider()

	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Fail: %s", err)
		return
	}

	inputFile := strings.Split(filepath.Base(cfg.InputFile), ".")[0]
	fileSuccess := filepath.Join(workDir, resultsDir, fmt.Sprintf("%s_%s_success.txt", inputFile, cfg.Parser))
	fileErrors := filepath.Join(workDir, resultsDir, fmt.Sprintf("%s_%s_errors.txt", inputFile, cfg.Parser))
	headerText := "Total: %d\nSuccess: %d\nErrors: %d\nInput file: %s\nParser: %s\nComment: %s\n"
	header := fmt.Sprintf(headerText, totalCount, successCount, errorsCount, cfg.InputFile, cfg.Parser, cfg.Comment)

	success = append([]string{header}, success...)
	errors = append([]string{header}, errors...)
	saveToFile(fileSuccess, success)
	saveToFile(fileErrors, errors)
}
