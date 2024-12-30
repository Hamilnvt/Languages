package NFA

import (
  "fmt"
  "strings"
  "math"
  "strconv"
)

const (
  EPS = "ε"
  FINAL = true
  NON_FINAL = false
)

type NFA struct {
  Sigma string
  States []State
  InitialState int
  n, m int
}

func (N NFA) GetStatesNum() int {
  return N.n
}

func (N NFA) GetStateByLabel(label string) *State {
  i := 0
  found := false
  var q *State = nil
  for !found && i < len(N.States) {
    q = &N.States[i]
    if strings.Compare(label, q.label) == 0 {
      found = true
    } else {
      i++
    }
  }
  return q
}

func (N *NFA) AddState(label string, isFinal bool) {
  q := State{
    Index: N.n,
    label: label,
    adjac: make([]*Transition, 0),
    isFinal: isFinal,
  }
  N.States = append(N.States, q)
  N.n++
}

func (N *NFA) AddTransition(label string, q1, q2 int) {
  if !strings.Contains(N.Sigma, label) && strings.Compare(label, EPS) != 0 {
    N.Sigma += label
  }
  q := &N.States[q1]
  t := Transition{
    label: label,
    src: q,
    dst: &N.States[q2],
  }
  q.adjac = append(q.adjac, &t)
  N.m++
}

func (N NFA) E_clos(states []int) []int {
  T := make([]int, len(states))
  copy(T, states)
  e_clos := make([]int, len(states))
  copy(e_clos, states)

  for len(T) > 0 {
    q := N.States[T[len(T)-1]]
    T = T[:len(T)-1]
    for _, t := range q.adjac {
      if (strings.Compare(t.label, EPS) == 0) {
        p := t.dst
        j := 0
        found := false
        for !found && j < len(e_clos) {
          if p.Index == e_clos[j] {
            found = true
          } else {
            j++
          }
        }
        if !found {
          T = append(T, p.Index)
          e_clos = append(e_clos, p.Index)
        }
      }
    }
  }

  return e_clos
}

func (N NFA) Move(states []int, a string) []int {
  if strings.Compare(a, EPS) == 0 {
    panic(fmt.Sprintf(a, "cannot be epsilon"))
  }
  move := make([]int, 0)
  for _, i := range states {
    q := N.States[i] 
    //fmt.Println("Considering", q.label, q.Index)
    for _, t := range q.adjac {
      if (strings.Compare(t.label, a) == 0) {
        //fmt.Println("transition", t)
        p := t.dst
        k := 0
        found := false
        for !found && k < len(move) {
          if (p.Index == move[k]) {
            found = true
          } else {
            k++
          }
        }
        if !found {
          //fmt.Println("added", p.Index)
          move = append(move, p.Index)
        }
      }
    }
  } 
  return move
}

