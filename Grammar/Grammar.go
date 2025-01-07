package Grammar

import (
  "fmt"
  "strings"
  _"unicode"
  "os"
  "bufio"
)

const (
  EPS = "ε"
)

type Terminal = string
type NonTerminal = string
type Prod = string

type Grammar struct {
  S string
  NT []NonTerminal
  T []Terminal
  R map[NonTerminal][]Prod
  FirstTable map[NonTerminal][]Terminal
  FollowTable map[NonTerminal][]Terminal
}

type Rule struct {
  A NonTerminal
  prods []Prod
}

func (rule Rule) String() (res string) {
  res += fmt.Sprintf("Rule: %v -> ", rule.A)
  for i, prod := range rule.prods {
    if i != len(rule.prods)-1 {
      res += fmt.Sprintf("%v | ", prod)
    } else {
      res += fmt.Sprintf("%v", prod)
    }
  }
  return
}

func makeRule(rule_str string) Rule {
  fmt.Println("Creating rule from:", rule_str)
  rule := Rule{} 
  parsed_rule := strings.Fields(rule_str)
  //fmt.Printf("%q\n", parsed_rule)
  if len(parsed_rule) < 3 {
    panic("Rule should have at least one prod (e.g. A -> a)")
  }
  rule.A = parsed_rule[0]
  if parsed_rule[1] != "->" {
    panic(fmt.Sprintf("ERROR: Malformed rule. Should be of the form: A -> a | ... | z"))
  }
  for i := 2; i < len(parsed_rule); i++ {
    prod := parsed_rule[i]
    if prod != "|" {
      //TODO contiene, non per forza all'inizio
      if len(prod) > 1 && prod[:2] == "\\|" {
        if len(prod) > 2 {
          prod = "|"+prod[2:]
        } else {
          prod = "|"
        }
      } else if prod == "\\eps" {
        prod = EPS
      }
      //fmt.Printf("%v -> %v\n", rule.A, parsed_rule[i])
      rule.prods = append(rule.prods, prod)
    }
  }
  return rule
}

func isStringIn(s string, strs []string) bool {
  found := false
  i := 0
  for !found && i < len(strs) {
    //fmt.Printf("%v == %v?\n", s, strs[i])
    if strings.Compare(s, strs[i]) == 0 {
      found = true
    }
    i++
  }
  return found
}

//TODO non inserisce correttamente i simboli. Ancora? Da testare
func MakeGrammar(rules []string, initialSymbol string, nonterminals []NonTerminal, terminals []Terminal) Grammar {
  parsed_rules := make([]Rule, 0)
  for _, rule := range rules {
    parsed_rules = append(parsed_rules, makeRule(rule))
  }
  R := make(map[NonTerminal][]Prod)
  G := Grammar{
    S: parsed_rules[0].A,
    NT: nonterminals,
    T: terminals,
    R: make(map[NonTerminal][]Prod),
  }
  for _, rule := range parsed_rules {
    fmt.Println(rule)
    G.addRule(rule)
  }
  fmt.Println("R:", R)
  return G
}

func ParseGrammar(grammar_path string) Grammar {
  if grammar_path[len(grammar_path)-2:] != ".g" {
    panic("File extension should be .g")
  }

  file, err := os.Open(grammar_path)
  if err != nil {
    panic(err)
  }

  fmt.Println("Scanning file:")
  clean_file := make([]string, 0)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    //fmt.Println(line)

    if len(line) > 0 && (line[0] != '#'){
      clean_file = append(clean_file, line)
    }
  }

  if err := scanner.Err(); err != nil {
    panic(err)
  }
  file.Close()
  fmt.Println("File scanning ended without errors.\n")
  for _, line := range clean_file {
    fmt.Println(line)
  }
  fmt.Println()

  fmt.Println("Parsing Definitions:")
	//symbolTable := make(map[string][]string)

  i := 0
  line := clean_file[i]
  if line != "DEFINE:" {
    if line != "GRAMMAR:" {
      panic("Invalid Definitions declaration, it should be of the form:\nDEFINE:\nDEF1 [def1]\n(You can omit it)")
    } else {
      fmt.Println("No Definitions declared")
    }
  } else {
    fmt.Println("TODO: DEFINITIONS")
  }

  grammar := Grammar{
    R: make(map[NonTerminal][]Prod),
    FirstTable: make(map[NonTerminal][]Terminal),
    FollowTable: make(map[NonTerminal][]Terminal),
  }

  fmt.Println("Parsing Grammar:")
  line = clean_file[i]
  if line != "GRAMMAR:" {
    panic("Invalid Grammar declaration, it should be of the form:\nGRAMMAR:\nA -> a")
  } else {
    i++
  }

  for j := i; j < len(clean_file); j++ {
    line := strings.Fields(clean_file[j])
    if len(line) <3 {
      panic("Invalid Rule declaration, it should have at least one right production:\nA -> a")
    }
    if line[1] != "->" {
      panic("Invalid Rule declaration, it should be of the form:\nA -> b_0 | ... | b_k")
    }
    nonTerminal := line[0]
    if j == 1 {
      grammar.S = nonTerminal
    }
    grammar.NT = union(grammar.NT, []string{nonTerminal})
    //fmt.Println(line)
  }

  for j := i; j < len(clean_file); j++ {
    line := clean_file[j]
    rule := makeRule(line)
    grammar.R[rule.A] = union(grammar.R[rule.A], rule.prods)
    //fmt.Println(grammar.R[rule.A])
  }

  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      for _, symbol := range prod {
        if !grammar.IsNonTerminal(string(symbol)) && string(symbol) != EPS {
          grammar.T = union(grammar.T, []string{string(symbol)})
        }
      }
    }
  }

  for _, nt := range grammar.NT {
    grammar.FirstTable[nt] = grammar.First(nt)
  }
  for _, nt := range grammar.NT {
    grammar.FollowTable[nt] = grammar.Follow(nt)
  }

  return grammar
}

