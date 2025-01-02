# Workflow

[x] NFA
[] DFA
- [x] fattorizzare, per ora è tutto nel main
- [x] minimizzare
[] Grammatica
    > stavo pensando a produzione del tipo A -> K1 | ... | K2 dove gli spazi sono necessari e dividono perfettamente le varie produzioni
- [] Simboli annulabili
- [] First
- [] Follow
- [] Semplificazione ? //Potrebbe essere troppo complicato e troppo poco utile
[] analizzatore lessicale
[] PDA
[] DPDA
[] parser LL(1)
[] parser LR(0)
[] parser SLR(1) ?
[] parser LR(1)
[] parser LALR(1)
[] cominciare a pensare al mio linguaggio

# TODO

## NFA
[x] Labels degli stati generati automaticamente
[x] Possibilità di ottenere gli stati con le labels
[x] serve un modo per ottenere lo stato di arrivo dallo stato di partenza e la label della transizione (insomma, serve la delta, così da poter fare delta(q, a) e ottenere l'insieme di stati su cui si arriva):
    > è un campo dello State, le cui transizioni voglio che siano una mappa invece che un array
    - [] delta func ? (nel DFA la delta potrebbe restituire un solo stato (però sembra difficile questa cosa perché non collabora con tutto il rest))
[] RemoveState
[x] RemoveTransition
[] NFA ha una mappa da int/string a stato:
    > sorge un dilemma: se voglio le mappe si può fare nei seguenti modi:
    1. 2 mappe: una con le label e una con gli indici (quindi una mappa e un array, suppongo)
    2. 1 mappa[string]: ha il doppio degli elementi e map[label] = elt = map[index], con qualche accortezza questo non è male (func getEltAt(i int) = map[label] = elt, dove getEltAt passa l'intero come stringa)
[] Sigma sarà uno slice di string, dove ogni stringa avrà lunghezza 1
[] func isDFA che controlla se non ci sono transizioni epsilon e se per ogni stato c'è una transizione per ogni terminale

## Grammar
[] String

# Cose da leggere

- https://www.reddit.com/r/ProgrammingLanguages/comments/15cxb1a/advice_on_building_an_lalrk_parser_generator/?rdt=62283
- paper sul lexical analyzer

# Cos'ho capito

- Nel lexical analizer (LA), ogni stato finale del DFA corrisponde ad un token (la cui descrizione viene data in fase di dichiarazione con un'espressione regolare), e quando il DFA, leggendo una stringa, termina su quello stato finale, restituisce il token e la stringa letta
