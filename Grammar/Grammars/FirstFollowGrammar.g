GRAMMAR:

S -> B C A | A B E
A -> a | a D b | b S c
B -> C | B b
C -> \eps | d C
D -> d D
E -> D | d E
