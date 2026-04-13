package paths

import (
	"os"
	"testing"
)

func TestConvert_AllFormats(t *testing.T) {
	// テスト用のHOME設定
	os.Setenv("USERPROFILE", `C:\Users\testuser`)
	defer os.Unsetenv("USERPROFILE")

	testPath := `C:\Users\testuser\Documents\test`

	for _, format := range DefaultOrder {
		t.Run(format.String(), func(t *testing.T) {
			result := Convert(testPath, format)
			t.Logf("Format: %s, Result: %s", format.String(), result)
			if result == "" {
				t.Errorf("Convert returned empty string for format %s", format.String())
			}
		})
	}
}
