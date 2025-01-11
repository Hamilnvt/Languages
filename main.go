package main

import (
  "fmt"
  "Languages/NFA"
  "Languages/Grammar"
  "Languages/LexicalAnalyzer"
  "Languages/Parsing"
  "Languages/Utils"
)

//TODO
// - controllare se uno stato esiste già, in quel caso la goto non creerà lo stato, ma ritornerà lo stato già esistente così che la delta lo prenda
func main() {
  G := Grammar.ParseGrammar("./Grammar/Grammars/ProvaLR0.g")
  fmt.Println(G)

  initial_rule := Grammar.MakeRule("InitialSymbol_LR0 -> S")
  G.NT = append(G.NT, "InitialSymbol_LR0")
  G.AddRule(initial_rule)
  initial_item := NFA.ItemLR0{
    A: "ITLR0",
    Prod: Grammar.Production{G.S},
    Dot: 0,
  }

  initial_state := NFA.Closure(&G, NFA.CA_State{initial_item})
  states := make([]NFA.CA_State, 1)

  type DeltaPair struct {
    state *NFA.CA_State
    term string
  }

  delta := make(map[DeltaPair]*NFA.CA_State)
  states[0] = initial_state
  fmt.Printf("Initial state:\n%v\n", initial_state)
  
  queue := Utils.Queue[NFA.CA_State]{}
  queue.Enqueue(initial_state)

  for !queue.IsEmpty() {
    current_state, _ := queue.Dequeue()
    fmt.Println(current_state)

    ttg_map := make(map[string]bool)
    fmt.Println(ttg_map)
    for _, item := range current_state {
      if item.Dot < len(item.Prod) {
        ttg_map[item.Prod[item.Dot]] = true
      }
    }
    terms_to_go := make(Grammar.Production, len(ttg_map))
    i := 0
    for term := range ttg_map {
      terms_to_go[i] = term
      i++
    }
    fmt.Println("Terms to go", terms_to_go)

    for _, term := range terms_to_go {
      new_state := NFA.Goto(&G, current_state, term)
      fmt.Printf("New state\n%v\n", new_state)
      states = append(states, new_state)
      //TODO se è uguale non aggiungerlo, giustamente, dico io
      queue.Enqueue(new_state)
      delta[DeltaPair{state:&states[0], term:term}] = &states[len(states)-1]
    }
    fmt.Println("States after to go:")
    for i, state := range states {
      fmt.Printf("State %v:\n%v\n", i, state)
    }
  }
  //fmt.Println("Delta after to go:")
  //for key, value := range delta {
  //  fmt.Printf("Delta: {\nfrom\n%v\nto\n%v\nwith '%v'\n}\n", key.state, value, key.term)
  //}
  //CA := NFA.MakeCanonicAutomatonLR0(&G)
  //fmt.Println(CA)
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
