package main

import (
  "github.com/keita0805carp/cacis/cmd"
  "github.com/keita0805carp/cacis/master"
  "github.com/keita0805carp/cacis/slave"
)

func main() {
    cmd.Run()
    master.Main()
    slave.Main()
}
