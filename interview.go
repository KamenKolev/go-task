// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {

	file, error := os.Open("customers.csv")

	if error != nil {
		log.Fatal(error)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		// fmt.Println(scanner.Text())
		line := scanner.Text()
		// idk := strings.SplitAfterN(",", line, 1)
		idk := strings.Split(line, ",")
		email := idk[2]
		fmt.Println(email)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
