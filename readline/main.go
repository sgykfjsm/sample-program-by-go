package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/chzyer/readline"
)

var (
	bucketName = []byte("history")
)

const (
	defaultPrompt = ""
)

type InputData struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Content   string    `json:"content"`
}

func NewDB(filename string) (*bolt.DB, error) {
	return bolt.Open(filename, 0600, nil)
}

func SaveInput(db *bolt.DB, data []string) (int, error) {
	var i int
	err := db.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}

		input := InputData{
			ID:        int(id),
			CreatedAt: time.Now(),
			Content:   strings.Join(data, "\n"),
		}

		buf, err := json.Marshal(input)
		if err != nil {
			return err
		}

		if err := b.Put(itob(input.ID), buf); err != nil {
			return err
		}
		i++

		return nil
	})

	return i, err
}

func itob(i int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))

	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func ShowHistory(db *bolt.DB) {
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(bucketName).Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%d, value=%s\n", btoi(k), v)
		}

		return nil
	})
}

func NewReadLine(prompt string) (*readline.Instance, error) {
	return readline.New(prompt)
}

func StoreInput(rl *readline.Instance, c chan string, quit chan struct{}) {
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt { // Ctrl+c
			quit <- struct{}{}
			break
		}
		c <- strings.TrimSpace(line)
	}

	return
}

func main() {
	db, err := NewDB("./history.db")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	rl, err := NewReadLine(defaultPrompt)
	if err != nil {
		log.Panicln(err)
	}
	defer rl.Close()

	var inputs []string
	inputChan := make(chan string)
	quitChan := make(chan struct{})
	fmt.Println("Paste it!")
	go StoreInput(rl, inputChan, quitChan)

inputLoop:
	for {
		select {
		case input := <-inputChan:
			inputs = append(inputs, input)
		case <-quitChan:
			break inputLoop
		}
	}
	saved, err := SaveInput(db, inputs)
	fmt.Printf("Saved %d lines.\n", saved)

	ShowHistory(db)
	// readline.ClearScreen(rl) // Ctrl+r
}
