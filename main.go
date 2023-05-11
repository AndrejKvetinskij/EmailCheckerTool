package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain,hasMX,hasSPF,sprRecord,hasDMARC,dmarcRecord\n")

	for scanner.Scan() {
		email := scanner.Text()
		coll := strings.Split(email, "@")
		if len(coll) == 2 {
			checkDomain(email, coll[1])
		} else {
			fmt.Printf("Email: %v is not correct", email)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error: could not read from intput: %v\n", err)
	}
}

func checkDomain(email, domain string) {
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)

	if err != nil {
		log.Printf("Error:%v\n", err)
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			spfRecord = record
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("ErrorL%v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			dmarcRecord = record
			break
		}
	}

	if hasMX || hasSPF || hasDMARC {
		fmt.Printf("Email: %v is correct\n", email)
	} else {
		fmt.Printf("Email: %v is not correct\n", email)
	}
	fmt.Printf("%v %v %v %v %v %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}
