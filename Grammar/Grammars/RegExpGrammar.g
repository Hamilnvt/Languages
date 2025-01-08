DEFINE:

letter "[a-z]"

GRAMMAR:

S -> \eps | R
R -> A P
P -> \| A P | \eps
A -> B Q
Q -> B | \eps
B -> C T
T -> * T | \eps

C -> a | b | c | ( R )
# C -> a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z | ( R )
# C -> letter | ( R )
