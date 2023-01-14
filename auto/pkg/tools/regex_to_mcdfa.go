package tools

import (
    "mipt_formal/auto/internal/common"
    "mipt_formal/auto/internal/doa"
    "mipt_formal/auto/internal/modify"
    "mipt_formal/auto/internal/pipeline"
    "mipt_formal/auto/internal/regex"
)

func RegexToMinCDFA(expr string, alpha string) string {
    var result stringWriter
    pipeline.New(
        common.NewLogger(),
        readString(expr),
        regex.NewCompiler(),
        []modify.Step{
            {
                Name: "Building minimal DFA",
                Func: modify.Minimize,
            },
            {
                Name: "Complete DFA",
                Func: modify.Complete(alpha),
            },
        },
        doa.NewWriter(&result),
    )()
    return result.String()
}
