package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s, This IS zeelang\n", user.Username)
	fmt.Print("feel free to start\n")
	repl.STart(os.Stdin, os.Stdout)

}
