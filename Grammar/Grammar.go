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
  //fmt.Println("Creating rule from:", rule_str)
  rule := Rule{} 
  parsed_rule := strings.Fields(rule_str)
  //fmt.Printf("%q\n", parsed_rule)
  if len(parsed_rule) < 3 {
    panic("Rule should have at least one prod (e.g. A -> a)")
  }
  rule.A = parsed_rule[0]
  if strings.Compare(parsed_rule[1], "->") != 0 {
    panic(fmt.Sprintf("ERROR: Malformed rule. Should be of the form: A -> a | ... | z"))
  }
  for i := 2; i < len(parsed_rule); i++ {
    prod := parsed_rule[i]
    if strings.Compare(prod, "|") != 0 {
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

//TODO se ad un certo punto è arrivato alla fine del file, ma non ha ancora parsato tutto: errore
func ParseGrammar(grammar_path string) Grammar {
  type Stage int
  const (
    s_BEGIN Stage = iota
    s_S
    s_T
    s_NT
    s_G
    s_END
  )
  stage := s_BEGIN
  fmt.Println("Begin stage:", stage)

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
        //grammar.R[rule.A] = union(grammar.R[rule.A], rule.prods)
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
  return grammar
}

func (G Grammar) addRule(rule Rule) {
  if G.R == nil {
    G.R = make(map[NonTerminal][]Prod)
  }
  if isStringIn(rule.A, G.NT) {
    missing := false
    var missing_symbol string
    for _, prod := range rule.prods {
      for _, symbol := range prod {
        if !(isStringIn(string(symbol), G.T) || isStringIn(string(symbol), G.NT)) {
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

  return
}

//TODO non è proprio ottimizzata, dovrebbe uscire quando len(nullSyms) == len(G.NT)
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
  if f == EPS {
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
    i := 0
    found := false
    for !found && i < len(first_X) {
      if string(first_X[i]) == EPS {
        found = true
        first_X[i] = first_X[len(first_X)-1]
        first_X = first_X[:len(first_X)-1]
      }
      i++
    }
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
