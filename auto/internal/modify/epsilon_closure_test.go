package modify

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "mipt_formal/internal/common"
	"mipt_formal/internal/nfa"
)

func TestEliminateEpsilonMoves(t *testing.T) {
	type testcase struct {
		name string
		got  *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "just works",
			got: nfa.NewMachine([]State{0}, []State{0}, []Transition{
				{From: 0, To: 0, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{0}, []Transition{}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			EliminateEpsilonMoves(c.got)
			assert.True(t, c.got.Equal(c.want))
		})
	}
}

func Test_buildTransitiveEpsilonClosure(t *testing.T) {
	type testcase struct {
		name string
		got  *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "transitive",
			got: nfa.NewMachine([]State{0}, []State{2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: Epsilon},
				{From: 0, To: 2, By: Epsilon},
			}),
		},
		{
			name: "loops",
			got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 0, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 0, By: Epsilon},
				{From: 0, To: 0, By: Epsilon},
				{From: 1, To: 1, By: Epsilon},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			buildTransitiveEpsilonClosure(c.got)
			assert.True(t, c.got.Equal(c.want))
		})
	}
}

func Test_compressAcceptances(t *testing.T) {
	type testcase struct {
		name string
		got  *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "basic",
			got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
				{From: 0, To: 1, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{0, 1}, []Transition{
				{From: 0, To: 1, By: Epsilon},
			}),
		},
		{
			name: "transitive closure",
			got: nfa.NewMachine([]State{0}, []State{2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: Epsilon},
				{From: 0, To: 2, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{0, 1, 2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: Epsilon},
				{From: 0, To: 2, By: Epsilon},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			compressAcceptances(c.got)
			assert.True(t, c.got.Equal(c.want))
		})
	}
}

func Test_addTransitiveEdges(t *testing.T) {
	type testcase struct {
		name string
		got  *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "transitive",
			got: nfa.NewMachine([]State{0}, []State{2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: "a"},
			}),
			want: nfa.NewMachine([]State{0}, []State{2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 2, By: "a"},
				{From: 0, To: 2, By: "a"},
			}),
		},
		{
			name: "loops",
			got: nfa.NewMachine([]State{0}, []State{1}, []Transition{
				{From: 0, To: 1, By: "b"},
				{From: 1, To: 0, By: Epsilon},
				{From: 0, To: 0, By: "a"},
			}),
			want: nfa.NewMachine([]State{0}, []State{1}, []Transition{
				{From: 0, To: 1, By: "b"},
				{From: 1, To: 0, By: Epsilon},
				{From: 0, To: 0, By: "a"},
				{From: 1, To: 0, By: "a"},
				{From: 1, To: 1, By: "b"},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			addTransitiveEdges(c.got)
			assert.True(t, c.got.Equal(c.want))
		})
	}
}

func Test_removeEpsilonMoves(t *testing.T) {
	type testcase struct {
		name string
		got  *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "regular case",
			got: nfa.NewMachine([]State{0}, []State{0, 1, 2}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 0, To: 2, By: "b"},
				{From: 2, To: 1, By: Epsilon},
				{From: 2, To: 3, By: "c"},
				{From: 3, To: 1, By: "a"},
			}),
			want: nfa.NewMachine([]State{0}, []State{0, 1, 2}, []Transition{
				{From: 0, To: 2, By: "b"},
				{From: 2, To: 3, By: "c"},
				{From: 3, To: 1, By: "a"},
			}),
		},
		{
			name: "remove all edges",
			got: nfa.NewMachine([]State{0}, []State{0, 1}, []Transition{
				{From: 0, To: 1, By: Epsilon},
				{From: 1, To: 1, By: Epsilon},
			}),
			want: nfa.NewMachine([]State{0}, []State{0, 1}, []Transition{}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			removeEpsilonMoves(c.got)
			assert.True(t, c.got.Equal(c.want))
		})
	}
}
