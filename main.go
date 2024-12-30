package main

import (
  "fmt"
  "NFA/NFA"
)

func main() {

}

// TESTS (forse farò un file separato ad un certo punto)
func test_minimize() {
  N := NFA.MakeNFA_minimize_example()
  fmt.Println(N)

  N.Minimize()
  fmt.Println(N)
}

func test_permutationString() {
  sigma := "abc"
  max_len := 3
  p_n := 3
  w := NFA.GetPermutationString(sigma, max_len, p_n)
  fmt.Println("w:", w)
}
