package Parsing

import (
  "fmt"
  "Languages/Grammar"
  "Languages/Utils"
  "errors"
  "os"
  "text/tabwriter"
  "strings"
)

type TDTablePair struct {
  A Grammar.NonTerminal
  a Grammar.Terminal
}

type TDParsingTable map[TDTablePair][]string

const EPS = Grammar.EPS

type Parser_LL1 struct {
  grammar Grammar.Grammar
  stack Utils.Stack[string]
  input []string
  table TDParsingTable
}

func (parser Parser_LL1) String() (res string) {
  res += "Printing parser:\n"
  res += "TODO"
  return
}

//TODO si può fare di meglio, ma così già ci sta
func (parser Parser_LL1) PrintTable() {
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
      pair := TDTablePair{A:nt, a:t}
      fmt.Fprintf(w, "%v\t|\t", parser.table[pair])
    }
    fmt.Fprintln(w, "")
  }
	w.Flush()
}

func MakeParserTopDownLL1(grammar Grammar.Grammar) (Parser_LL1, error) {
  parser := Parser_LL1{
    grammar: grammar,
    table: make(TDParsingTable),
  }

  //TODO, forse potrei controllare prima se è LL(1) con le intersezioni dei first
  //TODO ricontrollare, forse non funziona. E infatti non tiene conto degli spazi e senza definizioni non legge ad esempio "if", devo fare il LA
  for _, nt := range grammar.NT {
    for _, prod := range grammar.R[nt] {
      first := grammar.First(prod)
      fmt.Printf("First(%v) = %v\n", prod, first)
      for _, term := range first {
        if term == EPS {
          follow := grammar.FollowTable[nt] 
          for _, f_term := range follow {
            pair := TDTablePair{A:nt, a:f_term}
            if val, ok := parser.table[pair]; !ok || len(val) == 0 {
              parser.table[pair] = prod 
            } else {
              fmt.Printf("ERROR:\nTablePair: (%v, %v)\nTrying to insert: %v\nBut there is already: %v\n", pair.A, pair.a, prod, val)
              parser.PrintTable()
              return Parser_LL1{}, errors.New("Grammar is not LL(1) :(")
            } 
          }          
        } else {
          pair := TDTablePair{A:nt, a:term}
          if val, ok := parser.table[pair]; !ok || len(val) == 0 {
            parser.table[pair] = prod 
          } else {
            fmt.Printf("ERROR:\nTablePair: (%v, %v)\nTrying to insert: %v\nBut there is already: %v\n", pair.A, pair.a, prod, val)
            parser.PrintTable()
            return Parser_LL1{}, errors.New("Grammar is not LL(1) :(")
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

func (parser Parser_LL1) Parse(input string) (ParseTree, error) {
  parser.stack = Utils.Stack[string]{}
  parser.stack.Push(parser.grammar.S)
  for _, c := range input {
    parser.input = append(parser.input, string(c))
  }
  parser.input = append(parser.input, "$")
  ic := 0
  fmt.Println("Parsing", strings.Join(parser.input, ""))

  tree := makeParseTree(parser.grammar, parser.grammar.S)
  w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
  for X, err := parser.stack.Top(); err == nil; X, err = parser.stack.Top() {
    fmt.Fprintf(w, "stack: %v\t|  input: %v\t", parser.stack, strings.Join(parser.input, "")[ic:]) 
    if parser.grammar.IsTerminal(X) {
      if X == parser.input[ic] {
        parser.stack.Pop()
        ic++
        fmt.Fprintf(w, "|  match %v\t", X)
      } else {
        fmt.Fprintln(w)
        w.Flush()
        return ParseTree{}, errors.New(fmt.Sprintf("No match for Terminal at index %v, symbol %v.", ic, parser.input[ic]))
      }
    } else {
      prod, ok := parser.table[TDTablePair{A:X, a:parser.input[ic]}]
      //fmt.Fprintf(w, "|  ParsingTable[%v, %v] = %v\t", X, parser.input[ic], prod)
      if ok {
        parser.stack.Pop()
        if len(prod) != 0 && prod[0] != EPS && prod[0] != "" {
          for j := len(prod)-1; j >= 0; j-- {
            parser.stack.Push(string(prod[j]))
          }
        }
        fmt.Fprintf(w, "|  output: %v -> %v", X, prod)
        tree.addChildren(prod)
      } else {
        fmt.Fprintln(w)
        w.Flush()
        return ParseTree{}, errors.New(fmt.Sprintf("No match for NonTerminal at index %v, symbol %v.", ic, parser.input[ic]))
      }
    }
    fmt.Fprintln(w, "\t")
  }
  fmt.Fprintln(w, "stack:\t|  input: $\t|  String accepted!")

  w.Flush()
  return tree, nil
}
