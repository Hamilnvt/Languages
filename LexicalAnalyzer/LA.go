package LexicalAnalyzer

import (
  "Languages/NFA"
  "Languages/Parsing"
  "bufio"
  "strings"
  "os"
  "fmt"
)

type LexicalAnalyzer struct {
  dfa NFA.DFA
  symbolTable map[string]string
  parser Parsing.Parser_LL1
}

func MakeLexicalAnalyzer(la_path string) LexicalAnalyzer {
  if la_path[len(la_path)-3:] != ".la" {
    panic("File extension should be .la")
  }

  file, err := os.Open(la_path)
  if err != nil {
    panic(err)
  }

  LA := LexicalAnalyzer{
    symbolTable: make(map[string]string),
  }
  
  fmt.Println("Scanning file:\n")
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := strings.TrimSpace(scanner.Text())

    if len(line) > 0 && line[0] != '#' {
      fmt.Println("Parsing", line)
      is_new := false
      splitted := strings.Fields(line)
      for i, word := range splitted {
        fmt.Println(i, word)
        if i == 0 {
          if word[len(word)-1] != ':' {
            panic("Invalid declaration, should be of the form:\nIDENTIFIER: value")
          }
          new_identifier := word[:len(word)-1]
          if _, ok := LA.symbolTable[new_identifier]; ok {
            panic(fmt.Sprintf("Redeclaring %v", new_identifier))
          } else {
            fmt.Println("New identifier:", new_identifier)
            is_new = true
          }
        } else if len(word) > 1 && word[0] == '$' {
          fmt.Println("Found identifier to substitute:", word[1:]) 
          if value, ok := LA.symbolTable[word[1:]]; !ok {
            panic(fmt.Sprintf("Undeclared %v", word[1:]))
          } else {
            splitted[i] = "("+value+")"
          }
        }
      }
      // pensavo di tenerli come identificatori perché credo possa semplificare la costruzione del DFA
      if is_new {
        value := strings.Join(splitted[1:], " ")
        identifier := splitted[0][:len(splitted[0])-1]
        LA.symbolTable[identifier] = value
        fmt.Println("Added identifier", identifier, "with value", value)
      }
    }
  }

  if err := scanner.Err(); err != nil {
    panic(err)
  }
  file.Close()
  fmt.Println("\nScanning file ended without errors.\n")
  for identifier, value := range LA.symbolTable {
    fmt.Printf("%v: %v\n", identifier, value)
  }

  return LA
}
