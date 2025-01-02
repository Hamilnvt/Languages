package NFA

import ("fmt")

type State struct {
  Index int
  label string
  adjac map[string][]*State
  isFinal bool
}

func (q State) String() (res string) {
  var isFinal_str string
  if q.isFinal {
    isFinal_str = "final"
  } else {
    isFinal_str = "non-final"
  }
  res += fmt.Sprintf("State %v (%v) [%v]:\n", q.label, q.Index, isFinal_str)
  if len(q.adjac) > 0 {
    for a, P := range q.adjac {
      for _, p := range P {
        res += fmt.Sprintf(" -%v-> %v\n", a, p.label)
      }
    }
  } else {
    res += "No transitions for this state\n"
  }
  return 
}
