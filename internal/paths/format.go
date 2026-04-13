package paths

// Format は出力形式を表す型
type Format int

const (
	// Windows は Windows ネイティブ形式（例: C:\path\to\file）
	Windows Format = iota
	// Home はホームディレクトリを ~ で表す形式（例: ~/Documents/file）
	Home
	// DOpus は Directory Opus 用の形式（例: /profile/Documents/file）
	DOpus
	// VSCode は VSCode の環境変数形式（例: ${env:USERPROFILE}/Documents/file）
	VSCode
	// GitBash は Git Bash の環境変数形式（例: $USERPROFILE/Documents/file）
	GitBash
	// WSL は WSL 形式（例: /mnt/c/Users/user/Documents/file）
	WSL
	// Linux は Linux 風の形式（例: /c/Users/user/Documents/file）
	Linux
	// Escape は Windows パスのエスケープ形式（例: C:\\path\\to\\file）
	Escape
	// UNC は拡張長パス形式（例: \\?\C:\path\to\file）
	UNC
	// URL は file URL 形式（例: file:///C:/path/to/file）
	URL
)

// String はフォーマット名を返す
func (f Format) String() string {
	switch f {
	case Windows:
		return "Windows"
	case Home:
		return "Home"
	case DOpus:
		return "DOpus"
	case VSCode:
		return "VSCode"
	case GitBash:
		return "GitBash"
	case WSL:
		return "WSL"
	case Linux:
		return "Linux"
	case Escape:
		return "Escape"
	case UNC:
		return "UNC"
	case URL:
		return "URL"
	default:
		return "Unknown"
	}
}

// DefaultOrder はデフォルトの出力順序
var DefaultOrder = []Format{
	Windows,
	Home,
	DOpus,
	VSCode,
	GitBash,
	WSL,
	Linux,
	Escape,
	UNC,
	URL,
}
