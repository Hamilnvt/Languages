package main

import (
  "fmt"
  _"NFA/NFA"
  "NFA/Grammar"
)

func main() {
  rules := []string{"S -> aS | a"}
  g := Grammar.MakeGrammar(rules)
  fmt.Println(g)
}

