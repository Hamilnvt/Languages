GRAMMAR:

statement -> if condition then command else command
condition -> true | false
command -> a | b | c

# command -> c | command ; allwhitespaces command

#command -> c next_command
#next_command -> ; command next_command | \eps
