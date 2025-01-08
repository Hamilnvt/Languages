DEFINE:

a "aaa"
b "bb"
c (a "c" b)
# d (c)
# e (b b "a b c")

GRAMMAR:
S -> a S b | \eps | A
A -> c
