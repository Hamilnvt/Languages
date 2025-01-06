package Parsing

import (
  "fmt"
  "strings"
  "Languages/Grammar"
)

type DerivationTree struct {
  grammar Grammar.Grammar
  val string
  children []*DerivationTree
  parent *DerivationTree
}

//TODO
func (tree DerivationTree) String() string {
  return tree.FormattedTree(0)
}

func (tree DerivationTree) FormattedTree(pad int) (res string) {
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

func (tree DerivationTree) isLeaf() bool {
  return len(tree.children) == 0 && tree.grammar.IsTerminal(tree.val)
}

func (tree DerivationTree) getLeaves(leaves []string) []string {
  if len(tree.children) == 0 {
    return append(leaves, tree.val)
  } else {
    for _, child := range tree.children {
      for _, leaf := range child.getLeaves(make([]string, 0)) {
        leaves = append(leaves, leaf)
      }
    }
    return leaves
  }
}

func (tree DerivationTree) GetParsedString() string {
  return strings.Join(tree.getParsedString(make([]string, 0)), "")
}

func (tree DerivationTree) getParsedString(s []string) []string {
  if tree.isLeaf() {
    return append(s, tree.val)
  } else {
    for _, child := range tree.children {
      s = append(s, child.getParsedString(make([]string, 0))...)
    }
  }
  return s
}

func makeDerivationTree(grammar Grammar.Grammar, val string) DerivationTree {
  return DerivationTree{
    grammar: grammar,
    val: val,
    children: make([]*DerivationTree, 0),
    parent: nil,
  }
}

func newDerivationTree(grammar Grammar.Grammar, val string) *DerivationTree {
  tree := new(DerivationTree)
  tree.grammar = grammar
  tree.val = val
  tree.parent = nil
  return tree
}

func (tree *DerivationTree) findLeftMostNonTerminal() *DerivationTree {
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

func (root  *DerivationTree) addChildren(prod string) {
  tree := root.findLeftMostNonTerminal()
  for _, symbol := range prod {
    child := newDerivationTree(tree.grammar, string(symbol))
    child.parent = root
    tree.children = append(tree.children, child)
  }  
}

