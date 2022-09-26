package main

import (
    "fmt"
    "log"
    "mipt_formal/internal/regex"
)

func main() {
    var re string
    _, _ = fmt.Scanf("%s", &re)
    nfa, err := regex.NewCompiler().Compile(re)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(nfa.DOA())
}
