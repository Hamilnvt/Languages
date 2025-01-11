package Grammar

import (
  "fmt"
  "strings"
  _"unicode"
  "os"
  "bufio"
  "Languages/Utils"
)

const (
  EPS = "ε"
)

type Terminal = string
type NonTerminal = string
type Production = []string

type Grammar struct {
  S string
  NT []NonTerminal
  T []Terminal
  R map[NonTerminal][]Production
  FirstTable map[NonTerminal][]Terminal
  FollowTable map[NonTerminal][]Terminal
}

type Rule struct {
  A NonTerminal
  prods []Production
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

//TODO non fa l'unione, ma solo l'append (se ci sono due regole per lo stesso nonterminale fa casino)
func MakeRule(rule_str string) Rule {
  fmt.Println("Creating rule from:", rule_str)
  rule := Rule{
    prods: make([]Production, 0),
  } 
  parsed_rule := strings.Fields(rule_str)
  fmt.Printf("Parsed rule: %q\n", parsed_rule)
  if len(parsed_rule) < 3 {
    panic("Rule should have at least one prod (e.g. A -> a)")
  }
  rule.A = parsed_rule[0]
  if parsed_rule[1] != "->" {
    panic(fmt.Sprintf("ERROR: Invalid rule. Should be of the form: A -> a | ... | z"))
  }
  //TODO ' ', '\n', '\t' (controllare anche quando vengono stampati)
  parsed_rule = strings.Split(strings.Join(parsed_rule[2:], " "), " | ")
  fmt.Printf("Ready to check rule: %q\n", parsed_rule)
  for i := 0; i < len(parsed_rule); i++ {
    prod := strings.Fields(parsed_rule[i])
    for i, term := range prod {
      switch term {
      case "\\eps":
        prod[i] = EPS
      case "\\|":
        prod[i] = "|"
      case "\\ ":
        prod[i] = " "
      case "\\n":
        prod[i] = "\n"
      case "\\t":
        prod[i] = "\t"
      }
    }
    fmt.Printf("%v -> %v\n", rule.A, prod)
    rule.prods = append(rule.prods, prod)
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
    parsed_rules = append(parsed_rules, MakeRule(rule))
  }
  R := make(map[NonTerminal][]Production)
  G := Grammar{
    S: parsed_rules[0].A,
    NT: nonterminals,
    T: terminals,
    R: make(map[NonTerminal][]Production),
  }
  for _, rule := range parsed_rules {
    fmt.Println(rule)
    G.AddRule(rule)
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
	symbolTable := make(map[string]string)

  i := 0
  line := clean_file[i]
  if line != "DEFINE:" {
    if line != "GRAMMAR:" {
      panic("Invalid Definitions declaration, it should be of the form:\nDEFINE:\nDEF1 [def1]\n(You can omit it)")
    } else {
      fmt.Println("No Definitions declared")
    }
  } else {
    i++
    line = clean_file[i]
    for i < len(clean_file) && line != "GRAMMAR:" {
      definition := strings.SplitN(line, " ", 2)
      fmt.Println(definition)
      if len(definition) != 2 {
        panic("Invalid definition, it should be of the form:\nname \"def\"\nor:\nname: (\"def\" other_name ...)")
      }
      name := definition[0]
      if _, ok := symbolTable[name]; ok {
        panic(fmt.Sprintf("%v has already been declared\n", name))
      }
      body := definition[1]
      if body[0] == '"' {
        if !(body[0] == '"' && body[len(body)-1] == '"') {
          panic("Invalid definition, it should be of the form:\nname \"def\"")
        }
        symbolTable[name] = body[1:len(body)-1]
      } else if body[0] == '(' {
        fmt.Println("body of compound definition:", body)
        if !(body[0] == '(' && body[len(body)-1] == ')') {
          panic("Invalid definition, it should be of the form:\nname (\"def\" \"def\" ...)")
        }
        //TODO non si può usare Fields perché altrimenti non si possono mettere gli spazi nelle virgolette
        //TODO trovare un modo per parsare " senza fare casini (basta l'escape ma poi i vari join e replace fanno casino)
        splitted_body := strings.Fields(body[1:len(body)-1])
        fmt.Printf("splitted body: %q\n", splitted_body)
        for i, term := range splitted_body {
          fmt.Println("term:", term)
          //TODO ricontrolla questa condizione
          if term[0] == '"' {
            if !(term[0] == '"' && term[len(term)-1] == '"') {
              panic("Invalid definition, it should be of the form:\nname (\"def\" other_name ...)")
            }
            //TODO
          } else {
            if _, ok := symbolTable[term]; !ok {
              panic(fmt.Sprintf("Undeclared %v.\n", term))
            }
            splitted_body[i] = "\""+symbolTable[term]+"\""
          }
        }
        for i, term := range splitted_body {
          if i == len(splitted_body)-1 {
            splitted_body[i] = term[1:len(term)-1]
          } else {
            splitted_body[i] = term[1:]
          }
        }
        stringed_body := strings.Join(splitted_body, "")
        symbolTable[name] = strings.ReplaceAll(stringed_body, "\"", " ")
        fmt.Println("Final definition:", symbolTable[name])
      } else {
        panic("Invalid definition, it should be of the form:\nname \"def\"\nor:\nname: (\"def\" other_name ...)")
      }
      i++
      line = clean_file[i]
    }
  }

  grammar := Grammar{
    R: make(map[NonTerminal][]Production),
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
    if j == i {
      grammar.S = nonTerminal
    }
    grammar.NT = Utils.Union(grammar.NT, []string{nonTerminal})
    fmt.Println(line)
  }
  fmt.Println("Non terminals:", grammar.NT)

  for j := i; j < len(clean_file); j++ {
    line := clean_file[j]
    rule := MakeRule(line)
    grammar.R[rule.A] = append(grammar.R[rule.A], rule.prods...)
    fmt.Printf("%v -> %v\n", rule.A, grammar.R[rule.A])
  }

  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      for _, term := range prod {
        if !grammar.IsNonTerminal(term) && term != EPS {
          grammar.T = Utils.Union(grammar.T, []string{term})
        }
      }
    }
  }
  fmt.Println("Terminals:", grammar.T)

  for _, nt := range grammar.NT {
    grammar.FirstTable[nt] = grammar.First([]string{nt})
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
    R: make(map[NonTerminal][]Production),
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
        grammar.T = Utils.Union(grammar.T, terminals)
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
        grammar.NT = Utils.Union(grammar.NT, nonterminals)
      } else {
        done = true
      }
      i++
    }
  } else {
    panic("Invalid Nonterminals declaration, it should be of the form:\nNT: { <A_0> ... <A_n> }\n(new lines are allowed in the body)")
  }
  if inter := Utils.Intersection(grammar.T, grammar.NT); len(inter) > 0 {
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
        rule := MakeRule(line)
        //fmt.Println(rule)
        grammar.AddRule(rule)
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
        if used := usedTerminals[t]; !used && isStringIn(t, prod) { //TODO
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
    grammar.FirstTable[nt] = grammar.First([]string{nt})
  }
  for _, nt := range grammar.NT {
    grammar.FollowTable[nt] = grammar.Follow(nt)
  }

  return grammar
}

