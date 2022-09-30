# Regular expression and state machine Go library

## Formats

### Regular expression

|  operator  |      meaning      |
|:----------:|:-----------------:|
| a &#124; b |    or (union)     |
|     a+     | one or more times |
|     a*     |     any times     |
|    (a)     |   just brackets   |

### State machine in DOA

## Programs

### Regular expression to complete deterministic state machine compiler
- `main.go`: `cmd/r2a/main.go`
- command line argument: alphabet for completion (example: `go run main.go "abcd"`)
- stdin: regular expression