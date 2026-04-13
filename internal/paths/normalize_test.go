package paths

import (
	"testing"
)

// TestNormalizeInputPath は入力パスの正規化をテストする
func TestNormalizeInputPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		// WSL形式
		{`/mnt/c/hoge`, `C:\hoge`},
		{`/mnt/d/path/to/file`, `D:\path\to\file`},
		// Linux形式
		{`/c/hoge`, `C:\hoge`},
		{`/d/path/to/file`, `D:\path\to\file`},
		// file:// URL形式（Windows）
		{`file:///C:/hoge`, `C:\hoge`},
		{`file:///D:/path/to/file`, `D:\path\to\file`},
		// file:// URL形式（UNC）
		{`file://server/share/path`, `\\server\share\path`},
		// UNC拡張長パス
		{`\\?\C:\hoge`, `C:\hoge`},
		{`\\?\D:\path\to\file`, `D:\path\to\file`},
		// Windows形式（スラッシュ、小文字ドライブ→大文字化）
		{`c:/hoge`, `C:\hoge`},
		{`C:/path/to/file`, `C:\path\to\file`},
		// Windows形式（既に正規化済み）
		{`C:\hoge`, `C:\hoge`},
		// 相対パス
		{`./sub/path`, `.\sub\path`},
		{`../parent`, `..\parent`},
		{`.`, `.`},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeInputPath(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeInputPath(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestConvert_MultipleInputFormats は異なる入力形式が正しく変換されることをテストする
func TestConvert_MultipleInputFormats(t *testing.T) {
	// 複数の入力形式が同じパスを表す
	inputFormats := []string{
		`/mnt/c/hoge`,     // WSL形式
		`/c/hoge`,         // Linux形式
		`c:/hoge`,         // Windows相対形式
		`C:\hoge`,         // Windows絶対形式
		`\\?\C:\hoge`,     // UNC拡張長パス
		`file:///C:/hoge`, // file:// URL形式
	}

	tests := []struct {
		format   Format
		expected string
	}{
		{Windows, `C:\hoge`},
		{Home, `C:/hoge`}, // HOMEが /mnt/c/Users/... の場合を想定
		{DOpus, `C:/hoge`},
		{VSCode, `C:/hoge`},
		{GitBash, `/c/hoge`},
		{WSL, `/mnt/c/hoge`},
		{Linux, `/c/hoge`},
		{Escape, `C:\\hoge`},
		{UNC, `\\?\C:\hoge`},
		{URL, `file:///C:/hoge`},
	}

	for _, input := range inputFormats {
		t.Run("input:"+input, func(t *testing.T) {
			for _, tt := range tests {
				t.Run(tt.format.String(), func(t *testing.T) {
					result := Convert(input, tt.format)
					if result != tt.expected {
						t.Errorf("Convert(%q, %s) = %q, expected %q", input, tt.format.String(), result, tt.expected)
					}
				})
			}
		})
	}
}