func (G Grammar) AddRule(rule Rule) {
  fmt.Println(rule)
  if G.R == nil {
    G.R = make(map[NonTerminal][]Production)
  }
  if isStringIn(rule.A, G.NT) {
    missing := false
    var missing_term string
    for _, prod := range rule.prods {
      for _, term := range prod {
        if term != EPS && !(isStringIn(term, G.T) || isStringIn(term, G.NT)) {
          missing = true
          missing_term = term
          break
        }
      }
      if missing {
        break
      }
    }
    if !missing {
      //TODO per ora fa solo l'append
      G.R[rule.A] = append(G.R[rule.A], rule.prods...)
    } else {
      panic(fmt.Sprintf("%v is not a term of the grammar (Terminals %v, NonTerminals %v)\n", missing_term, G.T, G.NT))
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
      //TODO
      if len(prod) == 1 && prod[0] == EPS {
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

//TODO 
// - si potrebbe fare che se calcola il first di un nonterminale lo inserisce nella mappa e se lo deve ricacolare, prima di farlo controlla la tabella
func (G Grammar) first(f []string, first []string, nullSyms []string) []string {
  fmt.Println("Calculating first of", f)
  if len(f) == 0 || (len(f) == 1 && (f[0] == EPS || f[0] == "")) {
    return append(first, EPS)
  }
  first_term := string(f[0])
  if isStringIn(first_term, G.T) {
    return append(first, first_term)
  }
  if !isStringIn(first_term, nullSyms) {
    fmt.Println(first_term, "is not nullable")
    for _, prod := range G.R[first_term] {
      if prod[0] != first_term {
        first = G.first(prod, first, nullSyms)
      }
    }
  } else {
    fmt.Println(first_term, "is nullable")
    first_X := make([]string, 0)
    for _, prod := range G.R[first_term] {
      if prod[0] != first_term {
        first_X = G.first(prod, first_X, nullSyms)
      } else if len(prod) > 1 {
        i := 0
        for prod[i] == first_term && i < len(prod) {
          i++
        }
        if i < len(prod) {
          first_X = G.first(prod[i:], first_X, nullSyms)
        }
      } else {
        fmt.Println("C'è un altro caso???")
      }
    }
    fmt.Println(first_X)

    found := false
    first_X_wo_eps := make([]string, 0)
    for _, fi := range first_X {
      if fi != EPS {
        first_X_wo_eps = append(first_X_wo_eps, fi)
      } else {
        found = true
      }
    }

    first_X = first_X_wo_eps
    if found {
      //fmt.Println("without EPS:", first_X)
    }
    if len(first_X) > 0 {
      first = Utils.Union(first, first_X)
    }
    if len(f) > 1 {
      if f[1] != first_term {
        first = Utils.Union(first, G.First(f[1:]))
      }
    } else {
      first = append(first, EPS)
    }
  }
  return first
}

func (G Grammar) First(f []string) []string {
  first := G.first(f, make([]string, 0), G.NullableSymbols())
  //fmt.Println("First con simboli ripetuti:", first)
  //i := 0
  //for i < len(first)-1 {
  //  j := i+1
  //  removed := false
  //  for j < len(first) {
  //    if first[i] == first[j] {
  //      removed = true
  //      first[j] = first[len(first)-1]
  //      first = first[:len(first)-1]
  //    } else {
  //      j++
  //    }
  //  }
  //  if !removed {
  //    i++
  //  }
  //}
  return Utils.Union(first, first)
}

func (G Grammar) follow(Y NonTerminal, from []string) []string {
  if isStringIn(Y, from) {
    return []string{}
  }
  follow := make([]string, 0)
  if Y == G.S {
    follow = append(follow , "$")
  }
  fmt.Printf("Calculating Follow(%v)\n", Y)
  for _, X := range G.NT {
    for _, prod := range G.R[X] {
      for i := range prod {
        term := prod[i]
        if term == Y {
          fmt.Printf("%v in pos %v\n", prod, i)
          beta := make([]string, 0)
          if i != len(prod)-1 {
            beta = prod[i+1:]
          } 
          first := G.First(beta)
          fmt.Printf("First(%v) = %v\n", beta, first)

          if isStringIn(EPS, first) {
            follow = Utils.Union(follow, G.follow(X, append(from, Y)))
          }
          first_wo_eps := make([]string, 0)
          for _, f := range first {
            if f != EPS {
              first_wo_eps = append(first_wo_eps, f)
            }
          }
          fmt.Printf("First(%v) \\ {%v} = %v\n", beta, EPS, first_wo_eps)
          follow = Utils.Union(follow, first_wo_eps)
        }
      }
    } 
  }
  return follow
}

func (G Grammar) Follow(A NonTerminal) []string {
  return G.follow(A, make([]string, 0))
}

func (G Grammar) IsTerminal(X string) bool {
  return isStringIn(X, G.T)
}
func (G Grammar) IsNonTerminal(X string) bool {
  return isStringIn(X, G.NT)
}
