package parsers

type Parser interface {
    Check(word string) bool
}
