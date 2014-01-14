package main

import (
  "errors"
  "fmt"
)

func pick(thing string, max int) (choice int, err error) {
  if max == 1 {
    fmt.Printf("Only one %s; continue? [yN] ", thing)
    var yn string
    fmt.Scanf("%s", &yn)
    if yn == "y" {
      return 0, nil
    }
    return -1, errors.New("Aborted by user")
  }

  fmt.Printf("Choose a %s [1-%d]: ", thing, max)
  fmt.Scanf("%d", &choice)
  if choice < 1 || choice > max {
    err = errors.New("Invalid selection")
  }
  choice -= 1
  return
}
