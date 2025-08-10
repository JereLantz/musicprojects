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

// GetProgFromParams parses the parameter from url.Values and returns the
// requested progression in a Progression struct and nil.
// Or empty progression struct and an error
func GetProgFromParams(params url.Values) (Progression, error){
	//TODO:
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
	/*
	prog.Chords = params["chord"]
	*/
	prog.ChordNums = []int{1,2,3}
	return prog, nil
}
