package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type config struct {
	InputFile string
	Parser    string
	Comment   string
}

var parsers = []string{"vitessmysql", "blastrainmysql", "pgquery"}

func help() {
	fmt.Println("Please define input file by flag --input PATH/NAME")
	fmt.Printf("Please define parser by flag --parser OPTION, options: %+v\n", parsers)
	fmt.Println("You can add comment to every query at random place by flag --comment VALUE")
}

func checkArgValue(args []string, index int, flag string) {
	if len(args) < index+1 {
		fmt.Printf("Value for flag %s is missing\n", flag)
		os.Exit(1)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func proceedFlags(args []string) *config {
	cfg := new(config)

	for k, v := range args {
		if v[0:2] != "--" {
			continue
		}

		switch v {
		case "--input":
			checkArgValue(args, k+1, v)
			cfg.InputFile = args[k+1]
		case "--parser":
			checkArgValue(args, k+1, v)
			p := args[k+1]
			if !contains(parsers, p) {
				fmt.Printf("Unsupported parser %s, choose one of %+v\n", p, parsers)
				os.Exit(1)
			}

			cfg.Parser = p
		case "--comment":
			checkArgValue(args, k+1, v)
			cfg.Comment = args[k+1]
		default:
		}
	}

	return cfg
}

func getQueries(inputFile string) ([]string, error) {
	var res []string

	file, err := os.Open(inputFile)
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
	min, max := 0, 1
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

func getQueriesWithComment(inputFile, comment string) ([]string, error) {
	queries, err := getQueries(inputFile)
	if err != nil {
		return nil, err
	}
	for index, q := range queries {
		queries[index] = addCommentToRandomWord(q, fmt.Sprintf("/* %s */", comment))
	}

	return queries, nil
}

func getQueriesWithSimpleComment(inputFile, comment string) ([]string, error) {
	queries, err := getQueries(inputFile)
	if err != nil {
		return nil, err
	}
	for index, q := range queries {
		queries[index] = addCommentToRandomWord(q, fmt.Sprintf("--%s", comment))
	}

	return queries, nil
}

func title(name string) {
	devider()
	fmt.Printf("%s\n", strings.ToUpper(name))
	devider()
}

func devider() {
	fmt.Println("====================================")
}

func saveToFile(path string, lines []string) {
	f, err := os.Create(path)
	if err != nil {
		devider()
		fmt.Print(err)
		devider()
		return
	}
	defer f.Close()

	for _, line := range lines {
		_, err := f.WriteString(line + "\n")
		if err != nil {
			devider()
			fmt.Print(err)
			devider()
			return
		}
	}
}
