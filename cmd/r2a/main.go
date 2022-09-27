package main

import (
    "mipt_formal/internal/common"
    "mipt_formal/internal/doa"
    "mipt_formal/internal/modify"
    "mipt_formal/internal/pipeline"
    "mipt_formal/internal/regex"
)

func main() {
    pipeline.New(
        common.NewLogger(),
        common.NewStdinReader(),
        regex.NewCompiler(),
        []modify.Step{
            {
                Name: "Removing epsilon moves",
                Func: modify.EliminateEpsilonMoves,
            },
            {
                Name: "Removing unused states",
                Func: modify.RemoveUnusedStates,
            },
            {
                Name: "Building equal DFA",
                Func: modify.Determine,
            },
            {
                Name: "Removing extra states",
                Func: modify.RemoveUnusedStates,
            },
        },
        doa.NewStdoutWriter(),
    )()
}
