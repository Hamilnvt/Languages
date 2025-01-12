# Workflow

[x] NFA
[ ] DFA
- [x] fattorizzare, per ora è tutto nel main
- [x] minimizzare
[ ] Grammatica
    > stavo pensando a produzione del tipo A -> K1 | ... | K2 dove gli spazi sono necessari e dividono perfettamente le varie produzioni
- [x] Simboli annulabili
- [x] First
- [x] Follow
- [ ] Semplificazione ? //Potrebbe essere troppo complicato e troppo poco utile
[ ] lexical analyzer
[ ] PDA
[ ] DPDA
[x] parser LL(1)
[ ] parse tree
[ ] parser LR(0)
[ ] parser SLR(1) ?
[ ] parser LR(1)
[ ] parser LALR(1)
[ ] cominciare a pensare al mio linguaggio

# TODO

## NFA
[x] Labels degli stati generati automaticamente
[x] Possibilità di ottenere gli stati con le labels
[x] serve un modo per ottenere lo stato di arrivo dallo stato di partenza e la label della transizione (insomma, serve la delta, così da poter fare delta(q, a) e ottenere l'insieme di stati su cui si arriva):
    > è un campo dello State, le cui transizioni voglio che siano una mappa invece che un array
    - [ ] delta func ? (nel DFA la delta potrebbe restituire un solo stato (però sembra difficile questa cosa perché non collabora con tutto il rest))
[ ] RemoveState
[x] RemoveTransition
[ ] NFA ha una mappa da int/string a stato:
    > sorge un dilemma: se voglio le mappe si può fare nei seguenti modi:
    1. 2 mappe: una con le label e una con gli indici (quindi una mappa e un array, suppongo)
    2. 1 mappa[string]: ha il doppio degli elementi e map[label] = elt = map[index], con qualche accortezza questo non è male (func getEltAt(i int) = map[label] = elt, dove getEltAt passa l'intero come stringa)
[ ] Sigma sarà uno slice di string, dove ogni stringa avrà lunghezza 1
[ ] func isDFA che controlla se non ci sono transizioni epsilon e se per ogni stato c'è una transizione per ogni terminale
[ ] func cat: concatena due NFA se il primo ha un solo stato finale (fonde f1 e i2)

## Grammar
[x] String
[ ] First and Follow table in Grammar (as maps)
- [ ] si potrebbe fare che se calcola il first di un nonterminale lo inserisce nella mappa e se lo deve ricacolare, prima di farlo controlla la tabella
[o] Parsing grammar from file
- [x] fare in modo che accetti solo file con estensione .g (for the meme)
nah - [ ] commenti anche a metà riga (mettere # nei simboli slashati e migliorare il parser)
- [x] \eps = ε
- [x] | deve essere escaped
- [ ] definizioni all'inizio del file (tipo alias) ? 
- [ ] se C = a | ... | z , posso scrivere questa notazione e lo capisce
- [x] togliere S, T e NT e dare delle regole precise per parsare (S è il NonTerminale della prima regola, gli altri NonTerminali si prendono dalle regole successive e tutti gli altri simboli sono terminali)
- [ ] spazio per le definizioni regolari (i pattern dei token)

## Parsing
[ ] abstract parsing tree from the concrete one
[ ] i parser devono leggere i token, non i singoli caratteri
[ ] risolvere problema: quando lo stato esiste già non lo segna correttamente nella delta

# Cose da leggere

- paper sul lexical analyzer
- https://www.reddit.com/r/ProgrammingLanguages/comments/15cxb1a/advice_on_building_an_lalrk_parser_generator/?rdt=62283
- https://www.sciencedirect.com/topics/computer-science/parse-tree#:~:text=A%20parse%20tree%20is%20a,symbol%20used%20in%20the%20derivation.
- https://www.site.uottawa.ca/~bochmann/SEG-2106-2506/Notes/M2-3-SyntaxAnalysis/grammar-for-regular-expressions.html

Video abstract syntax tree construction
- https://www.youtube.com/watch?v=q4tdwlAU1-M

# Cos'ho capito

- Nel lexical analizer (LA), ogni stato finale del DFA corrisponde ad un token (la cui descrizione viene data in fase di dichiarazione con un'espressione regolare), e quando il DFA, leggendo una stringa, termina su quello stato finale, restituisce il token e la stringa letta
- Come procedere:
  - Si danno le definizioni regolari per i terminali
  - si crea il DFA che fungerà da LA che ritorna i token
  - Si definiscono le regole con terminali e non terminali
  - nel Parsing quando si ritorna un terminal gli si associa il token riconosciuto (e il valore che devo capire a che punto assegnarlo)
