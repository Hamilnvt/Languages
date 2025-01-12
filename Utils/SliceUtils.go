package Utils

func IsInSlice[S comparable](n S, slice []S) bool {
  i := 0
  found := false
  for !found && i < len(slice) {
    if slice[i] == n {
      found = true
    }
    i++
  }
  return found
}

func Reverse[S any](slice []S) []S {
  reversed := make([]S, len(slice))
  copy(reversed, slice)
  for i, j := 0, len(reversed)-1; i < j; i, j = i+1, j-1 {
    reversed[i], reversed[j] = reversed[j], reversed[i]
  }
  return reversed
}