//TODO se ad un certo punto è arrivato alla fine del file, ma non ha ancora parsato tutto: errore
func ParseGrammar_deprecated(grammar_path string) Grammar {
  if grammar_path[len(grammar_path)-2:] != ".g" {
    panic("File extension should be .g")
  }

  file, err := os.Open(grammar_path)
  if err != nil {
    panic(err)
  }

  fmt.Println("Scanning file:\n")
  clean_file := make([]string, 0)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())
    fmt.Println(line)

    if len(line) > 0 && (line[0] != '#'){
      clean_file = append(clean_file, line)
    }
  }

  if err := scanner.Err(); err != nil {
    panic(err)
  }
  file.Close()
  fmt.Println("\nScanning file ended without errors:\n")

  grammar := Grammar{
    R: make(map[NonTerminal][]Prod),
    FirstTable: make(map[NonTerminal][]Terminal),
    FollowTable: make(map[NonTerminal][]Terminal),
  }

  for _, line := range clean_file {
    fmt.Println(line)
  }
  fmt.Println()
  fmt.Println("Parsing Initial Symbol:")
  line := clean_file[0]
  if len(line) > 2 && line[:2] == "S:" {
    S := strings.Fields(strings.TrimSpace(line[2:]))
    if len(S) == 1 {
      grammar.S = S[0]
      fmt.Println("S:", S[0])
    } else {
      panic("Initial symbol must be unique")
    }
  } else {
    panic("Invalid Initial Symbol declaration, it should be of the form:\nS: <initial_symbol>")
  }

  fmt.Println("Parsing Terminals:")
  line = clean_file[1]
  i := 2
  if line == "T: {" {
    line = clean_file[i]
    for done := false; i < len(clean_file) && !done; line = clean_file[i] {
      if line != "}" {
        terminals := strings.Fields(line)
        //fmt.Println("Terminals:", terminals)
        for j, t := range terminals {
          if t == "\\}" {
            terminals[j] = "}"
          }
        }
        grammar.T = union(grammar.T, terminals)
      } else {
        done = true
      }
      i++
    }
  } else {
    panic("Invalid Terminals declaration, it should be of the form:\nT: { <a_0> ... <a_n> }\n(new lines are allowed in the body)")
  }
  fmt.Println("T:", grammar.T)

  fmt.Println("Parsing Nonterminals:")
  line = clean_file[i]
  if line == "NT: {" {
    i++
    line = clean_file[i]
    for done := false; i < len(clean_file) && !done; line = clean_file[i] {
      if line != "}" {
        nonterminals := strings.Fields(line)
        //fmt.Println("Nonterminals:", nonterminals)
        for j, t := range nonterminals {
          if t == "\\}" {
            nonterminals[j] = "}"
          }
        }
        grammar.NT = union(grammar.NT, nonterminals)
      } else {
        done = true
      }
      i++
    }
  } else {
    panic("Invalid Nonterminals declaration, it should be of the form:\nNT: { <A_0> ... <A_n> }\n(new lines are allowed in the body)")
  }
  if inter := intersection(grammar.T, grammar.NT); len(inter) > 0 {
    panic(fmt.Sprintf("Symbols cannot be both Terminals and Nonterminals: %v", inter))
  }
  fmt.Println("NT:", grammar.NT)

  fmt.Println("Parsing Rules:")
  line = clean_file[i]
  if line == "R: {" {
    i++
    line = clean_file[i]
    for done := false; i < len(clean_file) && !done; line = clean_file[i] {
      if line != "}" {
        rule := makeRule(line)
        //fmt.Println(rule)
        grammar.addRule(rule)
        i++
      } else {
        done = true
      }
    }
  } else {
    panic("Invalid Rules declaration, it should be of the form:\nR: { <A_i> -> a_0 | ... | a_k }")
  }

  usedTerminals := make(map[string]bool)
  for _, t := range grammar.T {
    usedTerminals[t] = false
  }
  //fmt.Println("UsedTerminals initialization:", usedTerminals)
  for _, nt := range grammar.NT {
    if len(grammar.R[nt]) == 0 {
      panic(fmt.Sprintf("There isn't a rule associated to the NonTerminal %v", nt))
    }
    for _, prod := range grammar.R[nt] {
      for _, t := range grammar.T {
        if used := usedTerminals[t]; !used && strings.Contains(prod, t) {
          //fmt.Println(t, "is used")
          usedTerminals[t] = true
        }
      }
    }
  }
  //fmt.Println("UsedTerminals:", usedTerminals)
  for t, used := range usedTerminals {
    if !used {
      panic(fmt.Sprintf("Terminal %v is not used in any rule", t))
    }
  }
  //TODO non li stampa in ordine (perché è una mappa)
  fmt.Println("R:")
  for _, rule := range grammar.R {
    fmt.Println(rule)
  }

  if i < len(clean_file)-1 {
    panic("Invalid end of file")
  }

  fmt.Println("Grammar parsed successfully")

  for _, nt := range grammar.NT {
    grammar.FirstTable[nt] = grammar.First(nt)
  }
  for _, nt := range grammar.NT {
    grammar.FollowTable[nt] = grammar.Follow(nt)
  }

  return grammar
}

