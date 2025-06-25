package handlers

import (
	"errors"
	"log"
	"math/rand/v2"
	"musiikkiProjektit/session"
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
	cookie, err := r.Cookie(session.SESSION_TOKEN_NAME)
	if err != nil{
		log.Printf("Failed to fetch the session for displaying the key quiz page. %s\n", err)
		w.WriteHeader(500)
		return
	}
	sessionData := session.Sessions[cookie.Value]

	keyquiz.KeyQuizPage(sessionData).Render(r.Context(), w)
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

	w.Header().Set("Content-Type","text/html")
	components.KeyQuiz(randomKey, quizNotes).Render(r.Context(), w)
}

func HandleCheckQuiz(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(400)
		log.Printf("error parsing the checkquiz form %s\n", err)
		return
	}

	accidentals := r.Form["accidentalSelector"]
	currentKey := r.FormValue("currentKey")

	answerCorrect, err := checkAnswer(currentKey, accidentals)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("Error parsing the answer: %s\n", err)
		return
	}
	components.KeyQuizCheckResp(answerCorrect).Render(r.Context(), w)
}

func checkAnswer(key string, accidentals []string) (bool, error){
	pureNote := strings.Trim(key, "#b")
	noteOffset := 0
	for _, note := range possibleNotes{
		if note == pureNote{
			break
		}
		noteOffset ++
	}

    switch key{
    case "C","G","D","A","E","B","F#","C#":
		accidentalsInThis := accidentalOrder.Sharps[:correctAccidentalAmounts[key]]

		for i, accidentalToCheck := range accidentals{
			shouldBeAccidental := false
			currentNote := possibleNotes[(noteOffset+i) % len(possibleNotes)]

			// loop every accidental and note to check if it should be accidental
			for _, acc := range accidentalsInThis{
				if currentNote == acc {
					shouldBeAccidental = true
				}

				// break if the note is correctly identified as flat
				if currentNote == acc && accidentalToCheck == "sharp"{
					break
				}
			}

			if accidentalToCheck == "natural" && shouldBeAccidental {
				return false, nil
			}

			if accidentalToCheck == "sharp" && !shouldBeAccidental {
				return false, nil
			}
		}
		return true, nil

    case "F","Bb","Eb","Ab","Db","Gb","Cb":
		accidentalsInThis := accidentalOrder.Flats[:correctAccidentalAmounts[key]]

		for i, accidentalToCheck := range accidentals{
			shouldBeAccidental := false
			currentNote := possibleNotes[(noteOffset+i) % len(possibleNotes)]

			// loop every accidental and note to check if it should be accidental
			for _, acc := range accidentalsInThis{
				if currentNote == acc {
					shouldBeAccidental = true
				}

				// break if the note is correctly identified as flat
				if currentNote == acc && accidentalToCheck == "flat"{
					break
				}
			}

			if accidentalToCheck == "natural" && shouldBeAccidental {
				return false, nil
			}

			if accidentalToCheck == "flat" && !shouldBeAccidental {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, errors.New("Invalid key")
	}
}
