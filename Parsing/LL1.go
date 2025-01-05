package Parsing

import (
  "fmt"
  "Languages/Grammar"
  "errors"
  "os"
  "text/tabwriter"
  "strings"
)

type TablePair struct {
  A Grammar.NonTerminal
  a Grammar.Terminal
}

type ParsingTable map[TablePair]string

const EPS = Grammar.EPS

type LL1 struct {
  grammar Grammar.Grammar
  stack Stack
  input []string
  table ParsingTable
}

func (parser LL1) String() (res string) {
  res += "Printing parser:\n"
  res += "TODO"
  return
}

//TODO si può fare di meglio, ma così già ci sta
func (parser LL1) PrintTable() {
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
  terminals := append(parser.grammar.T, "$")
  fmt.Fprint(w, "\t|\t")
  for _, t := range terminals {
    fmt.Fprintf(w, "%v\t|\t", t)
  }
  fmt.Fprintln(w, "")
  for _, nt := range parser.grammar.NT {
    fmt.Fprintf(w, "%v\t|\t", nt)
    for _, t := range terminals {
      pair := TablePair{A:nt, a:t}
      fmt.Fprintf(w, "%v\t|\t", parser.table[pair])
    }
    fmt.Fprintln(w, "")
  }
	w.Flush()
}

func MakeParserTopDownLL1(grammar Grammar.Grammar) (LL1, error) {
  parser := LL1{
    grammar: grammar,
    stack: Stack{},
    input: make([]string, 0),
    table: make(ParsingTable),
  }

  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      first := grammar.First(prod)
      //fmt.Printf("First(%v) = %v\n", prod, first)
      for _, symbol := range first {
        if symbol == EPS {
          follow := grammar.FollowTable[nt] 
          for _, f_symbol := range follow {
            pair := TablePair{A:nt, a:f_symbol}
            if _, ok := parser.table[pair]; !ok {
              parser.table[pair] = prod 
            } else {
              return LL1{}, errors.New("Grammar is not LL(1) :(")
            } 
          }          
        } else {
          pair := TablePair{A:nt, a:symbol}
          if _, ok := parser.table[pair]; !ok {
            parser.table[pair] = prod 
          } else {
            return LL1{}, errors.New("Grammar is not LL(1) :(")
          } 
        }
      }
    }
  } 
  //fmt.Println("\nParsing Table:\n")
  parser.PrintTable()
  //fmt.Println()

  return parser, nil
}

func (parser LL1) Parse(input string) (DerivationTree, error) {
  parser.stack = Stack{}
  parser.stack.Push(parser.grammar.S)
  for _, c := range input {
    parser.input = append(parser.input, string(c))
  }
  parser.input = append(parser.input, "$")
  ic := 0
  fmt.Println("Parsing", strings.Join(parser.input, ""))

  //TODO stampa cose strane (simboli non ASCII)
  for X := parser.stack.Top(); X != nil; X = parser.stack.Top() {
    fmt.Print(*X, "\t", strings.Join(parser.input, "")[ic:])
    output := false
    if parser.grammar.IsTerminal(*X) {
      if *X == parser.input[ic] {
        parser.stack.Pop()
        ic++
      } else {
        fmt.Println()
        return DerivationTree{}, errors.New(fmt.Sprintf("No match for Terminal at index %v, symbol %v.", ic, parser.input[ic]))
      }
    } else {
      prod, ok := parser.table[TablePair{A:*X, a:parser.input[ic]}]
      //fmt.Printf("\nTable[%v, %v] = %v\n", *X, parser.input[ic], prod)
      if ok {
        parser.stack.Pop()
        for j := len(prod)-1; j >= 0; j-- {
          parser.stack.Push(string(prod[j]))
        }
        output = true
        fmt.Printf("\t%v -> %v\n", *X, prod)
      } else {
        fmt.Println()
        return DerivationTree{}, errors.New(fmt.Sprintf("No match for NonTerminal at index %v, symbol %v.", ic, parser.input[ic]))
      }
    }
    if !output {
      fmt.Println()
    }
  }

  return DerivationTree{}, errors.New("TODO")
}

type Stack struct {
  stack []string
}
func (stack Stack) String() (res string) {
  for i := len(stack.stack)-1; i >= 0; i-- {
    res += fmt.Sprintf("%v ", stack.stack[i])
  }
  return
}
func (stack Stack) Top() *string {
  length := len(stack.stack)
  if length == 0 {
    return nil
  } else {
    return &stack.stack[length-1]
  }
}
func (stack Stack) isEmpty() bool {
  return len(stack.stack) > 0
}
func (stack *Stack) Push(elt string) {
  stack.stack = append(stack.stack, elt)
}
func (stack *Stack) Pop() (string, error) {
  length := len(stack.stack)
  if length == 0 {
    return "", errors.New("ERROR: Cannot pop from empty stack")
  } else {
    elt := stack.stack[length-1]
    stack.stack = stack.stack[:length-1]
    return elt, nil
  }
}
