package chorprog

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type Progression struct{
	Key string
	ModeNum int
	ModeText string
	ChordTypes string
	ChordNums []int
	ChordNames []string
}

// GetModeName returns the modes name based on its number
func GetModeName(num int) (string, error){
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
	loweredS := strings.ToLower(s)

	switch{
	case strings.Contains(loweredS, "ionian"):
		return 1, nil
	case strings.Contains(loweredS, "dorian"):
		return 2, nil
	case strings.Contains(loweredS, "phrygian"):
		return 3, nil
	//NOTE: mixolydian pitää testaa ennen lydian,koska muuten mixolydian myös palauttaisi aina 4
	case strings.Contains(loweredS, "mixolydian"):
		return 5, nil
	case strings.Contains(loweredS, "lydian"):
		return 4, nil
	case strings.Contains(loweredS, "aeolian"):
		return 6, nil
	case strings.Contains(loweredS, "locrian"):
		return 7, nil
	default:
		return 0, errors.New("no valid mode expression found")
	}
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
	prog.ModeText, err = GetModeName(prog.ModeNum)
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

// getChordNames takes the mode number, the key, the chord type and the requested chord numbers,
// and returns the names of the chords as a string array
func getChordNames(mode int, key string, chordType string, chordNums []int) ([]string, error){
	//TODO:
	return []string{}, errors.New("Not implemented")
}
