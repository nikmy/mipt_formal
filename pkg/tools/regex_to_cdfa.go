package tools

import (
	"mipt_formal/internal/common"
	"mipt_formal/internal/doa"
	"mipt_formal/internal/modify"
	"mipt_formal/internal/pipeline"
	"mipt_formal/internal/regex"
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
