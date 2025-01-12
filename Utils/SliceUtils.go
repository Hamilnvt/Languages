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