func (N NFA) ToDFA() DFA {
  /// Costruzione per sottoinsiemi
  S := MakeNtoDState("Q0", &N)
  S.States = append(S.States, N.InitialState)
  S.States = N.E_clos(S.States)
  if S.ContainsFinalState() {
    S.IsFinal = true
  }
  T := make([]NtoDState, 0) 
  T = append(T, S)
  Deltas := make([]Delta, 0)
  //fmt.Printf("T (%v):\n", len(T))
  //for i, t := range T {
  //  fmt.Printf("%v: %v\n", i, t.Label)
  //}
  //fmt.Println("Deltas:", Deltas)
  //fmt.Println("S:", S)
  //fmt.Printf("The sigma: %s\n", N.Sigma)

  state_counter := 1
  i := 0
  for i < len(T) {
    fmt.Println("iteration", i)
    fmt.Println("T len:", len(T))
    P := &T[i]
    fmt.Println(P)
    P.IsMarked = true
    for _, a := range N.Sigma {
      //fmt.Printf("%v: %v\n", i, string(a))
      R := MakeNtoDState("", &N)
      move := N.Move(P.States, string(a))
      fmt.Printf("Move (%v -%v->): %v\n", P.Label, string(a), move)
      R.States = N.E_clos(move)
      fmt.Println("e-clos:", R.States)
      if R.ContainsFinalState() {
        R.IsFinal = true
      }
      fmt.Println("R finito:\n", R)

      RinT := R.IsIn(T)
      if RinT == -1 {
        fmt.Println("R not in T")
        R.Label = fmt.Sprintf("Q%v", state_counter)
        state_counter++
        T = append(T, R)
      } else {
        R = T[RinT]
        fmt.Println("C'è già e si chiama", R.Label)
      }
      fmt.Println("Sto per creare il delta")
      fmt.Println("P:", P)
      fmt.Println("a:", string(a))
      fmt.Println("R:", R)
      D := Delta{
        P:*P,
        A:string(a),
        R:R,
      }
      fmt.Println(D)
      Deltas = append(Deltas, D)
    }

    i++
  }

  fmt.Println("\nPronti per assemblare il DFA")
  fmt.Println("Sigma:", N.Sigma)
  fmt.Println("T:\n", T)
  fmt.Println("Deltas:\n", Deltas)
  fmt.Println("Stato iniziale:\n", T[0])

  M := DFA{}
  M.Sigma = N.Sigma
  for _, t := range T {
    M.AddState(t.Label, t.IsFinal)
  }
  for _, d := range Deltas {
    p := M.GetStateByLabel(d.P.Label)
    r := M.GetStateByLabel(d.R.Label)
    M.AddTransition(d.A, p.Index, r.Index)
  }
  fmt.Println("M:\n", M)
  return M
}

/*

0 x
1 x x
2 x x x
  1 2 3

*/

type IntPair struct {
  a, b int
}

type StairTableEntry struct {
  IntPair
  mark int
}

//TODO per aumentare la base possibile (oltre 36) si possono introdurre le lettere maiuscole (come?)
func GetPermutationString(sigma string, str_len int, p_n int) string {
  // sostanzialmente:
  /*
se sigma = "ab":
 - con str_len = 1:
   - n = 0 -> "a"
   - n = 1 -> "b"
 - con str_len = 2:
   - n = 0 -> "a"
   - n = 1 -> "b"
   - n = 2 -> "aa"
   - n = 3 -> "ab"
   - n = 4 -> "ba"
   - n = 5 -> "bb"
  */
  // se len(sigma) = 2 -> base 2 (0 -> a, 1 -> b)
  // str_len = 3 -> 3 cifre
  // p_n = 3 = 011 in base 2 (p_n non può essere maggiore del numero totale di permutazioni per quella lunghezza)
  // w = "abb", ovvero la permutazione numero 3 (la quarta)

  base := len(sigma)
  if base > 36 {
    panic(fmt.Sprintf("len(sigma) cannot be greater than 36 (is %v)", base))
  }
  //fmt.Printf("%v -> base %v\n", sigma, base)

  if max_n := int(math.Pow(float64(base), float64(str_len)))-1; p_n > max_n {
    panic(fmt.Sprintf("p_n cannot be higher than len(sigma)^(str_len)-1\n(p_n = %v, max_n = %v^(%v)-1 = %v)", p_n, base, str_len, max_n))
  }

  converted := fmt.Sprintf("%0*v", str_len, strconv.FormatInt(int64(p_n), base))
  fmt.Printf("%v is %v in base %v with %v digits\n", p_n, converted, base, str_len)

  res := ""
  for _, c := range converted {
    i, err := strconv.ParseInt(string(c), base, 0)
    if err != nil {
      panic(err)
    }
    fmt.Printf("%v -> %v\n", string(c), string(sigma[i]))
    res += string(sigma[i])
  }

  return res
}

