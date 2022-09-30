package main

import (
    "errors"
    "fmt"
    "log"

    "mipt_formal/pkg/tools"
)

func main() {
    re, alpha, err := getInput()
    if err != nil {
        log.Fatal(err)
    }

    a := tools.RegexToMinCDFA(re, alpha)
    fmt.Println(a)
}

func getInput() (re string, alpha string, err error) {
    var n int
    fmt.Print("[INPUT] Enter alphabet: ")
    n, err = fmt.Scanf("%s\n", &alpha)
    if err == nil && n == 0 {
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
