package main

import (
  "fmt"
  _"NFA/NFA"
  "NFA/Grammar"
)

func main() {
  rules := []string{
    "S -> A ",
    "A -> ε | a",
  }
  g := Grammar.MakeGrammar(rules)
  fmt.Println(g)
  nullSyms := g.NullableSymbols()
  fmt.Println(nullSyms)
}

