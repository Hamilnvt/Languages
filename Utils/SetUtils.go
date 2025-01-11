package Utils

//func union(slice1 []string, slices ...[]string) []string {
//  // Create a map to store the elements of the union
//  values := make(map[string]bool)
//  for _, key := range slice1 { // for loop used in slice1 to remove duplicates from the values
//    values[key] = true
//  }
//  for _, slice := range slices {
//    for _, key := range slice { // for loop used in slice2 to remove duplicates from the values
//      values[key] = true
//    }
//  }
//  // Convert the map keys to a slice
//  output := make([]string, 0, len(values)) //create slice output
//  for val := range values {
//    output = append(output, val) //append values in slice output
//  }
//  return output
//}
func Union[S comparable](setA, setB []S) []S {
  //fmt.Printf("Union:\n%q\n%q\n", setA, setB)
  values := make(map[S]bool)
  for _, key := range setA {
    values[key] = true
  }
  for _, key := range setB {
    values[key] = true
  }
  output := make([]S, len(values))
  i := 0
  for val := range values {
    output[i] = val
    i++
  }
  return output
}
//func union(setA [][]string, sets ...[]string) [][]string {
//  values := make(map[string]bool)
//  for _, key := range slice1 {
//    values[key] = true
//  }
//  for _, slice := range slices {
//    for _, key := range slice {
//      values[key] = true
//    }
//  }
//  output := make([][]string, 0, len(values))
//  for val := range values {
//    output = append(output, val)
//  }
//  return output
//}

//func Intersection(slice1, slice2 []string) []string {
//  values := make(map[string]bool)
//  for _, key := range slice1 {
//    values[key] = true
//  }
//  output := make([]string, 0, len(values))
//  for _, key := range slice2 {
//    if values[key] {
//      output = append(output, key)
//    }
//  }
//  return output
//}

func Intersection[S comparable](SetA, SetB []S) []S {
  values := make(map[S]bool)
  for _, key := range SetA {
    values[key] = true
  }
  output := make([]S, 0, len(values))
  for i, key := range SetB {
    if values[key] {
      output[i] = key
    }
  }
  return output
}

