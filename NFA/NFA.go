package NFA

import (
  "fmt"
  "strings"
  "math"
  "strconv"
  "testing"
  "Languages/Grammar"
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

func (N *NFA) AddState(label string, isFinal bool) *State {
  q := State{
    Index: N.n,
    label: label,
    adjac: make(map[string][]*State),
    isFinal: isFinal,
  }
  N.States = append(N.States, q)
  N.n++

  return &N.States[q.Index]
}

func (N *NFA) removeState(label string) {
  //TODO
  // - cambiare gli indici (è molto lungo perché devo cambiare anche tutte le transizioni)
  //initial_state := N.States[N.InitialState]
  //i := 0
  stateToRemove := N.GetStateByLabel(label).Index
  if stateToRemove == N.InitialState {
    panic("Can't remove initial state")
  }
  N.States[stateToRemove] = N.States[len(N.States)-1]
  N.States = N.States[:len(N.States)-1]
  N.n--
  //TODO non funziona così, dovrò cambiare gli indici
}

func (N *NFA) AddTransition(label string, q1, q2 int) {
  if !strings.Contains(N.Sigma, label) && strings.Compare(label, EPS) != 0 {
    N.Sigma += label
  }
  q := &N.States[q1]
  p := &N.States[q2]
  if q.adjac[label] == nil {
    q.adjac[label] = make([]*State, 0)
  }
  q.adjac[label] = append(q.adjac[label], p)
  N.m++
}

func (N *NFA) removeTransition(label string, q1, q2 int) {
  //fmt.Printf("Removing transition: %v -%v-> %v\n", N.States[q1].label, label, N.States[q2].label)
  T := N.States[q1].adjac[label]
  for i, p := range T {
    if p.Index == q2 {
      T[i] = T[len(T)-1]
      N.States[q1].adjac[label] = T[:len(T)-1]
      N.m--
    }
    if len(T) == 0 {
      delete(N.States[q1].adjac, label)
    }
  }
}

// ritorna *State, State o int?
func (N DFA) Delta(q int, w string) int {
  //fmt.Printf("Delta(%v, '%v')\n", q, w)
  currentState := &N.States[q]
  for _, a := range w {
    adjac, ok := currentState.adjac[string(a)]
    //fmt.Println("nextState:", adjac)
    if !ok || len(adjac) == 0 {
      return -1
    } else {
      currentState = adjac[0]
    }
  }
  return currentState.Index
}

func (N NFA) E_clos(states []int) []int {
  T := make([]int, len(states))
  copy(T, states)
  e_clos := make([]int, len(states))
  copy(e_clos, states)

  for len(T) > 0 {
    q := N.States[T[len(T)-1]]
    T = T[:len(T)-1]
    for label, P := range q.adjac {
      for _, p := range P {
        if (strings.Compare(label, EPS) == 0) {
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
    for label, P := range q.adjac {
      for _, p := range P {
        if (strings.Compare(label, a) == 0) {
          //fmt.Println("transition", t)
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
  } 
  return move
}

// TODO ricontrollare, magari ora si può migliorare
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
  for i := 0; i < len(T); i++ {
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
  }

  fmt.Println("\nPronti per assemblare il DFA")
  fmt.Println("Sigma:", N.Sigma)
  fmt.Println("T:\n", T)
  fmt.Println("Deltas:\n", Deltas)
  fmt.Println("Stato iniziale:\n", T[0])

  M := DFA{}
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
  //fmt.Printf("%v is %v in base %v with %v digits\n", p_n, converted, base, str_len)

  res := ""
  for _, c := range converted {
    i, err := strconv.ParseInt(string(c), base, 0)
    if err != nil {
      panic(err)
    }
    //fmt.Printf("%v -> %v\n", string(c), string(sigma[i]))
    res += string(sigma[i])
  }

  return res
}

func (N NFA) Copy() NFA {
  M := NFA{}
  for _, state_n := range N.States {
    M.AddState(state_n.label, state_n.isFinal)
  }
  M.InitialState = M.GetStateByLabel(N.States[N.InitialState].label).Index
  for _, state_n := range N.States {
    state_m := M.GetStateByLabel(state_n.label)
    for a, delta := range state_n.adjac {
      for _, p_n := range delta {
        p_m := M.GetStateByLabel(p_n.label)
        M.AddTransition(a, state_m.Index, p_m.Index)
      }
    }
  }
  return M
}

func (N DFA) Minimize() DFA {
  table := make(map[IntPair]StairTableEntry)
  pairs := make([]IntPair, 0)
  //fmt.Println("Table Initialization")
  for i := 1; i < len(N.States); i++ {
    for j := 0; j < len(N.States)-1; j++ {
      if i != j {
        _, ok := table[IntPair{j, i}]
        var p IntPair
        if !ok {
          p = IntPair{i, j}
          if N.States[i].isFinal && !N.States[j].isFinal || !N.States[i].isFinal && N.States[j].isFinal {
            table[p] = StairTableEntry{p, 0}
            //fmt.Printf("(%v, %v): %v\n", p.a, p.b, 0)
          } else {
            table[p] = StairTableEntry{p, -1}
            pairs = append(pairs, p)
            //fmt.Printf("(%v, %v): %v\n", p.a, p.b, -1)
          }
        }
      }
    }
  }
  //fmt.Println("Pairs to mark")
  //for _, p := range pairs {
  //  fmt.Println(p)
  //}

  i := 1
  done := false
  for !done && len(pairs) > 0 {
    //fmt.Println("\nIteration", i)
    done = true
    pair_index := 0
    for pair_index < len(pairs) {
      pair := table[pairs[pair_index]]
      //fmt.Printf("\nConsidering pair (%v, %v)\n", pair.a, pair.b)
      found := false
      str_len := 1
      n := 0
      for !found && str_len <= i {
        perm_n := int(math.Pow(float64(len(N.Sigma)), float64(str_len)))
        //fmt.Printf("P(%v, %v) = %v, n=%v\n", len(N.Sigma), str_len, perm_n, n)
        for !found && n < perm_n {
          w := GetPermutationString(N.Sigma, str_len, n)
          //fmt.Printf("'%v', len=%v, n=%v, w='%v'\n", N.Sigma, str_len, n, w)

          p1 := N.Delta(pair.a, w)
          p2 := N.Delta(pair.b, w)
          if p1 == -1 || p2 == -1 {
            panic(fmt.Sprintf("There isn't a transition %v -%v-> %v\n", pair.a, w, pair.b))
          }
          //fmt.Printf("(%v, %v) -%v-> (%v, %v)\n", pair.a, pair.b, w, p1, p2)

          if p1 != p2 {
            delta_pair, ok := table[IntPair{p1, p2}]
            delta_pair_inverted, ok_inverted := table[IntPair{p2, p1}]
            var actual_delta_pair StairTableEntry
            if ok {
              actual_delta_pair = delta_pair 
            } else if ok_inverted {
              actual_delta_pair = delta_pair_inverted
            }
            //fmt.Println("delta pair:", actual_delta_pair)
            if actual_delta_pair.mark != -1 && actual_delta_pair.mark < i {
              pair.mark = i
              table[IntPair{pair.a, pair.b}] = pair
              //fmt.Println("Marked at", i)
              found = true
              pairs[pair_index] = pairs[len(pairs)-1]
              pairs = pairs[:len(pairs)-1]
              done = false

            } else {
              //fmt.Println("Can't be marked at", i)
            }
          } else {
            //fmt.Println("They're the same state")
          }
          n++
        }
        str_len++
      }
      pair_index++
    }
    //fmt.Println("\nTable after iteration", i)
    //for _, pair := range table {
    //  fmt.Printf("(%v, %v): %v\n", pair.a, pair.b, pair.mark)
    //}
    //fmt.Println("Pairs to mark after iteration", i)
    //for _, pair := range pairs {
    //  fmt.Println(pair)
    //}
    i++
  }
  //fmt.Println("Table")
  //for _, pair := range table {
  //  fmt.Printf("(%v, %v): %v\n", pair.a, pair.b, pair.mark)
  //}
  //fmt.Println("Stati indistinguibili")
  //for _, pair := range pairs {
  //  fmt.Println(pair)
  //}

  statesToFuse := make([][]int, 0)
  for _, pair := range pairs {
    if len(statesToFuse) == 0 {
      new_class := make([]int, 0)
      new_class = append(new_class, pair.a, pair.b)
      statesToFuse = append(statesToFuse, new_class)
    } else {
      found := false
      for i, class := range statesToFuse {
        if isInSlice(pair.a, class) || isInSlice(pair.b, class) {
          if isInSlice(pair.a, class) && !isInSlice(pair.b, class) {
            statesToFuse[i] = append(class, pair.b)
          } else if !isInSlice(pair.a, class) && isInSlice(pair.b, class) {
            statesToFuse[i] = append(class, pair.a)
          }
          found = true
          break
        }
      }
      if !found {
        new_class := make([]int, 0)
        new_class = append(new_class, pair.a, pair.b)
        statesToFuse = append(statesToFuse, new_class)
      }
    } 
  }  
  //fmt.Println("Classi di equivalenza in N:", statesToFuse)
  
  N_min := N.Copy()
  for i, class := range statesToFuse {
    new_class := make([]int, 0)
    for _, state := range class {
      new_class = append(new_class, N_min.GetStateByLabel(N.States[state].label).Index)
    }
    statesToFuse[i] = new_class
  }
  //fmt.Println("Classi di equivalenza in N_min:", statesToFuse)

  initialStateHasChanged := false
  for _, class := range statesToFuse {
    new_state_final := false
    for _, state := range class {
      if N_min.States[state].isFinal {
        new_state_final = true
        break
      }
    }
    new_state := N_min.AddState(fmt.Sprintf("q%v", N_min.n), new_state_final)
    if isInSlice(N_min.InitialState, class) {
      if initialStateHasChanged {
        panic("NFA can't have more than one initial state")
      } else {
        N_min.InitialState = new_state.Index
        initialStateHasChanged = true
      }
    }
    //fmt.Println("New state:", new_state)
    for _, q := range N_min.States {
      if !isInSlice(q.Index, class) {
        for _, a := range N_min.Sigma {
          d := N_min.Delta(q.Index, string(a))
          if d == -1 {
            if isInSlice(d, class) {
              N_min.removeTransition(string(a), q.Index, d)
              N_min.AddTransition(string(a), q.Index, new_state.Index)
            }
          }
        }
      }
    }
    for _, q := range class {
      for _, a := range N_min.Sigma {
        d := N_min.Delta(q, string(a))
        if d != -1 {
          N_min.removeTransition(string(a), q, d)
          _, ok := N_min.States[new_state.Index].adjac[string(a)]
          if !ok {
            if isInSlice(d, class) {
              N_min.AddTransition(string(a), new_state.Index, new_state.Index)
            } else {
              N_min.AddTransition(string(a), new_state.Index, d)
            }
          }
        }
      }
    }
  }
  for _, class := range statesToFuse {
    for _, state := range class {
      N_min.removeState(N_min.States[state].label)
    }
  }
  return N_min
}

func isInSlice[S comparable](n S, slice []S) bool {
  i := 0
  found := false
  for !found && i < len(slice) {
    if slice[i] == n {
      found = true
    }
    i++
  }
  return found
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

type ItemLR0 struct {
  A Grammar.NonTerminal
  Prod Grammar.Production
  Dot int
}

func (item ItemLR0) String() string {
  var prod_with_dot string
  if len(item.Prod) == 0 {
    prod_with_dot = "."
  } else if item.Dot == len(item.Prod) {
    prod_with_dot = strings.Join(item.Prod, " ")+" ."
  } else if item.Dot == 0 {
    prod_with_dot = strings.Join(item.Prod[:item.Dot], " ")+". "+strings.Join(item.Prod[item.Dot:], " ")
  } else {
    prod_with_dot = strings.Join(item.Prod[:item.Dot], " ")+" . "+strings.Join(item.Prod[item.Dot:], " ")
  }
  return fmt.Sprintf("[%v -> %v]", item.A, prod_with_dot)
}

type CA_State []ItemLR0
func (Items CA_State) String() (res string) {
  for i, item := range Items {
    res += fmt.Sprintf("%v. %v\n", i, item)
  }
  return
}

func MakeCanonicAutomatonLR0(grammar *Grammar.Grammar) DFA {
  CA := DFA{}
  fmt.Println(CA)

  return CA
}

func (item ItemLR0) IsIn(Items CA_State) bool {
  fmt.Printf("Is %v in\n%v\n", item, Items)
  for _, comp_item := range Items {
    //fmt.Printf("comparing to\n%v\n", comp_item)
    if item.A == comp_item.A && item.Dot == comp_item.Dot {
      same := true
      if len(item.Prod) != len(comp_item.Prod) {
        same = false
      } else {
        for i := 0; i < len(item.Prod); i++ {
          if item.Prod[i] != comp_item.Prod[i] {
            same = false
          }
        }
      }
      if same {
        return true
      }
    }
  }
  return false
}

func Closure(grammar *Grammar.Grammar, Items CA_State) CA_State {
  already_closed_non_terminal := make([]Grammar.NonTerminal, 0)
  fmt.Printf("Taking the closure of\n[\n%v]\n", Items)
  i := 0
  for i < len(Items) {
    item := Items[i]
    if len(item.Prod) != 0 && item.Dot < len(item.Prod) {
      if X := item.Prod[item.Dot]; grammar.IsNonTerminal(X) && !isInSlice(X, already_closed_non_terminal) {
        fmt.Println("closure", i, ": cosidering item", item)
        for _, X_prod := range grammar.R[X] {
          new_prod := X_prod
          if len(X_prod) == 1 && X_prod[0] == EPS {
            new_prod = make(Grammar.Production, 0)
          }
          new_item := ItemLR0{
            A: X,
            Prod: new_prod,
            Dot: 0,
          }
          if !new_item.IsIn(Items) {
            fmt.Println(new_item, "non c'è, lo aggiungo")
            Items = append(Items, new_item)
          } else {
            fmt.Println(new_item, "c'è già")
          }
        }
        already_closed_non_terminal = append(already_closed_non_terminal, X)
      }
    }
    i++
  } 
  return Items
}

func Goto(grammar *Grammar.Grammar, Items CA_State, X Grammar.NonTerminal) CA_State {
  fmt.Println("Goto with", X)
  J := make(CA_State, 0)
  for _, item := range Items {
    fmt.Printf("goto: Considering item\n%v\n", item)
    if item.Dot < len(item.Prod) && item.Prod[item.Dot] == X {
      new_item := ItemLR0{
        A: item.A,
        Prod: item.Prod,
        Dot: item.Dot+1,
      }
      fmt.Printf("New item:\n%v\n", new_item)
      J = append(J, new_item)
    }
  }
  fmt.Println("Goto without closure:", J)
  return Closure(grammar, J)
}

// TESTS

func TestAddState(t *testing.T) {
  N := new(NFA)
  N.AddState("A", FINAL)
  if N.n != 1 {
    t.Fatalf("States number is not 1\n")
  }
  A := &N.States[0]
  if A == nil {
    t.Fatalf("State is null\n")
  }
  if strings.Compare(A.label, "0") != 0 {
    t.Fatalf("Label is not '0'")
  }
}

func Test_E_clos() {
  N := Test_MakeNFA_star()
  states_e_clos := []int{0}
  e_clos := N.E_clos(states_e_clos)
  fmt.Println("States:", states_e_clos)
  fmt.Println("E_clos:", e_clos)
}

func Test_Move() {
  N := Test_MakeNFA_star()
  states_move := []int{1}
  a := "a"
  move := N.Move(states_move, a)
  fmt.Println("States:", states_move)
  fmt.Printf("Move with '%v': %v\n", a, move)

}

func test_NFA2DFA() {
  N := Test_MakeNFA_star()
  M := N.ToDFA()
  fmt.Println("Il mio bel DFA\n", M)
}

func Test_minimize() {
  N := Test_MakeNFA_to_minimize()
  //N := Test_MakeNFA_to_minimize_mapiùfacile()
  fmt.Println(N)

  M := N.Minimize()
  fmt.Println(M)
}

func Test_permutationString() {
  sigma := "abc"
  max_len := 3
  p_n := 3
  w := GetPermutationString(sigma, max_len, p_n)
  fmt.Println("w:", w)
}

// NFA from regexp a*
func Test_MakeNFA_star() NFA {
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

func Test_MakeNFA_to_minimize() NFA {
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

func Test_MakeNFA_to_minimize_mapiùfacile() NFA {
  N := NFA{}

  var A, B, C, D int

  A = N.GetStatesNum()
  N.AddState("A", NON_FINAL)
  B = N.GetStatesNum()
  N.AddState("B", NON_FINAL)
  C = N.GetStatesNum()
  N.AddState("C", NON_FINAL)

  D = N.GetStatesNum()
  N.AddState("D", FINAL)

  N.AddTransition("a", A, B)
  N.AddTransition("b", A, C)
  N.AddTransition("a", B, C)
  N.AddTransition("b", B, D)
  N.AddTransition("a", C, B)
  N.AddTransition("b", C, D)
  N.AddTransition("a", D, D)
  N.AddTransition("b", D, D)

  N.InitialState = A

  return N
}

func TestCopyNFA(t *testing.T) {
  N := Test_MakeNFA_star()
  fmt.Println("N:", N)
  M := N.Copy()
  fmt.Println("M:", M)
}
