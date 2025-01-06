S: S

T: {
 a b c d e f g h i j k l m n o p q r s t u v w x y z ( )
 | *
}

NT: {
  S R P A Q B T C
}

R: {
  S -> ε | R
  R -> AP
  P -> \|AP | ε
  A -> BQ
  Q -> B | ε
  B -> CT
  T -> *T | ε
  C -> a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z | (R)
}
