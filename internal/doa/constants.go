package doa

const (
    Version          = "DOA: v1\n"
    StartFormat      = "Start: %v\n"
    AcceptanceFormat = "Acceptance: %v\n"

    StateFormat = "State: %v\n"
    EdgeFormat  = "    -> %v %v\n"

    Begin = "--BEGIN--"
    End   = "--END--"

    StateConj = " & "
    Epsilon   = "EPS"

    MinimalLength = len(Version) + len(StartFormat) + len(AcceptanceFormat) + len(Begin) + len(End)
)
