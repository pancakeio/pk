package main

import (
  "fmt"
)

func white(s string) string {
  return fmt.Sprintf("\x1b[1;37m%s\x1b[0m", s)
}
