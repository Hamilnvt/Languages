package Parsing

import (
  _"fmt"
  "Languages/NFA"
  "Languages/Utils"
)

type CanonicAutomatonLR0 struct {
  dfa NFA.DFA
  states map[int]*NFA.State
}

type BUTablePair struct {
  state int
  term string
}

type BUParsingTable map[BUTablePair][]string  

type Parser_LR0 struct {
  table BUParsingTable 
  stack Utils.Stack[string]
}

