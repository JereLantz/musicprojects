package handlers

import (
	"musiikkiProjektit/views/chordProgress"
	"net/http"
)


func HandleServeChordProg(w http.ResponseWriter, r *http.Request){
	chordprogress.ChordProgPage().Render(r.Context(), w)
}
