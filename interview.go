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
	// TODO slice for return type
	freqCountMap := make(map[string]int)

	file, error := os.Open("customers.csv")

	if error != nil {
		log.Fatal(error)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		email := strings.Split(line, ",")[2]

		var domain string
		emailParts := strings.Split(email, "@")
		if len(emailParts) == 2 {
			domain = emailParts[1]
			fmt.Println(domain)

			_, domainAlreadyAdded := freqCountMap[domain]
			if domainAlreadyAdded {
				freqCountMap[domain] = freqCountMap[domain] + 1
			} else {
				freqCountMap[domain] = 1
			}

		} else {
			// TODO error handling
			// domain = ""
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
