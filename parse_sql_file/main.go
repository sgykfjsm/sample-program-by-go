package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func parseQueryLine(text string, queries []string) []string {
	if strings.HasPrefix(text, "--") || len(queries) == 0 {
		queries = append(queries, text)
		return queries
	}

	lines := strings.Split(text, ";")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if i == 0 && len(queries) > 0 && !strings.HasSuffix(queries[len(queries)-1], ";") {
			queries[len(queries)-1] += " " + line
		} else {
			queries = append(queries, line)
		}

		if i != len(lines)-1 {
			queries[len(queries)-1] += ";"
		}
	}

	return queries
}

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Need to pass SQL file name")
	}

	fName := args[1]
	f, err := os.Open(fName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	var queries []string

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		queries = parseQueryLine(text, queries)
	}

	for _, q := range queries {
		fmt.Println(q)
	}
}
