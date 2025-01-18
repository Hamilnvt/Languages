DEFINE:

Symbol "[a-z]"

GRAMMAR:

S -> \eps | R
R -> A P
P -> \| A P | \eps
A -> B Q
Q -> B Q | \eps
B -> C Unary
Unary -> * Unary | + Unary | ? Unary | \eps

# sto provando a inserire le parentesi quadre
C -> a | b | c | ( R ) | [ D ]
D -> a | b | c | D - E | \eps
E -> a | b | c

# C -> a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z | ( R )

# in alternativa, ma è molto simile

# RE -> Concat Union | Union
# Union -> \| Concat Union | \eps
# Concat -> B Concat'
# Concat' -> B Concat' | \eps
# B -> Sym Star | ( RE ) B'
# Star -> * Star | \eps
# Sym -> a | b | c
