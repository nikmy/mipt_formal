package regex

import "mipt_formal/internal/common"

func _or(args ...*fragment) *fragment {
	result := &fragment{
		Start:  newIntrusiveState(common.Epsilon),
		Accept: newIntrusiveState(common.Epsilon),
	}
	for _, f := range args {
		result.Start.precede(f.Start)
		f.Accept.precede(result.Accept)
	}
	return result
}

func _concat(args ...*fragment) *fragment {
	result := args[0]
	for _, f := range args[1:] {
		result.Accept.precede(f.Start)
		result.Accept = f.Accept
	}
	return result
}

func _kleene(args ...*fragment) *fragment {
	f := args[0]
	start := newIntrusiveState(common.Epsilon, f.Start, f.Accept)
	f.Accept.precede(start)
	f.Start = start
	return f
}

func _id(symbol byte) func(args ...*fragment) *fragment {
	return func(...*fragment) *fragment {
		accept := newIntrusiveState(common.Epsilon)
		start := newIntrusiveState(common.Word(symbol), accept)
		return &fragment{
			Start:  start,
			Accept: accept,
		}
	}
}

func _eps(_ ...*fragment) *fragment {
	state := newIntrusiveState(common.Epsilon)
	return &fragment{
		Start:  state,
		Accept: state,
	}
}
