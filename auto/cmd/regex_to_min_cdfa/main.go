package main

import (
	"errors"
	"fmt"
	"log"

	"mipt_formal/pkg/tools"
)

func main() {
	regExpr, alphabet, err := getInput()
	if err != nil {
		log.Fatal(err)
	}

	outputDOA := tools.RegexToMinCDFA(regExpr, alphabet)
	fmt.Println(outputDOA)
}

func getInput() (re string, alpha string, err error) {
	var bytesRead int
	fmt.Print("[INPUT] Enter alphabet: ")
	bytesRead, err = fmt.Scanf("%s\n", &alpha)
	if err == nil && bytesRead == 0 {
		err = errors.New("empty alphabet")
	}
	if err != nil {
		return "", "", err
	}

	fmt.Print("[INPUT] Enter regular expression: ")
	_, err = fmt.Scanf("%s\n", &re)
	if err != nil {
		return "", "", err
	}

	return
}
