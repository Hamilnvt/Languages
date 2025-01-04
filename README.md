Work in progress

# What is this?
It's a personal project about my university studies on programming languages.\
It has the ambition to recreate some sort of compiler compiler like yacc-lex (also, I'm learning Go, I don't know why, should've gone with Rust? Let me know).\
Someday the set of tools will comprehend:
- Lexical Analysis:
  -  [x] Creation of NFA from regular expressions
  -  [x] Creation of DFA from NFA (and its minimization)
  -  [ ] Lexical Analyzer
- Grammars:
  - [x] Grammar parsing from file
  - [x] Nullable symbols
  - [x] First
  - [ ] Follow
  - [ ] Generators
  - [ ] Reachable symbols
  - [ ] Singular symbols
  - [ ] Grammars simplification:
    - [ ] ε-productions
    - [ ] singular productions
    - [ ] useless symbols
    - [ ] left recursion (direct and indirect)
    
- Syntax Analysis:
  - [ ] Creation of PDA from Grammar
  - [ ] Creation of DPDA from Grammar
  - [ ] Top-Down parser LL(1)
  - [ ] Bottom-Up parsers:
    - [ ] LR(0)
    - [ ] SLR(1)
    - [ ] LR(1)
    - [ ] LALR(1)

# Parsing Grammar from file

You can parse a grammar by invoking the function Grammar.ParseGrammar("/path/to/grammar.g").
This function returns a Grammar (see below, or above?)

Any file .g should follow this syntax:

```
# This is a comment and can be placed only at the beginning of a line (for the moment)
# Blank lines will be ignored

S: initial_symbol

# List of Terminals
T: {
  a_00 ... a_0k
  a_10 ...
  .
  .
  .
  a_i0 ... a_ik
}

# white spaces are the separators, so you can use every character you want
# except for '}' and '#' which you'll need to escape like this \} and \#

# List of NonTerminals
NT: {
  A_00 ... A_0k
  A_10 ...
  .
  .
  .
  A_i0 ... A_ik
}

# List of Rules (Spaces between right productions are important!)
R: {
  A_0 -> b_00 | ... | b_0k
  A_1 -> ...
  .
  .
  .
  A_i -> b_i0 | ... | b_ik
}
```
Anything other than that will throw a panic error at your face.\

A Grammar is a type:\
```go
type Grammar struct {
  S  string              // initial symbol
  NT []string            // non terminals
  T  []string            // terminals
  R  map[string][]string // rules (A -> a0 | ... | ak is stored as Grammar.R[A] = [a0 ... ak])
}
```
