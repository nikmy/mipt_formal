# Context-Free grammar parsers library

## Formats

### Grammar

**Only CF-grammars are supported**

```
  |   alternation
  .   concatenation (omit spaces)
"..." raw string
 { }  repeating

non-terminal ::= A | B | ... | Z
terminal     ::= a | b | ... | z
rule         ::= non-terminal . " ::= " . { non-terminal | terminal }
```

## Parsers
