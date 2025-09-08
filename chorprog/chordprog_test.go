package chorprog

import (
	"strings"
	"testing"
)

func Test_GetModeName(t *testing.T){
	data := []struct{
		name string
		modeNum int
		expected string
		errorMsg string
	}{
		{"mode 1. Ionian", 1, "ionian", ""},
		{"mode 2. Dorian", 2, "dorian", ""},
		{"mode 3. Phrygian", 3, "phrygian", ""},
		{"mode 4. Lydian", 4, "lydian", ""},
		{"mode 5. Mixolydian", 5, "mixolydian", ""},
		{"mode 6. Aeolian", 6, "aeolian", ""},
		{"mode 7. Locrian", 7, "locrian", ""},
		{"error, lower that possible", 0, "", "not valid mode number"},
		{"error, lower that possible", -1, "", "not valid mode number"},
		{"error, higher that possible", 8, "", "not valid mode number"},
	}

	for _, d := range data{
		t.Run(d.name, func(t *testing.T) {
			result, err := GetModeName(d.modeNum)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.errorMsg{
				t.Errorf("expected error message '%s', got '%s'", d.errorMsg, errMsg)
			}

			if !strings.Contains(strings.ToLower(result), d.expected){
				t.Errorf("expected result '%s', but got '%s'", d.expected, result)
			}
		})
	}
}

func Test_GetModeNumber(t *testing.T){
	data := []struct{
		name string
		inputString string
		expected int
		errorMsg string
	}{
		{"mode 1. Ionian", "ionian", 1, ""},
		{"mode 2. Dorian", "dorian", 2, ""},
		{"mode 3. Phrygian", "phrygian", 3, ""},
		{"mode 4. Lydian", "lydian", 4, ""},
		{"mode 5. Mixolydian", "mixolydian", 5, ""},
		{"mode 6. Aeolian", "aeolian", 6, ""},
		{"mode 7. Locrian", "locrian", 7, ""},
		{"error, random text", "fdjkla", 0, "no valid mode expression found"},
		{"error, random text 2", "jflsJFKLDASFJPOI(", 0, "no valid mode expression found"},
		{"mode 1. Ionian, case insensitive", "IoNian", 1, ""},
		{"mode 2. Dorian, case insensitive", "dOriAn", 2, ""},
		{"mode 3. Phrygian, case insensitive", "PHRYgian", 3, ""},
		{"mode 4. Lydian, case insensitive", "LYDIan", 4, ""},
		{"mode 5. Mixolydian, case insensitive", "mIXolydIAN", 5, ""},
		{"mode 6. Aeolian, case insensitive", "aeOLIan", 6, ""},
		{"mode 7. Locrian, case insensitive", "Locrian", 7, ""},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result, err := GetModeNumber(d.inputString)

			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}

			if errMsg != d.errorMsg{
				t.Errorf("expected error msg: '%s', got '%s'", d.errorMsg, errMsg)
			}

			if result != d.expected {
				t.Errorf("expected result: '%d', got: '%d'", d.expected, result)
			}
		})
	}
}
