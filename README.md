# Regular expression and state machine Go library

## Formats

### Regular expression

| operator |    meaning    |
|:--------:|:-------------:|
|  a + b   |  or (union)   |
|    a*    |   any times   |
|   (a)    | just brackets |
|    1     |  empty word   |

### State machine in DOA

## Programs

### Regular expression to complete deterministic state machine compiler
- `main.go`: `cmd/regex_to_cdfa/main.go`

### Regular expression to minimal complete deterministic state machine compiler
- `main.go`: `regex_to_mcdfa/main.go`