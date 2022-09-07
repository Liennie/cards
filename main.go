package main

import (
	"log"

	"golang.org/x/exp/slices"
)

const (
	pickedCards = 5
	deckSize    = 124
)

func combinationsRecursive(i, depth int, comb []int, ch chan<- []int) {
	for j := i; j < deckSize; j++ {
		if depth == pickedCards {
			ch <- append(comb[:len(comb):len(comb)], j)
		} else {
			combinationsRecursive(j+1, depth+1, append(comb, j), ch)
		}
	}
}

func combinations() <-chan []int {
	ch := make(chan []int)

	go func() {
		defer close(ch)

		combinationsRecursive(0, 1, nil, ch)
	}()

	return ch
}

func sum(is []int) int {
	sum := 0
	for _, i := range is {
		sum += i
	}
	return sum
}

func choose(comb []int) (int, []int) {
	i := sum(comb) % len(comb)
	return comb[i], append(comb[:i:i], comb[i+1:]...)
}

func possibilities(cards []int) []int {
	res := []int{}

	last := 0
	sch := sum(cards)

	for i, c := range append(cards, deckSize) {
		for j := last; j < c; j++ {
			if (sch+j)%pickedCards == i {
				res = append(res, j)
			}
		}

		last = c + 1
	}

	return res
}

var permlut = [][]int{
	{0, 1, 2, 3},
	{1, 0, 2, 3},
	{2, 0, 1, 3},
	{0, 2, 1, 3},
	{1, 2, 0, 3},
	{2, 1, 0, 3},
	{2, 1, 3, 0},
	{1, 2, 3, 0},
	{3, 2, 1, 0},
	{2, 3, 1, 0},
	{1, 3, 2, 0},
	{3, 1, 2, 0},
	{3, 0, 2, 1},
	{0, 3, 2, 1},
	{2, 3, 0, 1},
	{3, 2, 0, 1},
	{0, 2, 3, 1},
	{2, 0, 3, 1},
	{1, 0, 3, 2},
	{0, 1, 3, 2},
	{3, 1, 0, 2},
	{1, 3, 0, 2},
	{0, 3, 1, 2},
	{3, 0, 1, 2},
}

func permutate(cards []int, i int) []int {
	perm := permlut[i]
	return []int{cards[perm[0]], cards[perm[1]], cards[perm[2]], cards[perm[3]]}
}

func encode(card int, rest []int) []int {
	p := possibilities(rest)
	i, ok := slices.BinarySearch(p, card)
	if !ok {
		log.Panicf("Chosen card %d (%v) not found in possibilities (%v)", card, rest, p)
	}
	return permutate(rest, i)
}

func decode(cards []int) int {
	sorted := slices.Clone(cards)
	slices.Sort(sorted)

perm:
	for i := range permlut {
		for j, ec := range permutate(sorted, i) {
			if ec != cards[j] {
				continue perm
			}
		}
		return i
	}
	log.Panic("Not found")
	return 0
}

func guess(cards []int) int {
	i := decode(cards)
	slices.Sort(cards)
	return possibilities(cards)[i]
}

func main() {
	tested := 0
	testedlog := 0

	for comb := range combinations() {
		card, rest := choose(comb)

		g := guess(encode(card, rest))
		if g != card {
			log.Panicf("Card guessed wrong, %d != %d, for combination %v", g, card, comb)
		}

		tested++
		testedlog++
		if testedlog == 100000 {
			testedlog = 0
			log.Printf("Tested %d combinations", tested)
		}
	}
	log.Printf("Tested all %d combinations", tested)
	log.Printf("No mistakes found, congratulations!")
}
