package main

import (
	"log"
	"math/rand"
	"os"
	"strings"
)

const sourceText = "moby.txt"

type MarkovChain struct {
	chain map[string][]string
}

func NewBabbler() *MarkovChain {
	text, err := os.ReadFile(sourceText)
	if err != nil {
		log.Fatal("Could not read source text for babbler", err)
	}
	words := strings.Fields(string(text))
	chain := make(map[string][]string)
	for i := 0; i < len(words)-1; i++ {
		chain[words[i]] = append(chain[words[i]], words[i+1])
	}
	return &MarkovChain{chain: chain}
}

func (m *MarkovChain) Generate(maxWords int) string {
	if len(m.chain) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m.chain))
	for k := range m.chain {
		keys = append(keys, k)
	}
	word := keys[rand.Intn(len(keys))]
	if _, ok := m.chain[word]; !ok {
		keys := make([]string, 0, len(m.chain))
		for k := range m.chain {
			keys = append(keys, k)
		}
		word = keys[rand.Intn(len(keys))]
	}

	var result []string
	result = append(result, word)

	for i := 1; i < maxWords; i++ {
		next, ok := m.chain[word]
		if !ok || len(next) == 0 {
			break
		}
		word = next[rand.Intn(len(next))]
		result = append(result, word)
	}

	return strings.Join(result, " ")
}
