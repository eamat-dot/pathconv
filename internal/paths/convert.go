package paths

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// getHome はホームディレクトリのパスを取得する
func getHome() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if userprofile := os.Getenv("USERPROFILE"); userprofile != "" {
		return userprofile
	}
	return ""
}

// expandTilde は ~ で始まるパスをホームディレクトリに展開する
func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}
	home := getHome()
	if home == "" {
		return path
	}
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return filepath.Join(home, path[2:])
	}
	return path
}

// isRelative はパスが相対パスかどうかを判定する
func isRelative(path string) bool {
	// 絶対パスのパターンをチェック
	if filepath.IsAbs(path) {
		return false
	}
	// Windows ドライブレター形式（c:\ や c:/ など）
	if len(path) >= 2 && path[1] == ':' {
		return false
	}
	// UNC パスのチェック
	if strings.HasPrefix(path, "\\\\") || strings.HasPrefix(path, "//") {
		return false
	}
	// WSL/Linux 形式の絶対パスのチェック
	if strings.HasPrefix(path, "/mnt/") {
		return false
	}
	if strings.HasPrefix(path, "/") && len(path) >= 2 && len(string(path[1])) == 1 {
		// /c/... 形式のチェック
		return false
	}
	return true
}

// normalizeInputPath は様々な形式の入力パスをWindows形式に正規化する
// （相対パスはそのまま返す）
func normalizeInputPath(path string) string {
	// file:// URL 形式を最初にチェック（相対判定より前）
	if strings.HasPrefix(path, "file:///") {
		// file:///C:/path → C:\path
		urlPath := path[8:] // "file:///" を削除
		normalized := strings.ReplaceAll(urlPath, "/", "\\")
		return normalizeWindowsPath(normalized)
	}
	if strings.HasPrefix(path, "file://") {
		// file://server/share/path → \\server\share\path
		netloc := path[7:] // "file://" を削除
		return "\\\\" + strings.ReplaceAll(netloc, "/", "\\")
	}

	// UNC 拡張長パスも早めにチェック
	if strings.HasPrefix(path, "\\\\?\\") {
		normalized := path[4:] // "\\?\" を削除
		return normalizeWindowsPath(normalized)
	}

	// 相対パスの場合はスラッシュを正規化するだけ
	if isRelative(path) {
		return strings.ReplaceAll(path, "/", "\\")
	}

	// WSL パス: /mnt/c/path → C:\path
	if strings.HasPrefix(path, "/mnt/") {
		parts := strings.Split(path, "/")
		if len(parts) >= 3 {
			drive := strings.ToUpper(parts[2])
			rest := strings.Join(parts[3:], "\\")
			if rest == "" {
				return drive + ":\\"
			}
			return drive + ":\\" + rest
		}
	}

	// Linux 形式: /c/path → C:\path
	if strings.HasPrefix(path, "/") {
		parts := strings.Split(path, "/")
		if len(parts) >= 2 && len(parts[1]) == 1 {
			drive := strings.ToUpper(parts[1])
			rest := strings.Join(parts[2:], "\\")
			if rest == "" {
				return drive + ":\\"
			}
			return drive + ":\\" + rest
		}
	}

	// Windows 形式（既に正規化済み）: C:/path または C:\path または c:/path
	return normalizeWindowsPath(path)
}

// normalizeWindowsPath は Windows パスのドライブレターを大文字化し、スラッシュをバックスラッシュに統一
func normalizeWindowsPath(path string) string {
	// スラッシュをバックスラッシュに統一
	path = strings.ReplaceAll(path, "/", "\\")

	// ドライブレターが存在する場合は大文字化（c:\ → C:\）
	if len(path) >= 2 && path[1] == ':' {
		path = strings.ToUpper(string(path[0])) + path[1:]
	}

	return path
}

// QuoteIfNeeded はスペースを含むパスを引用符で囲む
func QuoteIfNeeded(path string) string {
	if strings.Contains(path, " ") {
		return `"` + path + `"`
	}
	return path
}

// Convert は指定されたパスを指定フォーマットに変換する
func Convert(inputPath string, format Format) string {
	// 入力パスを正規化（様々な形式を Windows 形式に）
	path := normalizeInputPath(inputPath)

	// ~ 展開（フルパス系フォーマットの場合のみ）
	if strings.HasPrefix(path, "~") {
		switch format {
		case Windows, WSL, Linux, Escape, UNC, URL:
			path = expandTilde(path)
		}
	}

	// 正規化後に相対パスを再判定
	relative := isRelative(path)

	home := getHome()
	var result string

	switch format {
	case Windows:
		result = convertToWindows(path, relative)
	case Home:
		result = convertToHome(path, home, relative)
	case DOpus:
		result = convertToDOpus(path, home, relative)
	case VSCode:
		result = convertToVSCode(path, home, relative)
	case GitBash:
		result = convertToGitBash(path, home, relative)
	case WSL:
		result = convertToWSL(path, relative)
	case Linux:
		result = convertToLinux(path, relative)
	case Escape:
		result = convertToEscape(path, relative)
	case UNC:
		result = convertToUNC(path, relative)
	case URL:
		result = convertToURL(path, relative)
	default:
		result = path
	}

	return QuoteIfNeeded(result)
}

