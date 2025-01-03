package Grammar

import (
  "fmt"
  "strings"
  _"unicode"
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
  fmt.Println("Creating rule from:", rule_str)
  rule := Rule{} 
  parsed_rule := strings.Fields(rule_str)
  fmt.Printf("%q\n", parsed_rule)
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
      fmt.Printf("%v -> %v\n", rule.A, parsed_rule[i])
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

//TODO non inserisce correttamente i simboli
func MakeGrammar(rules []string, initialSymbol string, nonterminals []NonTerminal, terminals []Terminal) Grammar {
  parsed_rules := make([]Rule, 0)
  for _, rule := range rules {
    parsed_rules = append(parsed_rules, makeRule(rule))
  }
  R := make(map[NonTerminal][]Prod)
  for _, rule := range parsed_rules {
    fmt.Println(rule)
    R[rule.A] = rule.prods
  }
  fmt.Println("R:", R)
  G := Grammar{
    S: parsed_rules[0].A,
    NT: nonterminals,
    T: terminals,
    R: R,
  }
  return G
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
   // Convert the map keys to a sliceq5
   output := make([]string, 0, len(values)) //create slice output
   for val := range values {
      output = append(output, val) //append values in slice output
   }
   return output
}

//TODO 
// - si potrebbe fare che se calcola il first di un nonterminale lo inserisce nella mappa e se lo deve ricacolare, prima di farlo controlla la tabella
/*
ci sono quasi, non ha fatto questa cosa:

B -> C | Bb
C -> eps | d

first(B) dovrebbe essere { eps, b, d } e invece è solo { eps, d }, come se non si fosse annullato (probabilmente dove skippo la produzione uguale al simbolo iniziale se il simbolo è annullabile devo fare attenzione)

*/
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
