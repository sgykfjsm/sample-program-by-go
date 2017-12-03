package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
)

const defaultPrompt = "> "

func main() {
	rl, err := readline.New(defaultPrompt)
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	rl.SetPrompt("")
	var inputs []string
	fmt.Println("Paste it!")
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Ctrl+c
			break
		}
		line = strings.TrimSpace(line)
		inputs = append(inputs, line)
	}

	readline.ClearScreen(rl) // Ctrl+r
	fmt.Printf("You input %d lines.\n", len(inputs))
}
