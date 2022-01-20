package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

type CovidTestCase struct {
	Threads  int    `json:"thread"`
	Zipcode  string `json:"zipcode"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Expected string `json:"expected"`
}

func GetTests() []CovidTestCase {

	file, err := os.Open("tests.json")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open results.json")
		os.Exit(1)
	}

	dec := json.NewDecoder(file)
	tests := make([]CovidTestCase, 100)

	for {
		if err := dec.Decode(&tests); err != nil {
			break
		}
		break
	}
	if err := file.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not close file: tests.json\n")
		os.Exit(1)
	}
	return tests
}
func TestCovid(t *testing.T) {

	tests := GetTests()
	for num, test := range tests {
		testname := fmt.Sprintf("T=%v", num)
		t.Run(testname, func(t *testing.T) {
			var err error
			cmd := exec.Command("go", "run", ".", "program/covid", fmt.Sprint(test.Threads), fmt.Sprint(test.Zipcode), fmt.Sprint(test.Month), fmt.Sprint(test.Year))
			// cmd := exec.Command("go", "run", "program/bounded.go", "program/worker.go", "program/covid.go", fmt.Sprint(test.Threads), fmt.Sprint(test.Zipcode), fmt.Sprint(test.Month), fmt.Sprint(test.Year))
			out, err := cmd.Output()
			sOut := strings.TrimSpace(string(out))

			if err != nil || test.Expected != sOut {
				t.Errorf("\nRan:%s\nExpected:%s\nGot:%s", cmd, test.Expected, sOut)
			}
		})
	}
}
