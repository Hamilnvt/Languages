Work in progress

# What is this?
It's a personal project about my university studies on programming languages. 
It has the ambition to recreate some sort of compiler compiler like yacc-lex (also, I'm learning Go, I don't know why, should've gone with Rust? Let me know). 
Someday the set of tools will comprehend:
- Lexical Analysis:
  -  [x] Creation of NFA from regular expressions
  -  [x] Creation of DFA from NFA (and its minimization)
  -  [ ] Lexical Analyzer
- Grammars:
  - [x] Grammar parsing from file
    - [ ] command-line
  - [x] Nullable symbols
  - [x] First
  - [x] Follow
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
  - [ ] Parse Tree
    - [x] Concrete
    - [ ] Abstract
  - [x] Top-Down parser LL(1)
  - [ ] Bottom-Up parsers:
    - [x] LR(0)
    - [x] SLR(1)
    - [ ] LR(1)
    - [ ] LALR(1)

# Parsing Grammar from file

You can parse a grammar by invoking the function Grammar.ParseGrammar("/path/to/grammar.g").  
This function returns a Grammar (see below).

TODO: someday there will be a command-line for this

Any file .g should follow this syntax:  
```
# This is a comment and can only be placed at the beginning of a line

# Blank lines and comments will be ignored.
```

## Definitions of Terminals definitions which you can use later in the Grammar declaration.

_Syntax for this is work in progress..._ 
```
DEFINE:

name1 "def1"
name2 ("def2" name1)
...

# name2 will be "def1def2"
```

## List of Rules of the Grammar.

White spaces separate the terms, while '|' separate the productions, so you can use every character you want except for ' ' and '|', which you'll need to escape.

> Please, I beg you, don't use '#' as a NonTerminal, it'll be recognized as a comment, I didn't bother myself implementing the escaping

```
GRAMMAR:

  A_0 -> b_00 | ... | b_0k
  A_1 -> ...
  .
  .
  .
  A_i -> b_i0 | ... | b_ik

```

Anything other than that will throw a panic error at your face.

--------------

A Grammar is a type:
```go
type Grammar struct {
  S  string              // initial symbol
  NT []string            // non terminals
  T  []string            // terminals
  R  map[string][]string // rules (A -> a_0 | ... | a_k is stored as Grammar.R[A] = [a_0 ... a_k])
}
```
