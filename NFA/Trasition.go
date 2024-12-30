package NFA

import ("fmt")

type Transition struct {
  label string
  src *State
  dst *State
}

func (t Transition) String() (res string) {
  res += fmt.Sprintf("-%v-> %v", t.label, t.dst.label)
  return
}
