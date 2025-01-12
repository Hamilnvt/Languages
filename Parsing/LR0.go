package Parsing

import (
  "fmt"
  "os"
  "strings"
  "text/tabwriter"
  "Languages/NFA"
  "Languages/Grammar"
  "Languages/Utils"
)

type BUTableKey struct {
  state int
  term string
}

type LRAction int
const (
  SHIFT LRAction = iota
  REDUCE
  ACCEPT
  GOTO
  BLANK = -1
)
func (action LRAction) String() string {
  switch action {
    case SHIFT:  return "SHIFT"
    case REDUCE: return "REDUCE"
    case ACCEPT: return "ACCEPT"
    case GOTO:   return "GOTO"
    case BLANK:  return "BLANK"
    default:     return "HUH?"
  }  
}

type BUTableEntry struct {
  action LRAction
  num int
}

func (entry BUTableEntry) String() string {
  if entry.action == BLANK {
    return "[]"
  } else {
    return fmt.Sprintf("[%v %v]", entry.action, entry.num)
  }
}

type BUParsingTable map[BUTableKey]BUTableEntry

func (parser Parser_LR0) PrintTable(grammar *Grammar.Grammar, CA *NFA.CALR0) {
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
  grammar_terms := append(grammar.T, "$")
  grammar_terms = append(grammar_terms, grammar.NT...)
  fmt.Fprint(w, "\t|\t")
  for _, term := range grammar_terms {
    fmt.Fprintf(w, "%v\t|\t", term)
  }
  fmt.Fprintln(w, "")
  for _, state := range CA.States {
    fmt.Fprintf(w, "%v\t|\t", state.Index)
    for _, term := range grammar_terms {
      key := BUTableKey{state:state.Index, term:term}
      fmt.Fprintf(w, "%v\t|\t", parser.table[key])
    }
    fmt.Fprintln(w, "")
  }
	w.Flush()
}

type Parser_LR0 struct {
  table BUParsingTable 
  terms_stack Utils.Stack[string]
  states_stack Utils.Stack[int]
}

type NumberedProdKey struct {
  A Grammar.NonTerminal
  prod string
}

func MakeParserBottomUpLR0(grammar Grammar.Grammar) Parser_LR0 {
  parser := Parser_LR0{
    table: make(BUParsingTable),
    stack: Utils.Stack[string]{},
  }
  CA := NFA.MakeCanonicAutomatonLR0(&grammar)
  fmt.Println(CA)

  grammar_terms := make([]string, 0)
  grammar_terms = append(grammar_terms, grammar.T...)
  grammar_terms = append(grammar_terms, "$")
  grammar_terms = append(grammar_terms, grammar.NT...)
  //fmt.Println("Grammar terms:", grammar_terms)

  numbered_prods := make(map[NumberedProdKey]int)
  counter := 1
  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      fmt.Println(counter, prod)
      key := NumberedProdKey{
        A: nt,
        prod: strings.Join(prod, " "),
      }
      numbered_prods[key] = counter
      counter++
    }
  }
  fmt.Println(numbered_prods)

  for _, state := range CA.States {
    for _, term := range grammar_terms {
      // SHIFT
      //fmt.Println("Shift:")
      if grammar.IsTerminal(term) {
        delta_key := NFA.DeltaKey{
          State: state.Index,
          Term: term,
        } 
        if t, delta_ok := CA.Delta[delta_key]; delta_ok {
          key := BUTableKey{
            state: state.Index,
            term: term,
          }
          if _, ok := parser.table[key]; ok {
            panic("Grammar is not LR(0) :(")
          } else {
            entry := BUTableEntry{
              action: SHIFT,
              num: t,
            }
            parser.table[key] = entry
            fmt.Println(key, "=", entry)
          }
        }
      }

      //TODO per ora prende solo il primo item del tipo A -> a. (ma dovrebbe prenderli tutti)
      // REDUCE
      //fmt.Println("Reduce:")
      if grammar.IsTerminal(term) || term == "$" {
        found := false
        var prod_key NumberedProdKey
        i := 0
        for !found && i < len(state.Items) {
          item := state.Items[i]
          if item.A != "ITLR0" && item.Dot == len(item.Prod) {
            found = true    
            prod_key = NumberedProdKey{
              A: item.A,
              prod: strings.Join(item.Prod, " "),
            }
          } else {
            i++
          }
        }
        if found {
          key := BUTableKey{
            state: state.Index,
            term: term,
          }
          if _, ok := parser.table[key]; ok {
            panic("Grammar is not LR(0) :(")
          } else {
            entry := BUTableEntry{
              action: REDUCE,
              num: numbered_prods[prod_key],
            }
            parser.table[key] = entry
            fmt.Println(key, "=", entry)
          }
        }
      }

      // ACCEPT
      //fmt.Println("Accept")
      if term == "$" {
        found := false
        i := 0
        for !found && i < len(state.Items) {
          item := state.Items[i]
          if item.A == "ITLR0" && len(item.Prod) == 1 && item.Prod[0] == grammar.S && item.Dot == len(item.Prod) {
            found = true    
          } else {
            i++
          }
        }
        if found {
          key := BUTableKey{
            state: state.Index,
            term: "$",
          }
          if _, ok := parser.table[key]; ok {
            panic("Grammar is not LR(0) :(")
          } else {
            entry := BUTableEntry{
              action: ACCEPT,
              num: 0,
            }
            parser.table[key] = entry
            fmt.Println(key, "=", entry)
          }
        }
      }

      // GOTO
      //fmt.Println("Goto")
      if grammar.IsNonTerminal(term) {
        delta_key := NFA.DeltaKey{
          State: state.Index,
          Term: term,
        } 
        if t, delta_ok := CA.Delta[delta_key]; delta_ok {
          key := BUTableKey{
            state: state.Index,
            term: term,
          }
          if _, ok := parser.table[key]; ok {
            panic("Grammar is not LR(0) :(")
          } else {
            entry := BUTableEntry{
              action: GOTO,
              num: t,
            }
            parser.table[key] = entry
            fmt.Println(key, "=", entry)
          }
        }
      }
    }
  }

  for _, state := range CA.States {
    for _, term := range grammar_terms {
      key := BUTableKey{
        state: state.Index,
        term: term,
      }
      if _, ok := parser.table[key]; !ok {
        parser.table[key] = BUTableEntry{action: BLANK, num: -1}
      }
    }
  }

  parser.PrintTable(&grammar, &CA)

  return parser
}

func (parser Parser_LR0) Parse(input string) (ParseTree, error) {
  parser.stack = Utils.Stack[string]{}
  parser.stack.Push(parser.grammar.S)

  return ParseTree{}, nil
}
