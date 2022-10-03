package modify

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "mipt_formal/internal/common"
	"mipt_formal/internal/nfa"
)

func TestMinimize(t *testing.T) {
	type testcase struct {
		name string
		init *nfa.Machine
		want *nfa.Machine
	}

	cases := [...]testcase{
		{
			name: "just works",
			init: nfa.New([]State{0}, []State{0, 1}, []Transition{
				{From: 0, To: 1, By: "a"},
				{From: 1, To: 1, By: "a"},
			}),
			want: nfa.New([]State{0}, []State{0}, []Transition{
				{From: 0, To: 0, By: "a"},
			}),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			Minimize(c.init)
			assert.True(t, c.init.Equal(c.want), c.init.DOA())
		})
	}
}
