package NFA

import (
  "fmt"
)

type DFA = NFA

type NtoDState struct {
  Label string
  States []int
  IsFinal bool
  IsMarked bool
  Nfa *NFA
}

func MakeNtoDState(label string, N *NFA) NtoDState {
  Q := NtoDState{
    Label: label,
    States: make([]int, 0),
    IsFinal: false,
    IsMarked: false,
    Nfa: N,
  }
  return Q
}

func (Q NtoDState) String() (res string) {
  var isFinal_str string
  if Q.IsFinal {
    isFinal_str = "final"
  } else {
    isFinal_str = "non-final"
  }
  var isMarked_str string
  if Q.IsMarked {
    isMarked_str = "marked"
  } else {
    isMarked_str = "unmarked"
  }
  res += fmt.Sprintf("Printing NtoDState %v (%v states) [%v] {%v}:\n", Q.Label, len(Q.States), isFinal_str, isMarked_str)
  if len(Q.States) > 0 {
    for _, q := range Q.States {
      res += fmt.Sprintf("- State %v (%v)\n", Q.Nfa.States[q].label, Q.Nfa.States[q].Index)
      //res += Q.Nfa.States[q].String()
    }
  }
  return
}

func (Q NtoDState) ContainsFinalState() (found bool) {
  i := 0
  found = false
  for !found && i < len(Q.States) {
    q := Q.Nfa.States[Q.States[i]]
    if q.isFinal {
      found = true;
    }
    i++
  }
  return 
}

func isStateIn(r int, Q NtoDState) bool {
  i := 0
  found := false
  for !found && i < len(Q.States) {
    q := Q.States[i]
    if q == r {
      found = true
    } else {
      i++
    }
  }
  return found
}

func (R NtoDState) IsIn(T []NtoDState) int {
  i := 0
  found := false
  for !found && i < len(T) {
    Q := T[i]
    same := true
    j := 0
    for same && j < len(R.States) {
      r := R.States[j]
      same = same && isStateIn(r, Q)
      j++
    }
    if same {
      found = true
    } else {
      i++
    }
  }
  if found {
    return i
  } else {
    return -1
  }
}

type Delta struct {
  P NtoDState
  A string
  R NtoDState
}

func (D Delta) String() string {
  return fmt.Sprintf("Delta %v -%v-> %v\n", D.P.Label, D.A, D.R.Label)
}
