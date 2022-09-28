package regex

const (
    kleeneStar = '*'
    orOperator = '|'
    oneOrMore  = '+'
    lBracket   = '('
    rBracket   = ')'

    fewArgumentsErrorFormat = "few arguments for %v operator"
)

var priority = map[byte]int{
    kleeneStar: 10,
    oneOrMore:  10,
    orOperator: 6,
    rBracket:   4,
    lBracket:   2,
}
