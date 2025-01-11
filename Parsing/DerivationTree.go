package Parsing

import (
  "fmt"
  "strings"
  "Languages/Grammar"
)

type Token[S any] struct {
  pattern string
  name string
  val S
}

type ParseTree struct {
  grammar Grammar.Grammar
  val string
  children []*ParseTree
  parent *ParseTree
}

type AbstractParseTree struct {
  ParseTree
  val Token[any]
}

//TODO
func (tree ParseTree) String() string {
  return tree.FormattedTree(0)
}

func (tree ParseTree) FormattedTree(pad int) (res string) {
  begin := "- "
  if tree.parent == nil {
    begin = ""
  }
  res += fmt.Sprintf("%#v% *v%v\n", pad, pad*2+1, begin, tree.val)
  for _, child := range tree.children {
    res += child.FormattedTree(pad+1)
  }
  return 
}

func (tree ParseTree) isLeaf() bool {
  return len(tree.children) == 0 && tree.grammar.IsTerminal(tree.val)
}

func (tree ParseTree) GetParsedString() string {
  return strings.Join(tree.getParsedString(make([]string, 0)), " ")
}

func (tree ParseTree) getParsedString(s []string) []string {
  if tree.isLeaf() {
    return append(s, tree.val)
  } else {
    for _, child := range tree.children {
      s = append(s, child.getParsedString(make([]string, 0))...)
    }
  }
  return s
}

func makeParseTree(grammar Grammar.Grammar, val string) ParseTree {
  return ParseTree{
    grammar: grammar,
    val: val,
    children: make([]*ParseTree, 0),
    parent: nil,
  }
}

func newParseTree(grammar Grammar.Grammar, val string) *ParseTree {
  tree := new(ParseTree)
  tree.grammar = grammar
  tree.val = val
  tree.parent = nil
  return tree
}

//TODO guarda il video
func (tree ParseTree) Abstract() AbstractParseTree {
  //abs := AbstractParseTree{
  //  grammar: tree.grammar,
  //}
  return AbstractParseTree{}
}

func (tree *ParseTree) findLeftMostNonTerminal() *ParseTree {
  if len(tree.children) == 0 {
    if tree.grammar.IsNonTerminal(tree.val) {
      return tree
    } else {
      return nil
    }
  } else {
    for _, child := range tree.children {
      if leftMost := child.findLeftMostNonTerminal(); leftMost != nil {
        return leftMost
      } 
    }
  }
  return nil
}

func (root  *ParseTree) addChildren(prod []string) {
  tree := root.findLeftMostNonTerminal()
  for _, term := range prod {
    child := newParseTree(tree.grammar, term)
    child.parent = root
    tree.children = append(tree.children, child)
  }  
}

