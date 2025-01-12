package main

import (
  "fmt"
  "Languages/NFA"
  "Languages/Grammar"
  "Languages/LexicalAnalyzer"
  "Languages/Parsing"
  _"Languages/Utils"
)

//TODO
// - controllare se uno stato esiste già, in quel caso la goto non creerà lo stato, ma ritornerà lo stato già esistente così che la delta lo prenda
func main() {
  G := Grammar.ParseGrammar("./Grammar/Grammars/ProvaLR0.g")
  fmt.Println(G)
  parser := Parsing.MakeParserBottomUpLR0(G)
  fmt.Println("parser:", parser)
}

func canonic_automaton_LR0() {
  G := Grammar.ParseGrammar("./Grammar/Grammars/ProvaLR0.g")
  fmt.Println(G)
  CA := NFA.MakeCanonicAutomatonLR0(&G)
  fmt.Println(CA)
}

func ParsingAndLL1() {
  //G := Grammar.ParseGrammar("./Grammar/Grammars/LL1_prova.g")
  //G := Grammar.ParseGrammar("./Grammar/Grammars/GrammarIfStatement.g")
  //G := Grammar.ParseGrammar("./Grammar/Grammars/FirsFollowG2.g")
  G := Grammar.ParseGrammar("./Grammar/Grammars/RegExpGrammar.g")
  fmt.Println(G)
  parser, err := Parsing.MakeParserTopDownLL1(G)
  if err != nil {
    panic(err)
  }
  //input := "(a+(a*b)/(b-a))"
  //input := "(ab)?a"
  //input := "if true then a else b"
  input := "a*"
  tree, err := parser.Parse(input)
  if err != nil {
    panic(err)
  }
  fmt.Println("Parsed String:", tree.GetParsedString())
  fmt.Println(tree)
}

func makeLA() {
  la := LexicalAnalyzer.MakeLexicalAnalyzer("./LexicalAnalyzer/LAs/prova.la") 
  fmt.Println(la)
}

func grammar_first_and_follow() {
  G := Grammar.ParseGrammar("./Grammar/Grammars/RegExpGrammar.g")
  fmt.Println(G)

  fmt.Println("First:")
  for _, nt := range G.NT {
    fmt.Println(nt,G.First([]string{nt}))
  }

  fmt.Println("\nFollow:")
  for _, nt := range G.NT {
    fmt.Println(nt, G.Follow(nt))
  }
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
    fmt.Println(G.NT[i], ":", G.First([]string{G.NT[i]}))
  }
}
