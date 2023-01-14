package tools

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/doa"
    "mipt_formal/auto/internal/modify"
    "mipt_formal/auto/internal/pipeline"
    "mipt_formal/auto/internal/regex"
)

func RegexToCDFA(expr string, alpha string) string {
    var result stringWriter
    pipeline.New(
        common.NewLogger(),
        readString(expr),
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
                Name: "Building equivalent DFA",
                Func: modify.Determine,
            },
            {
                Name: "Removing unused states",
                Func: modify.RemoveUnusedStates,
            },
            {
                Name: "Completing DFA",
                Func: modify.Complete(alpha),
            },
        },
        doa.NewWriter(&result),
    )()
    return result.String()
}
