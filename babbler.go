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
	words []string
}

func NewBabbler() *MarkovChain {
	text, err := os.ReadFile(sourceText)
	if err != nil {
		log.Fatal("Could not read source text for babbler", err)
	}
	fields := strings.Fields(string(text))

	chain := make(map[string][]string)
	for i := 0; i < len(fields)-1; i++ {
		chain[fields[i]] = append(chain[fields[i]], fields[i+1])
	}

	seen := make(map[string]struct{})
	words := make([]string, 0, len(chain))
	for _, field := range fields {
		slug := Slugify(field)
		if slug == "" {
			continue
		}
		if _, dup := seen[slug]; dup {
			continue
		}
		seen[slug] = struct{}{}
		words = append(words, slug)
	}

	return &MarkovChain{chain: chain, words: words}
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

func (m *MarkovChain) Links(n int) []string {
	if len(m.words) == 0 {
		return nil
	}
	if n > len(m.words) {
		n = len(m.words)
	}
	seen := make(map[string]struct{}, n)
	links := make([]string, 0, n)
	for len(links) < n {
		word := m.words[rand.Intn(len(m.words))]
		if _, dup := seen[word]; dup {
			continue
		}
		seen[word] = struct{}{}
		links = append(links, word)
	}
	return links
}
