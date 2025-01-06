S: E

T: {
  a b
  + - *
 ( ) [ ] { \}
}

NT: {
 E F T D A
}

R: {
  E -> TF
  F -> ε | +E | -E
  T -> AD
  D -> ε | *T
  A -> a | b | (E) | [E] | {E}
}
