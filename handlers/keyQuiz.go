package handlers

import (
	"math/rand/v2"
	"musiikkiProjektit/views/components"
	"musiikkiProjektit/views/keyQuiz"
	"net/http"
	"strings"
)

type accidentalOrderStruct struct {
	Sharps []string
	Flats []string
}
var accidentalOrder = accidentalOrderStruct{
    Sharps: []string{"F", "C", "G", "D", "A", "E", "B"},
    Flats: []string{"B", "E", "A", "D", "G", "C", "F"},
}

var possibleNotes = []string{"A","B","C","D","E","F","G"}
var possibleKeys = []string{"F", "C", "G", "D", "A", "E", "B", "F#", "C#", "Bb", "Eb", "Ab", "Db", "Gb", "Cb"}
var correctAccidentalAmounts = map[string]int{"C":0,"G":1,"F":1,"D":2,"Bb":2,"A":3,"Eb":3,"E":4,"Ab":4,"B":5,"Db":5,"F#":6,"Gb":6,"C#":7,"Cb":7}


func HandleServeKeyQuiz(w http.ResponseWriter, r *http.Request){
	keyquiz.KeyQuizPage().Render(r.Context(), w)
}

func HandleStartKeyQuiz(w http.ResponseWriter, r *http.Request){
	var quizNotes []string
	randomKeyIndex := rand.IntN(len(possibleKeys))
	randomKey := possibleKeys[randomKeyIndex]

	pureNote := strings.Trim(randomKey, "#b")
	noteOffset := 0
	for _, note := range possibleNotes{
		if note == pureNote{
			break
		}
		noteOffset ++
	}

	for i := range possibleNotes{
		quizNotes = append(quizNotes, possibleNotes[(i + noteOffset) % len(possibleNotes)])
	}

	components.KeyQuiz(randomKey, quizNotes).Render(r.Context(), w)
}
