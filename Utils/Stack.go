package Utils

import (
  "fmt"
  "errors"
)

type Stack[S any] struct {
  stack []S
}
func (stack Stack[S]) String() (res string) {
  for i := len(stack.stack)-1; i >= 0; i-- {
    res += fmt.Sprintf("%v ", stack.stack[i])
  }
  return
}
func (stack Stack[S]) Top() (S, error) {
  if stack.isEmpty() {
    var empty S
    return empty, errors.New("Empty stack")
  } else {
    return stack.stack[len(stack.stack)-1], nil
  }
}
func (stack Stack[S]) isEmpty() bool {
  return len(stack.stack) == 0
}
func (stack *Stack[S]) Push(elt S) {
  stack.stack = append(stack.stack, elt)
}
func (stack *Stack[S]) Pop() (S, error) {
  if stack.isEmpty() {
    var empty S
    return empty, errors.New("ERROR: Cannot pop from empty stack")
  } else {
    elt := stack.stack[len(stack.stack)-1]
    stack.stack = stack.stack[:len(stack.stack)-1]
    return elt, nil
  }
}
