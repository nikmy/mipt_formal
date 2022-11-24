package parsers

import (
    "mipt_formal/gram/internal/cf"
    "mipt_formal/tools"
    "strings"
)

func proxyStart(g *cf.Grammar) {
    g.Rules = append(g.Rules, cf.Rule{
        Left:  byte(1),
        Right: string(cf.Start),
    })
    n := len(g.Rules) - 1
    g.Rules[0], g.Rules[n] = g.Rules[n], g.Rules[0]

    for i := range g.Rules {
        g.Rules[i].Right = strings.Replace(g.Rules[i].Right, "_", "", -1)
    }
}

func Earley(g *cf.Grammar, w string) bool {
    proxyStart(g)

    D := make([]tools.Set[Situation], len(w)+1)
    for i := range D {
        D[i] = tools.NewSet[Situation]()
    }
    D[0].Insert(InitSituation(g))

    newSituations := 1
    for newSituations > 0 {
        newSituations = Predict(D, 0, g) + Complete(D, 0)
    }

    for i := 1; i < len(w)+1; i++ {
        newSituations = Scan(D, i-1, w)
        for newSituations > 0 {
            newSituations = Predict(D, i, g)
            newSituations += Complete(D, i)
        }
    }

    return D[len(w)].Has(FinalSituation(g))
}

/*
	Scan
	  ( A -> u.xv , i ) in D[j]
	 ----------------------------- w[j] = x
	  ( A -> ux.v , i ) in D[j+1]
*/
func Scan(D []tools.Set[Situation], j int, w string) int {
    newSituations := 0
    for situation := range D[j] {
        if situation.Next() == w[j] {
            if D[j+1].Insert(situation.ReadNext()) {
                newSituations++
            }
        }
    }
    return newSituations
}

/*
	Predict
	  ( A -> u.Bv , i ) in D[j]
	 --------------------------- (B -> t) in rules
	  ( B -> .t , j ) in D[j]
*/
func Predict(D []tools.Set[Situation], j int, g *cf.Grammar) int {
    newSituations := 0
    for situation := range D[j] {
        for i, rule := range g.Rules {
            if rule.Left == situation.Next() {
                if D[j].Insert(ParseRule(&g.Rules[i], j)) {
                    newSituations++
                }
            }
        }
    }
    return newSituations
}

/*
	Complete
	  ( B -> t. , k ) in D[j]
	 --------------------------- ( A -> u.Bv , i ) in D[k]
	  ( A -> uB.v , i ) in D[j]
*/
func Complete(D []tools.Set[Situation], j int) int {
    newSituations := 0
    for curr := range D[j] {
        if !curr.Finished() {
            continue
        }

        k := curr.NextPos()
        for prev := range D[k] {
            if curr.Rule.Left == prev.Next() {
                cand := prev.ReadNext()
                if D[j].Insert(cand) {
                    newSituations++
                }
            }
        }
    }
    return newSituations
}

type Situation struct {
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

func ParseRule(rule *cf.Rule, position int) Situation {
    /*
    	Input:  (A -> u), i
    	Output: (A -> .u, i)
    */
    return Situation{
        Rule:    rule,
        RulePos: 0,
        WordPos: position,
    }
}

func InitSituation(g *cf.Grammar) Situation {
    /*
    	Returns (S' -> .S , 0)
    */
    return Situation{
        Rule:    &g.Rules[0],
        RulePos: 0,
        WordPos: 0,
    }
}

func FinalSituation(g *cf.Grammar) Situation {
    /*
    	Input:  w
    	Returns (S' -> S. , 0)
    */
    return Situation{
        Rule:    &g.Rules[0],
        RulePos: 1,
        WordPos: 0,
    }
}

func (s Situation) NextPos() int {
    /*
       (A -> u.v , i) => i
    */
    return s.WordPos
}

func (s Situation) Next() byte {
    /*
    	If s like (A -> u.xv, i) returns x
    	If s like (A -> u.  , i) returns nil
    */
    if s.Finished() {
        return byte(0)
    }
    return s.Rule.Right[s.RulePos]
}

func (s Situation) ReadNext() Situation {
    /*
    	(A -> u.xv, i) => (A -> ux.v, i)
    */
    ret := s
    if !s.Finished() {
        ret.RulePos++
    }
    return ret
}

func (s Situation) Finished() bool {
    /*
    	If s like (A -> u. , i) returns true
    	If s like (A -> u.w, i) returns false
    */
    return s.RulePos == len(s.Rule.Right)
}
