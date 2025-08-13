package chorprog

import (
	"errors"
	"net/url"
	"strconv"
)

type Progression struct{
	Key string
	ModeNum int
	ModeText string
	ChordTypes string
	ChordNums []int
	ChordNames []string
}

// GetMode returns the modes name based on its number
func GetMode(num int) (string, error){
	switch num{
		case 1:
			return "Ionian (major)", nil
		case 2:
			return "Dorian", nil
		case 3:
			return "Phrygian", nil
		case 4:
			return "Lydian", nil
		case 5:
			return "Mixolydian", nil
		case 6:
			return "Aeolian (minor)", nil
		case 7:
			return "Locrian", nil
	default:
		return "", errors.New("not valid mode number")
	}
}

// GetModeNumber takes a string representing a mode and returns the corresponding
// int that represents the mode.
func GetModeNumber(s string) (int, error){
	//TODO:
	return 0, errors.New("not yet implemented")
}

// GetProgFromParams parses the parameter from url.Values and returns the
// requested progression in a Progression struct and nil.
// Or empty progression struct and an error
func GetProgFromParams(params url.Values) (Progression, error){
	var prog Progression
	var err error
	prog.ChordTypes = params.Get("type")
	prog.Key = params.Get("key")
	prog.ModeNum, err = strconv.Atoi(params.Get("mode"))
	if err != nil {
		return Progression{}, err
	}
	prog.ModeText, err = GetMode(prog.ModeNum)
	if err != nil {
		return Progression{}, err
	}

	reqChords := params["chord"]
	var chordNums []int
	for _, c := range reqChords{
		newChord, _ := strconv.Atoi(c)
		if newChord >= 1 && newChord <= 7 {
			chordNums = append(chordNums, newChord)
		}
	}

	prog.ChordNums = chordNums

	//TODO:
	prog.ChordNames = []string{"Ab", "Bb", "C"}
	return prog, nil
}

// getChordNames takes the chod numbers a
func getChordNames(mode int, key string, chordNums []int) ([]string, error){
	//TODO:
	return []string{}, errors.New("Not implemented")
}
