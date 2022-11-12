# Context-Free grammar parsers library

## Formats

### Grammar

- **Only CF-grammars are supported**
- **S is start symbol**
- **_ is empty word**

```
  |   alternation
  .   concatenation (omit spaces)
"..." raw string
 { }  repeating

non-terminal ::= A | B | ... | Z
terminal     ::= a | b | ... | z | _
rule         ::= non-terminal . " -> " . { non-terminal | terminal } . "\n"
grammar      ::= "S" . " ::= " . { non-terminal | terminal } . "\n" . { rule }
```

## Parsers
