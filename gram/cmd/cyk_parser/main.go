package main

import (
    "bufio"
    "fmt"
    "io"
    "log"
    "mipt_formal/gram/pkg/parsers"
    "os"
    "strings"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: ./cyk_parser <grammar_file> <words_file>")
        fmt.Println()
    }

    rules, err := readLines(os.Args[1])
    if err != nil {
        log.Fatal(err)
    }

    words, err := readLines(os.Args[2])
    if err != nil {
        log.Fatal(err)
    }

    out, err := os.Create("output.cyk")
    if err != nil {
        log.Fatal(err)
    }

    parser, err := parsers.NewCYKParser(rules)
    if err != nil {
        log.Fatal(err)
    }

    for _, word := range words {
        if parser.Check(word) {
            _, _ = fmt.Fprintln(out, word)
        }
    }

    log.Println("CYK parsing: Success")
}

func readLines(filename string) ([]string, error) {
    grammarFin, err := os.Open(filename)
    if err != nil {
        log.Fatalf("couldn't open grammar file: %s", err)
    }
    gr, input := bufio.NewReader(grammarFin), true
    rules := make([]string, 0, 1)
    for input {
        line, err := gr.ReadString('\n')
        if err == io.EOF {
            input = false
        } else if err != nil {
            return nil, err
        }
        line = strings.TrimSuffix(line, "\n")
        line = strings.TrimSuffix(line, "\r")
        rules = append(rules, line)
    }
    return rules, nil
}
