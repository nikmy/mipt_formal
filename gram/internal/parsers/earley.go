package parsers

import (
    "strings"

    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
)

const (
    endOfString = byte(0)
    fakeStart   = byte(1)
)

func addProxyStart(grammar *cf.Grammar) {
    if grammar.Rules[0].Left == fakeStart {
        return
    }

    grammar.Rules = append(grammar.Rules, cf.Rule{
        Left:  fakeStart,
        Right: string(cf.Start),
    })
    n := len(grammar.Rules) - 1
    grammar.Rules[0], grammar.Rules[n] = grammar.Rules[n], grammar.Rules[0]

    for i := range grammar.Rules {
        grammar.Rules[i].Right = strings.Replace(grammar.Rules[i].Right, "_", "", -1)
    }
}

func Earley(g *cf.Grammar, word string) bool {
    addProxyStart(g)

    history := newSituationSet(len(word) + 1)
    delta := newSituationSet(len(word) + 1)

    history.Add(0, initSituation(g))
    delta.Add(0, initSituation(g))

    for delta.Size() > 0 {
        nextDelta := newSituationSet(len(word) + 1)
        predict(delta, history, 0, g, nextDelta)
        complete(delta, history, 0, nextDelta)
        delta = nextDelta
    }

    for i := 1; i < len(word)+1; i++ {
        delta = scan(history, i-1, word)
        for delta.Size() > 0 {
            nextDelta := newSituationSet(len(word) + 1)
            predict(delta, history, i, g, nextDelta)
            complete(delta, history, i, nextDelta)
            delta = nextDelta
        }
    }

    return !history.Add(len(word), finalSituation(g))
}

/*
	scan
	  ( A -> u.xv , i ) in D[j]
	 ----------------------------- w[j] = x
	  ( A -> ux.v , i ) in D[j+1]
*/
func scan(history *situationSet, j int, w string) *situationSet {
    delta := newSituationSet(len(w) + 1)
    situations := history.Get(j, w[j])
    for situation := range situations {
        next := situation.ReadNext()
        if history.Add(j+1, next) {
            delta.Add(j+1, next)
        }
    }
    return delta
}

/*
	predict
	  ( A -> u.Bv , i ) in D[j]
	 --------------------------- (B -> t) in rules
	  ( B -> .t , j ) in D[j]
*/
func predict(delta *situationSet, history *situationSet, j int, g *cf.Grammar, newDelta *situationSet) {
    for _, situations := range delta.GetMappedSet(j) {
        for situation := range situations {
            for i, rule := range g.Rules {
                if rule.Left == situation.Next() {
                    newSituation := parseRule(&g.Rules[i], j)
                    if history.Add(j, newSituation) {
                        newDelta.Add(j, newSituation)
                    }
                }
            }
        }
    }
}

/*
	complete
	  ( B -> t. , k ) in D[j]
	 --------------------------- ( A -> u.Bv , i ) in D[k]
	  ( A -> uB.v , i ) in D[j]
*/
func complete(delta *situationSet, history *situationSet, j int, newDelta *situationSet) {
    for _, situations := range delta.GetMappedSet(j) {
        for situation := range situations {
            if !situation.Finished() {
                continue
            }

            k := situation.NextPos()
            parents, exists := history.GetMappedSet(k)[situation.Rule.Left]
            if !exists {
                continue
            }
            for parent := range parents {
                if situation.Rule.Left == parent.Next() {
                    cand := parent.ReadNext()
                    if history.Add(j, cand) {
                        newDelta.Add(j, cand)
                    }
                }
            }
        }
    }
}

func newSituationSet(nSets int) *situationSet {
    data := make([]map[byte]tools.Set[earleySituation], nSets)
    for i := range data {
        data[i] = make(map[byte]tools.Set[earleySituation])
    }
    return &situationSet{
        data: data,
        size: 0,
    }
}

type situationSet struct {
    data []map[byte]tools.Set[earleySituation]
    size int
}

func (s *situationSet) Size() int {
    return s.size
}

func (s *situationSet) Add(setID int, situation earleySituation) bool {
    nextSymbol := situation.Next()
    if _, has := s.data[setID][nextSymbol]; !has {
        s.data[setID][nextSymbol] = tools.NewSet[earleySituation]()
    }
    if s.data[setID][nextSymbol].Insert(situation) {
        s.size++
        return true
    }
    return false
}

func (s *situationSet) Get(setID int, nextSymbol byte) tools.Set[earleySituation] {
    return s.GetMappedSet(setID)[nextSymbol]
}

func (s *situationSet) GetMappedSet(setID int) map[byte]tools.Set[earleySituation] {
    return s.data[setID]
}

type earleySituation struct {
    Rule    *cf.Rule
    RulePos int
    WordPos int
}

/*
	Notation
		x    single terminal or non-terminal
		.    current position in the rule
		u, v words consists of terminals and non-terminals
*/

func parseRule(rule *cf.Rule, position int) earleySituation {
    /*
    	Input:  (A -> u), i
    	Output: (A -> .u, i)
    */
    return earleySituation{
        Rule:    rule,
        RulePos: 0,
        WordPos: position,
    }
}

func initSituation(g *cf.Grammar) earleySituation {
    /*
    	Returns (S' -> .S , 0)
    */
    return earleySituation{
        Rule:    &g.Rules[0],
        RulePos: 0,
        WordPos: 0,
    }
}

func finalSituation(g *cf.Grammar) earleySituation {
    /*
    	Input:  w
    	Returns (S' -> S. , 0)
    */
    return earleySituation{
        Rule:    &g.Rules[0],
        RulePos: 1,
        WordPos: 0,
    }
}

func (s earleySituation) NextPos() int {
    /*
       (A -> u.v , i) => i
    */
    return s.WordPos
}

func (s earleySituation) Next() byte {
    /*
    	If s like (A -> u.xv, i) returns x
    	If s like (A -> u.  , i) returns nil
    */
    if s.Finished() {
        return endOfString
    }
    return s.Rule.Right[s.RulePos]
}

func (s earleySituation) ReadNext() earleySituation {
    /*
    	(A -> u.xv, i) => (A -> ux.v, i)
    */
    ret := s
    if !s.Finished() {
        ret.RulePos++
    }
    return ret
}

func (s earleySituation) Finished() bool {
    /*
    	If s like (A -> u. , i) returns true
    	If s like (A -> u.w, i) returns false
    */
    return s.RulePos == len(s.Rule.Right)
}
