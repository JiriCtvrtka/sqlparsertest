package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
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

func printDevider() {
	fmt.Println("====================================")
}
