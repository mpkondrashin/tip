package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	ErrUnknownType     = errors.New("unknown type")
	ErrUnsupportedType = errors.New("unsupported type")
)

func parseKind(s string) (ListKind, error) {
	switch s {
	case "email-src", "ip-dst", "btc", "hostname", "domain", "filename":
		return 0, fmt.Errorf("%s: %w", s, ErrUnsupportedType)
	case "md5", "sha1", "sha256":
		return Files, nil
	case "url":
		return URLs, nil
	default:
		return 0, ErrUnknownType
	}
}

func csvRead(fileName string) (*TipList, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	//r.Comma = ';'
	//r.Comment = '#'
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	stat := make(map[string]int)
	tip := NewTipList()
	for i, line := range records {
		if i == 0 {
			continue
		}
		if len(line) < 4 {
			log.Fatalf("Can not parse line: %s", strings.Join(line, ","))
		}
		kindStr := line[2]
		stat[kindStr] += 1
		kind, err := parseKind(kindStr)
		if err != nil {
			if errors.Is(err, ErrUnsupportedType) {
				continue
			}
		}
		value := line[3]
		tip.Add(kind, value)
	}
	for key, value := range stat {
		log.Printf("%s: %d", key, value)
	}
	return tip, nil
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s infilename outfilename", os.Args[0])
	}
	log.Print("TIP started")
	tip, err := csvRead(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if err := tip.GenerateResult(os.Args[2]); err != nil {
		log.Fatal(err)
	}
}
