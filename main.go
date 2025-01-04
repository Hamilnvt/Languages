package main

import (
  "fmt"
  _"NFA/NFA"
  "NFA/Grammar"
)

func main() {
  G := Grammar.ParseGrammar("./Grammar/Grammars/SimpleGrammar.g")
  fmt.Println(G)
}

func grammar_example() {
  rules := []string{
    "S -> BCA | ABE ",
    "A -> a | aDb | bSc ",
    "B -> C | Bb ",
    "C -> ε | dC ",
    "D -> dD ",
    "E -> D | dE ",
  }
  initialSymbol := "S"
  nonterminals := []Grammar.NonTerminal{initialSymbol, "A", "B", "C", "D", "E"}
  terminals := []Grammar.Terminal{"a", "b", "c", "d"}

  //rules := []string{
  //  "S -> S",
  //}
  //initialSymbol := "S"
  //nonterminals := []Grammar.NonTerminal{initialSymbol}
  //terminals := []Grammar.Terminal{}

  G := Grammar.MakeGrammar(rules, initialSymbol, nonterminals, terminals)
  fmt.Println(G)

  nullSyms := G.NullableSymbols()
  fmt.Println("N(G) =", nullSyms)

  fmt.Println("Calculanting the first of all nonterminals")
  for i := len(G.NT)-1; i >= 0; i-- {
    fmt.Println(G.NT[i], ":", G.First(G.NT[i]))
  }
}
