// package customerimporter reads from the given customers.csv file and returns a
// sorted (data structure of your choice) of email domains along with the number
// of customers with e-mail addresses for each domain.  Any errors should be
// logged (or handled). Performance matters (this is only ~3k lines, but *could*
// be 1m lines or run on a small machine).
package main

import (
	"bufio"
	"log"
	"net/mail"
	"os"
	"sort"
	"strings"
)

type result struct {
	domain string
	count  int
}

func main() {
	file, fileOpenError := os.Open("customers.csv")
	if fileOpenError != nil {
		log.Fatal(fileOpenError) // TODO custom message?
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	domainCount := make(map[string]int)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		email, invalidEmailError := mail.ParseAddress(strings.Split(line, ",")[2])

		if invalidEmailError != nil {
			log.Println("Line", i, "could not be parsed")
			log.Print(invalidEmailError)
		} else {

			domain := strings.Split(email.Address, "@")[1]

			_, domainAlreadyAdded := domainCount[domain]
			if domainAlreadyAdded {
				domainCount[domain] = domainCount[domain] + 1
			} else {
				domainCount[domain] = 1
			}

		}

		i++
	}

	if scannerError := scanner.Err(); scannerError != nil {
		log.Fatal(scannerError) // TODO custom message?
	}

	domains := make([]string, 0, len(domainCount))
	for k := range domainCount {
		domains = append(domains, k)
	}
	sort.Strings(domains)

	results := make([]result, 0, len(domains))
	for _, k := range domains {
		results = append(results, result{
			domain: k,
			count:  domainCount[k],
		})
	}

	// fmt.Println(results)
	// return results
}