func (G Grammar) addRule(rule Rule) {
  fmt.Println(rule)
  if G.R == nil {
    G.R = make(map[NonTerminal][]Prod)
  }
  if isStringIn(rule.A, G.NT) {
    missing := false
    var missing_symbol string
    for _, prod := range rule.prods {
      for _, symbol := range prod {
        if string(symbol) != EPS && !(isStringIn(string(symbol), G.T) || isStringIn(string(symbol), G.NT)) {
          missing = true
          missing_symbol = string(symbol)
          break
        }
      }
      if missing {
        break
      }
    }
    if !missing {
      G.R[rule.A] = union(G.R[rule.A], rule.prods)
    } else {
      panic(fmt.Sprintf("%v is not a symbol of the grammar (Terminals %v, NonTerminals %v)\n", missing_symbol, G.T, G.NT))
    }
  } else {
    panic(fmt.Sprintf("%v isn't a NonTerminal in %v\n", rule.A, G.NT))
  }
}

func (G Grammar) String() (res string) {
  res += fmt.Sprintf("\nPrinting Grammar:\n")
  res += fmt.Sprintf("Initial Symbol: %v\n", G.S)
  res += fmt.Sprintf("Non-Terminals: %v\n", G.NT)
  res += fmt.Sprintf("Terminals: %v\n", G.T)
  res += fmt.Sprintf("Rules:\n")
  for _, nt := range G.NT {
    res += fmt.Sprintf("%v -> ", nt)
    for i, prod := range G.R[nt] {
      if i != len(G.R[nt])-1 {
        res += fmt.Sprintf("%v | ", prod)
      } else {
        res += fmt.Sprintf("%v\n", prod)
      }
    } 
  }
  res += fmt.Sprintf("First:\n")
  for _, nt := range G.NT {
    res += fmt.Sprintf("%v: %v\n", nt, G.FirstTable[nt])
  }

  res += fmt.Sprintf("Follow:\n")
  for _, nt := range G.NT {
    res += fmt.Sprintf("%v: %v\n", nt, G.FollowTable[nt])
  }

  return
}

//TODO non è proprio ottimizzata, dovrebbe uscire quando len(nullSyms) == len(G.NT) (fa più iterazioni del necessario)
func (G Grammar) NullableSymbols() []NonTerminal {
  nullSyms := make([]NonTerminal, 0)
  for _, nt := range G.NT {
    for _, prod := range G.R[nt] {
      if strings.Compare(EPS, prod) == 0 {
        nullSyms = append(nullSyms, nt)
      }
    } 
  }
  done := false
  for !done && len(nullSyms) < len(G.NT) {
    done = true
    for _, nt := range G.NT {
      //fmt.Println(i, "NonTerminal", nt)
      for _, prod := range G.R[nt] {
        //fmt.Println(j, "Prod", prod)
        found := false
        for _, c := range prod {
          //fmt.Println(k, "symbol", string(c))
          if !isStringIn(string(c), nullSyms) {
            found = true 
            break
          }
        }
        if !found && !isStringIn(nt, nullSyms) {
            nullSyms = append(nullSyms, nt)
            done = false
            break
        }
      } 
      //fmt.Println("Nullable symbols:", nullSyms)
    }
  }
  
  return nullSyms
}

