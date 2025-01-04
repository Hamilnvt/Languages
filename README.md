Work in progress

# Parsing Grammar from file

You can parse a grammar by invoking the function Grammar.ParseGrammar("/path/to/grammar.g")

Any file .g should follow this syntax:

```
# This is a comment and can be placed only at the beginning of a line (for the moment)
# Blank lines will be ignored

S: <initial_symbol>

# List of Terminals
T: {
  a_00 ... a_0k
  a_10 ...
  .
  .
  .
  a_i0 ... a_ik
}

# white spaces are the separators, so you can use every character you want
# except for '}' and '#' which you'll need to escape like this \} and \#

# List of NonTerminals
NT: {
  A_00 ... A_0k
  A_10 ...
  .
  .
  .
  A_i0 ... A_ik
}

# List of Rules (Spaces between right productions are important!)
R: {
  A_0 -> b_00 | ... | b_0k
  A_1 -> ...
  .
  .
  .
  A_i -> b_i0 | ... | b_ik
}
```
Anything other than that will throw a panic error at your face.