// convertToWindows は Windows 形式に変換
func convertToWindows(path string, relative bool) string {
	// 入力が `/` 区切りでも `\` に正規化する（相対・絶対問わず）
	return strings.ReplaceAll(path, "/", "\\")
}

// convertToHome はホーム形式に変換
func convertToHome(path string, home string, relative bool) string {
	result := path
	// HOME 置換を試みる
	if home != "" && strings.HasPrefix(path, home) {
		result = "~" + path[len(home):]
	}
	// スラッシュに変換
	result = strings.ReplaceAll(result, "\\", "/")
	return result
}

// convertToDOpus は DOpus 形式に変換
func convertToDOpus(path string, home string, relative bool) string {
	result := path
	// HOME を /profile に置換
	if home != "" && strings.HasPrefix(path, home) {
		result = "/profile" + path[len(home):]
	}
	// スラッシュに変換
	result = strings.ReplaceAll(result, "\\", "/")
	return result
}

// convertToVSCode は VSCode 形式に変換
func convertToVSCode(path string, home string, relative bool) string {
	result := path
	// HOME を ${env:USERPROFILE} に置換
	if home != "" && strings.HasPrefix(path, home) {
		result = "${env:USERPROFILE}" + path[len(home):]
	}
	// スラッシュに変換
	result = strings.ReplaceAll(result, "\\", "/")
	return result
}

// convertToGitBash は Git Bash 形式に変換
func convertToGitBash(path string, home string, relative bool) string {
	result := path
	// HOME を $USERPROFILE に置換
	if home != "" && strings.HasPrefix(path, home) {
		result = "$USERPROFILE" + path[len(home):]
	}
	// スラッシュに変換
	result = strings.ReplaceAll(result, "\\", "/")
	// ドライブレターを /<lower> に変換（相対パスでない場合）
	if !relative {
		re := regexp.MustCompile(`^([A-Za-z]):`)
		result = re.ReplaceAllStringFunc(result, func(match string) string {
			drive := strings.ToLower(string(match[0]))
			return "/" + drive
		})
	}
	return result
}

// convertToWSL は WSL 形式に変換
func convertToWSL(path string, relative bool) string {
	if relative {
		// 相対パスはスラッシュに変換
		return strings.ReplaceAll(path, "\\", "/")
	}
	// スラッシュに変換
	result := strings.ReplaceAll(path, "\\", "/")
	// ドライブレターを /mnt/<lower> に変換
	re := regexp.MustCompile(`^([A-Za-z]):`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		drive := strings.ToLower(string(match[0]))
		return "/mnt/" + drive
	})
	return result
}

// convertToLinux は Linux 形式に変換
func convertToLinux(path string, relative bool) string {
	if relative {
		// 相対パスはスラッシュに変換
		return strings.ReplaceAll(path, "\\", "/")
	}
	// スラッシュに変換
	result := strings.ReplaceAll(path, "\\", "/")
	// ドライブレターを /<lower> に変換
	re := regexp.MustCompile(`^([A-Za-z]):`)
	result = re.ReplaceAllStringFunc(result, func(match string) string {
		drive := strings.ToLower(string(match[0]))
		return "/" + drive
	})
	return result
}

// convertToEscape は Windows エスケープ形式に変換
func convertToEscape(path string, relative bool) string {
	// \ を \\ にエスケープ
	return strings.ReplaceAll(path, "\\", "\\\\")
}

// convertToUNC は UNC 形式に変換
func convertToUNC(path string, relative bool) string {
	if relative {
		// 相対パスは UNC 非対応のため Windows 形式に準拠
		return strings.ReplaceAll(path, "/", "\\")
	}
	// UNC ネットワークパスはそのまま
	if strings.HasPrefix(path, "\\\\") {
		return path
	}
	// 絶対 Windows パスは拡張長プレフィックスを付与
	if matched, _ := regexp.MatchString(`^[A-Za-z]:`, path); matched {
		return `\\?\` + path
	}
	return path
}

// convertToURL は file URL 形式に変換
func convertToURL(path string, relative bool) string {
	if relative {
		// 相対パスはスラッシュに変換（URLエンコードなし）
		result := strings.ReplaceAll(path, "\\", "/")
		return result
	}

	// UNC ネットワークパス (\\server\share\path) を file://server/share/path に変換
	if strings.HasPrefix(path, "\\\\") || strings.HasPrefix(path, "//") {
		cleaned := strings.TrimPrefix(strings.TrimPrefix(path, "\\\\"), "//")
		cleaned = strings.ReplaceAll(cleaned, "\\", "/")
		return "file://" + cleaned
	}

	// Windows 絶対パス (C:\path) を file:///C:/path に変換
	if matched, _ := regexp.MatchString(`^[A-Za-z]:`, path); matched {
		result := strings.ReplaceAll(path, "\\", "/")
		return "file:///" + result
	}

	// それ以外（Unix 絶対パス等）
	result := strings.ReplaceAll(path, "\\", "/")
	if strings.HasPrefix(result, "/") {
		return "file://" + result
	}
	return "file:///" + result
}