func union(slice1, slice2 []string) []string {
  // Create a map to store the elements of the union
  values := make(map[string]bool)
  for _, key := range slice1 { // for loop used in slice1 to remove duplicates from the values
    values[key] = true
  }
  for _, key := range slice2 { // for loop used in slice2 to remove duplicates from the values
    values[key] = true
  }
  // Convert the map keys to a slice
  output := make([]string, 0, len(values)) //create slice output
  for val := range values {
    output = append(output, val) //append values in slice output
  }
  return output
}
func intersection(slice1, slice2 []string) []string {
  values := make(map[string]bool)
  for _, key := range slice1 {
    values[key] = true
  }
  output := make([]string, 0, len(values))
  for _, key := range slice2 {
    if values[key] {
      output = append(output, key)
    }
  }
  return output
}

//TODO 
// - si potrebbe fare che se calcola il first di un nonterminale lo inserisce nella mappa e se lo deve ricacolare, prima di farlo controlla la tabella
func (G Grammar) first_interface(f string, first []string, nullSyms []string) []string {
  //fmt.Println("Calculating first of", f)
  if f == EPS || f == "" {
    return append(first, EPS)
  }
  first_symbol := string(f[0])
  if isStringIn(first_symbol, G.T) {
    return append(first, first_symbol)
  }
  if !isStringIn(first_symbol, nullSyms) {
    //fmt.Println(first_symbol, "is not nullable")
    for _, prod := range G.R[first_symbol] {
      if string(prod[0]) != first_symbol {
        first = G.first_interface(prod, first, nullSyms)
      }
    }
  } else {
    //fmt.Println(first_symbol, "is nullable")
    first_X := make([]string, 0)
    for _, prod := range G.R[first_symbol] {
      if string(prod[0]) != first_symbol {
        first_X = G.first_interface(prod, first_X, nullSyms)
      } else if len(prod) > 1 {
        first_X = G.first_interface(prod[1:], first_X, nullSyms)
      }
    }
    //fmt.Println(first_X)

    found := false
    first_X_wo_eps := make([]string, 0)
    for _, f := range first_X {
      if f != EPS {
        first_X_wo_eps = append(first_X_wo_eps, f)
      } else {
        found = true
      }
    }

    first_X = first_X_wo_eps
    if found {
      //fmt.Println("without EPS:", first_X)
    }
    if len(first_X) > 0 {
      first = union(first, first_X)
    }
    if len(f) > 1 {
      if string(f[1]) != first_symbol {
        first = union(first, G.First(f[1:]))
      }
    } else {
      first = append(first, EPS)
    }
  }
  return first
}

func (G Grammar) First(f string) []string {
  //TODO sfruttando le mappe (come in union) è molto più facile
  first := G.first_interface(f, make([]string, 0), G.NullableSymbols())
  i := 0
  for i < len(first)-1 {
    j := i+1
    removed := false
    for j < len(first) {
      if first[i] == first[j] {
        removed = true
        first[j] = first[len(first)-1]
        first = first[:len(first)-1]
      } else {
        j++
      }
    }
    if !removed {
      i++
    }
  }
  return first
}

func (G Grammar) follow_interface(Y NonTerminal, from []string) []string {
  if isStringIn(Y, from) {
    return []string{}
  }
  follow := make([]string, 0)
  if Y == G.S {
    follow = append(follow , "$")
  }
  //fmt.Printf("Calculating Follow(%v)\n", Y)
  for _, X := range G.NT {
    for _, prod := range G.R[X] {
      for i := range prod {
        symbol := string(prod[i])
        if symbol == Y {
          //fmt.Printf("%v in pos %v\n", prod, i)
          beta := ""
          if i != len(prod)-1 {
            beta = prod[i+1:]
          } 
          first := G.First(beta)
          //fmt.Printf("First(%v) = %v\n", beta, first)

          if isStringIn(EPS, first) {
            follow = union(follow, G.follow_interface(X, append(from, Y)))
          }
          first_wo_eps := make([]string, 0)
          for _, f := range first {
            if f != EPS {
              first_wo_eps = append(first_wo_eps, f)
            }
          }
          //fmt.Printf("First(%v) \\ {%v} = %v\n", beta, EPS, first_wo_eps)
          follow = union(follow, first_wo_eps)
        }
      }
    } 
  }
  return follow
}

func (G Grammar) Follow(A NonTerminal) []string {
  return G.follow_interface(A, make([]string, 0))
}

func (G Grammar) IsTerminal(X string) bool {
  return isStringIn(X, G.T)
}
func (G Grammar) IsNonTerminal(X string) bool {
  return isStringIn(X, G.NT)
}
