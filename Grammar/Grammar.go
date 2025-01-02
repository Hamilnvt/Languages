package Grammar

import (
  "fmt"
  "strings"
  "unicode"
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
      fmt.Printf("%v. %v\n", i, parsed_rule[i])
      rule.prods = append(rule.prods, prod)
    }
  }
  return rule
}

func isStringIn(s string, strs []string) bool {
  found := false
  i := 0
  for !found && i < len(strs) {
    if strings.Compare(s, strs[i]) == 0 {
      found = true
    }
    i++
  }
  return found
}

func MakeGrammar(rules []string) Grammar {
  parsed_rules := make([]Rule, 0)
  for _, rule := range rules {
    parsed_rules = append(parsed_rules, makeRule(rule))
  }
  NT := make([]NonTerminal, 0)
  T := make([]Terminal, 0)
  R := make(map[NonTerminal][]Prod)
  for _, rule := range parsed_rules {
    fmt.Println(rule)
    if !isStringIn(rule.A, NT) {
      NT = append(NT, rule.A)    
    }
    R[rule.A] = rule.prods
    for _, prod := range rule.prods {
      prod = prod 
      for _, c := range prod {
        if (unicode.IsLower(c) || !unicode.IsLetter(c)) && !isStringIn(string(c), T) {
            T = append(T, string(c))    
        } else if (unicode.IsLetter(c) && unicode.IsUpper(c)) && !isStringIn(string(c), NT) {
          NT = append(NT, rule.A)    
        }
      }
    }
  }
  fmt.Println("NT:", NT)
  fmt.Println("T:", T)
  fmt.Println("R:", R)
  G := Grammar{
    S: parsed_rules[0].A,
    NT: NT,
    T: T,
    R: R,
  }
  return G
}

func (G Grammar) String() (res string) {
  fmt.Println("Printing Grammar: TODO")
  return
}
