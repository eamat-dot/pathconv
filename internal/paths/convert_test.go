package paths

import (
	"testing"
)

func TestQuoteIfNeeded(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `C:\Users\test`},
		{`C:\Program Files\test`, `"C:\Program Files\test"`},
		{`/mnt/c/test`, `/mnt/c/test`},
		{`/mnt/c/my files`, `"/mnt/c/my files"`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := QuoteIfNeeded(tt.input)
			if result != tt.expected {
				t.Errorf("QuoteIfNeeded(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvert_Windows(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `C:\Users\test`},
		{`./sub/path`, `.\sub\path`},
		{`../parent`, `..\parent`},
		{`c:/test`, `C:\test`}, // 小文字ドライブは大文字化される
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Convert(tt.input, Windows)
			// クオートなし版で比較
			if result != tt.expected {
				t.Errorf("Convert(%q, Windows) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvert_Escape(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `C:\\Users\\test`}, // \ を \\ にエスケープ
		{`.\sub\path`, `.\\sub\\path`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Convert(tt.input, Escape)
			if result != tt.expected {
				t.Errorf("Convert(%q, Escape) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvert_Linux(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `/c/Users/test`},
		{`D:\Projects`, `/d/Projects`},
		{`./sub/path`, `./sub/path`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Convert(tt.input, Linux)
			if result != tt.expected {
				t.Errorf("Convert(%q, Linux) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvert_UNC(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `\\?\C:\Users\test`},
		{`\\server\share\path`, `\\server\share\path`},
		{`./sub/path`, `.\sub\path`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Convert(tt.input, UNC)
			if result != tt.expected {
				t.Errorf("Convert(%q, UNC) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// (duplicate TestConvert_UNC removed)

func TestConvert_URL(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`C:\Users\test`, `file:///C:/Users/test`},
		{`\\server\share\path`, `file://server/share/path`},
		{`./sub/path`, `./sub/path`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Convert(tt.input, URL)
			if result != tt.expected {
				t.Errorf("Convert(%q, URL) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
