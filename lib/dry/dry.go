package dry

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Dictionary struct {
	Words   []string
	Phrases []string
	Names   []string
}

func NewDict(files ...string) (*Dictionary, error) {
	var d Dictionary
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		words := strings.Split(string(data), " ")
		fmt.Printf("word count: %d\n", len(words))
		d.Words = append(d.Words, words...)
	}
	return &d, nil
}

func RandString(n int) string {
	return strings.Repeat("x", n)
}

func (d *Dictionary) RandWord() string {
	l := len(d.Words)
	rand.New(rand.NewSource(time.Now().UnixNano()))
	i := rand.Intn(l)
	return d.Words[i]
}