func (N DFA) Minimize() {
  table := make(map[IntPair]StairTableEntry)
  for i := 0; i < len(N.States)-1; i++ {
    for j := 1; j < len(N.States); j++ {
      if i == j {
        continue
      }
      p := IntPair{i, j}
      if N.States[i].isFinal && !N.States[j].isFinal || !N.States[i].isFinal && N.States[j].isFinal {
        table[p] = StairTableEntry{p, 0}
      } else {
        table[p] = StairTableEntry{p, -1}
      }
    }
  }
  fmt.Println("Initialized table:", table)

  i := 1
  done := false
  for !done {
    done = true
    for _, pair := range table {
      if pair.mark == -1 {
        found := false
        str_len := 1
        n := 0
        for !found && str_len <= i {
          w := GetPermutationString(N.Sigma, str_len, n)

          q1 := N.States[pair.a]
          q2 := N.States[pair.b]
          //fmt.Println("Stato attuale:", q1, q2, w)
          
          var p1, p2 *State
          for k := 0; k < len(q1.adjac); k++ {
            //fmt.Println("Cercando in q1", k, q1.adjac[k])
            if strings.Compare(w, q1.adjac[k].label) == 0 {
              p1 = q1.adjac[k].dst
              break
            }
          }
          for k := 0; k < len(q2.adjac); k++ {
            if strings.Compare(w, q2.adjac[k].label) == 0 {
              p2 = q2.adjac[k].dst
              break
            }
          }
          //fmt.Println("Dopo la transizione w:", p1, p2)
          if p1 == nil || p2 == nil || p1.Index == p2.Index {
            fmt.Println("Guarda zio, c'è stato un problema, sono uguali o uno dei due è nil")
          } else {
            //TODO dovrei probabilmente controllare anche la coppia (p2, p1) ?
            new_pair, ok := table[IntPair{p1.Index, p2.Index}]
            if ok {
              fmt.Println("New pair", new_pair)
              if new_pair.mark != -1 {
                pair.mark = i
                fmt.Println("Pair marked", pair)
                done = true
              }
            }
          }
           
          //TODO togli
          found = true
        }
      }
    }
  }
}

func (N NFA) String() (res string) {
  res = fmt.Sprintf("Printing NFA (%v states, %v transitions):\n", N.n, N.m)
  res += fmt.Sprintf("- Sigma: %v\n", N.Sigma)
  res += "- States:\n"
  for _, q := range N.States {
    if q.Index == N.InitialState {
      res += "Initial "
    }
    res += q.String()
  }
  return 
}

// NFA from regexp a*
func MakeNFA_star_example() NFA {
  N := NFA{}

  var q_first, q_pre_add, q_post_add, q_final int

  q_first = N.GetStatesNum()
  N.AddState(fmt.Sprintf("q%v", q_first), NON_FINAL)

  q_pre_add = N.GetStatesNum()
  N.AddState(fmt.Sprintf("q%v", q_pre_add), NON_FINAL)

  q_post_add = N.GetStatesNum()
  N.AddState(fmt.Sprintf("q%v", q_post_add), NON_FINAL)

  q_final = N.GetStatesNum()
  N.AddState(fmt.Sprintf("q%v", q_final), FINAL)

  N.AddTransition(EPS, q_first,    q_pre_add)
  N.AddTransition(EPS, q_first,    q_final)
  N.AddTransition("a",     q_pre_add,  q_post_add)
  N.AddTransition(EPS, q_post_add, q_pre_add)
  N.AddTransition(EPS, q_post_add, q_final)
  fmt.Println(N)

  return N
}

func makeNFA_minimize_example() NFA {
  N := NFA{}

  var A, B, C, D, E, F int

  A = N.GetStatesNum()
  N.AddState("A", NON_FINAL)
  B = N.GetStatesNum()
  N.AddState("B", NON_FINAL)
  F = N.GetStatesNum()
  N.AddState("F", NON_FINAL)

  C = N.GetStatesNum()
  N.AddState("C", FINAL)
  D = N.GetStatesNum()
  N.AddState("D", FINAL)
  E = N.GetStatesNum()
  N.AddState("E", FINAL)

  N.AddTransition("0", A, B)
  N.AddTransition("1", A, C)
  N.AddTransition("0", B, A)
  N.AddTransition("1", B, D)
  N.AddTransition("0", C, E)
  N.AddTransition("1", C, F)
  N.AddTransition("0", D, E)
  N.AddTransition("1", D, F)
  N.AddTransition("0", E, E)
  N.AddTransition("1", E, F)
  N.AddTransition("0", F, F)
  N.AddTransition("1", F, F)

  N.InitialState = A

  return N
}
