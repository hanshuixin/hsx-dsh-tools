package dshutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CreateDesktopShortcut 在桌面创建"启动 DeepSeek Harness"快捷方式，指向
// startExe，图标直接使用 startExe 内嵌的应用图标。PowerShell 脚本写入
// 临时的 UTF-8 BOM .ps1 文件，避免中文字符在命令行上被转码破坏。
func CreateDesktopShortcut(startExe string) error {
	script := fmt.Sprintf(`$WshShell = New-Object -comObject WScript.Shell
$Desktop = [System.Environment]::GetFolderPath('Desktop')
$Shortcut = $WshShell.CreateShortcut($Desktop + '\%s.lnk')
$Shortcut.TargetPath = '%s'
$Shortcut.WorkingDirectory = '%s'
$Shortcut.Description = '%s'
$Shortcut.IconLocation = '%s',0
$Shortcut.Save()`,
		psEscape(ShortcutName), psEscape(startExe), psEscape(filepath.Dir(startExe)), psEscape(ShortcutDescription), psEscape(startExe))

	tmp, err := os.CreateTemp("", "dsh-shortcut-*.ps1")
	if err != nil {
		return fmt.Errorf("创建临时脚本失败：%w", err)
	}
	defer os.Remove(tmp.Name())

	// PowerShell 5.1 只有在带 BOM 时才把 UTF-8 文件当作 UTF-8 解释。
	if _, err := tmp.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		tmp.Close()
		return fmt.Errorf("写入临时脚本失败：%w", err)
	}
	if _, err := tmp.WriteString(script); err != nil {
		tmp.Close()
		return fmt.Errorf("写入临时脚本失败：%w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("写入临时脚本失败：%w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("创建快捷方式失败：%v：%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// psEscape 对路径进行转义，使其可安全嵌入 PowerShell 单引号字符串中。
func psEscape(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
