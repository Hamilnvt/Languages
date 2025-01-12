package Parsing

import (
  "fmt"
  "os"
  "strings"
  "errors"
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

type NumberedProdKey struct {
  A Grammar.NonTerminal
  prod string
}

type NumberedProdTable map[NumberedProdKey]int

func (table NumberedProdTable) getProd(num int) (NumberedProdKey, error) {
  for key, value := range table {
    if value == num {
      return key, nil
    }
  }
  var empty NumberedProdKey
  return empty, errors.New("There isn't a production associated with this number")
}

type Parser_LR0 struct {
  table BUParsingTable 
  terms_stack Utils.Stack[string]
  states_stack Utils.Stack[int]
  input []string
  numbered_prods NumberedProdTable
}

func MakeParserBottomUpLR0(grammar Grammar.Grammar) Parser_LR0 {
  parser := Parser_LR0{
    table: make(BUParsingTable),
    numbered_prods: make(map[NumberedProdKey]int),
  }

  counter := 1
  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      fmt.Println(counter, prod)
      key := NumberedProdKey{
        A: nt,
        prod: strings.Join(prod, " "),
      }
      parser.numbered_prods[key] = counter
      counter++
    }
  }
  fmt.Println(parser.numbered_prods)

  CA := NFA.MakeCanonicAutomatonLR0(&grammar)
  fmt.Println(CA)

  grammar_terms := make([]string, 0)
  grammar_terms = append(grammar_terms, grammar.T...)
  grammar_terms = append(grammar_terms, "$")
  grammar_terms = append(grammar_terms, grammar.NT...)
  //fmt.Println("Grammar terms:", grammar_terms)

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
              num: parser.numbered_prods[prod_key],
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
  parser.terms_stack = Utils.Stack[string]{}

  parser.states_stack = Utils.Stack[int]{}
  parser.states_stack.Push(0)

  for _, c := range input {
    parser.input = append(parser.input, string(c))
  }
  parser.input = append(parser.input, "$")
  ic := 0
  fmt.Printf("Parsing '%v'\n", input)

  var current_state int
  var err error
  accepted := false
  w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
  defer w.Flush()
  fmt.Fprintln(w, "STATES\t|  TERMS\t|  INPUT\t|  ACTION\t|  OUTPUT")
  for !accepted {
    //output := false
    current_state, err = parser.states_stack.Top()
    fmt.Fprintf(w, "%v\t|  %v\t|  %v\t|  ", parser.states_stack, parser.terms_stack, strings.Join(parser.input, "")[ic:]) 
    if err != nil {
      panic("Input didn't match (Empty stack)")
    }
    key := BUTableKey{
      state: current_state,
      term: parser.input[ic],
    }
    cell := parser.table[key]
    //fmt.Println(key, cell)
    switch cell.action {
      case ACCEPT:
        fmt.Fprintln(w, "ACCEPT\t|  String accepted!")
        accepted = true
      case SHIFT:
        fmt.Fprintf(w, "SHIFT %v\t|\n", cell.num)
        parser.states_stack.Push(cell.num)
        parser.terms_stack.Push(parser.input[ic])
        ic++
      case REDUCE:
        prod_key, err := parser.numbered_prods.getProd(cell.num)
        if err != nil {
          panic(err)
        }
        //fmt.Println("Prod key:", prod_key)
        prod := Grammar.Production(strings.Split(prod_key.prod, " "))
        //fmt.Println("Prod:", prod)
        popped_prod := make(Grammar.Production, len(prod))
        for i := 0; i < len(prod); i++ {
          _, err_s := parser.states_stack.Pop()
          if err_s != nil {
            panic("Input didn't match (Empty stack)")
          }
          term, err_t := parser.terms_stack.Pop()
          if err_t != nil {
            panic("Input didn't match (Empty stack)")
          }
          popped_prod[i] = term
        }
        popped_prod = popped_prod.Reverse()
        //fmt.Println("Popped prod:", popped_prod)
        if !prod.Equals(popped_prod) {
          panic("Input didn't match (Prods are not equal)")
        }
        from_state, err_top := parser.states_stack.Top() 
        if err_top != nil {
          panic("Input didn't match (Empty stack)")
        }
        goto_key := BUTableKey{
          state: from_state,
          term: prod_key.A,
        }
        goto_state := parser.table[goto_key]
        if goto_state.action != GOTO {
          panic("Input didn't match (Mismatch cell action: want GOTO)")
        }
        parser.states_stack.Push(goto_state.num)
        parser.terms_stack.Push(prod_key.A)
        fmt.Fprintf(w, "REDUCE %v\t|  Output: %v -> %v\n", cell.num, prod_key.A, strings.Join(prod, ""))

      default:
        panic("Input  didn't match (Blank cell or something else)")
    }
  }

  return ParseTree{}, nil
}
