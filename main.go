package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type room int

const (
	M room = iota + 1
	K
	Q
	B
	L
)

func (r room) String() string {
	switch r {
	case M:
		return "Master Suite"
	case K:
		return "King Suite"
	case Q:
		return "Queen Bunk Room"
	case B:
		return "Bunk Room"
	case L:
		return "Living Room"
	}
	return fmt.Sprintf("Room(%d)", r)
}

var vrbo = [...]room{M, K, K, Q, Q, Q, Q, Q, Q, B, B, B, B, B, L}

type guy struct {
	name          string
	roomPrefScore map[room]int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	f, err := os.Open("prefs.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	// skip the header
	if _, err = csvReader.Read(); err != nil {
		log.Fatal(err)
	}

	// load prefs
	guys := make([]guy, 0, 14)
	for {
		if record, err := csvReader.Read(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			g := guy{record[0], make(map[room]int)}
			g.roomPrefScore[M], _ = strconv.Atoi(record[M])
			g.roomPrefScore[K], _ = strconv.Atoi(record[K])
			g.roomPrefScore[Q], _ = strconv.Atoi(record[Q])
			g.roomPrefScore[B], _ = strconv.Atoi(record[B])
			g.roomPrefScore[L], _ = strconv.Atoi(record[L])
			guys = append(guys, g)
		}
	}

	// brute force search for the best assignment
	maxScore := 0
	bestAssignments := make([]int, len(guys))
	for n := 0; n < 1000000; n++ {
		score := 0
		assignments := rand.Perm(len(guys))
		for i := range guys {
			room := vrbo[assignments[i]]
			score += guys[i].roomPrefScore[room]
		}
		if score > maxScore {
			maxScore = score
			copy(bestAssignments, assignments)
			fmt.Println("New best score:", maxScore, "n:", n)
		}
	}

	fmt.Println("Best score:", maxScore, "assignments:")
	for i := range guys {
		fmt.Println(guys[i].name, ":", vrbo[bestAssignments[i]])
	}
}
