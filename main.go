package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"pathconv/internal/paths"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "pathconv",
		Usage: "カレントディレクトリまたは指定されたパスを様々な形式で表示します",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "windows",
				Aliases: []string{"W"},
				Usage:   "Windows 形式だけを出力 (例: C:\\Users\\...)",
			},
			&cli.BoolFlag{
				Name:    "home",
				Aliases: []string{"H"}, // -h ではなく -H を使用（help と競合回避）
				Usage:   "Home 形式だけを出力 (例: ~/Documents/...)",
			},
			&cli.BoolFlag{
				Name:    "dopus",
				Aliases: []string{"d"},
				Usage:   "DOpus 形式だけを出力 (例: /profile/Documents/...)",
			},
			&cli.BoolFlag{
				Name:    "vscode",
				Aliases: []string{"v"},
				Usage:   "VSCode 形式だけを出力 (例: ${env:USERPROFILE}/Documents/...)",
			},
			&cli.BoolFlag{
				Name:    "gitbash",
				Aliases: []string{"g"},
				Usage:   "Git Bash 形式だけを出力 (例: $USERPROFILE/Documents/...)",
			},
			&cli.BoolFlag{
				Name:    "wsl",
				Aliases: []string{"w"},
				Usage:   "WSL 形式だけを出力 (例: /mnt/c/Users/...)",
			},
			&cli.BoolFlag{
				Name:    "linux",
				Aliases: []string{"l"},
				Usage:   "Linux 形式だけを出力 (例: /c/Users/...)",
			},
			&cli.BoolFlag{
				Name:    "escape",
				Aliases: []string{"e"},
				Usage:   "Windows エスケープ形式だけを出力 (例: C:\\\\Users\\\\...)",
			},
			&cli.BoolFlag{
				Name:    "unc",
				Aliases: []string{"u"},
				Usage:   "UNC 形式だけを出力 (例: \\\\?\\C:\\Users\\...)",
			},
			&cli.BoolFlag{
				Name:    "url",
				Aliases: []string{"U"},
				Usage:   "file URL 形式だけを出力 (例: file:///C:/Users/...)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			// パスの取得（引数なしの場合はカレントディレクトリ）
			targetPath := ""
			if cmd.Args().Len() > 1 {
				return fmt.Errorf("引数が多すぎます。パスは0個または1個のみ指定できます")
			}
			if cmd.Args().Len() == 1 {
				targetPath = cmd.Args().Get(0)
			} else {
				wd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("カレントディレクトリの取得に失敗しました: %w", err)
				}
				targetPath = wd
			}

			// フラグ名とフォーマットのマッピング（エイリアスも含む）
			flagMap := map[string]paths.Format{
				"windows": paths.Windows,
				"W":       paths.Windows,
				"home":    paths.Home,
				"H":       paths.Home, // h から H に変更（help と競合回避）
				"dopus":   paths.DOpus,
				"d":       paths.DOpus,
				"vscode":  paths.VSCode,
				"v":       paths.VSCode,
				"gitbash": paths.GitBash,
				"g":       paths.GitBash,
				"wsl":     paths.WSL,
				"w":       paths.WSL,
				"linux":   paths.Linux,
				"l":       paths.Linux,
				"escape":  paths.Escape,
				"e":       paths.Escape,
				"unc":     paths.UNC,
				"u":       paths.UNC,
				"url":     paths.URL,
				"U":       paths.URL,
			}

			// os.Argsから指定順を取得（ただし、cmd.Bool()で実際にセットされたかを確認）
			formats := []paths.Format{}
			seen := make(map[paths.Format]bool)

			for _, arg := range os.Args[1:] {
				// パス引数はスキップ
				if !strings.HasPrefix(arg, "-") {
					continue
				}
				// フラグを正規化
				flagName := strings.TrimPrefix(strings.TrimPrefix(arg, "--"), "-")

				if format, ok := flagMap[flagName]; ok {
					if !seen[format] {
						// cmd.Bool()で実際にセットされているか確認
						// （flagNameは長い形式に正規化）
						var longName string
						switch format {
						case paths.Windows:
							longName = "windows"
						case paths.Home:
							longName = "home"
						case paths.DOpus:
							longName = "dopus"
						case paths.VSCode:
							longName = "vscode"
						case paths.GitBash:
							longName = "gitbash"
						case paths.WSL:
							longName = "wsl"
						case paths.Linux:
							longName = "linux"
						case paths.Escape:
							longName = "escape"
						case paths.UNC:
							longName = "unc"
						case paths.URL:
							longName = "url"
						}

						if cmd.Bool(longName) {
							formats = append(formats, format)
							seen[format] = true
						}
					}
				}
			}

			// フラグが指定されていない場合はデフォルト順序
			if len(formats) == 0 {
				formats = paths.DefaultOrder
			}

			// 各フォーマットで変換して出力
			for _, format := range formats {
				result := paths.Convert(targetPath, format)
				fmt.Println(result)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
}
