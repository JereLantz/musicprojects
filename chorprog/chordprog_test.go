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
	}
}
