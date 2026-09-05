package dshutil

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// systemEnvPathKey 是系统级环境变量的注册表键。
const systemEnvPathKey = `SYSTEM\CurrentControlSet\Control\Session Manager\Environment`

// userEnvPathKey 是用户级环境变量的注册表键。
const userEnvPathKey = `Environment`

// MergePath 将若干条 PATH 字符串合并成一条去重的 PATH 字符串。条目会去除
// 首尾空白与尾部反斜杠，大小写不敏感去重，并保持首次出现的顺序。
func MergePath(paths ...string) string {
	seen := make(map[string]struct{})
	var out []string
	for _, path := range paths {
		for _, entry := range strings.Split(path, ";") {
			entry = normalizePathEntry(entry)
			if entry == "" {
				continue
			}
			key := strings.ToLower(entry)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, entry)
		}
	}
	return strings.Join(out, ";")
}

// normalizePathEntry 规范化单个 PATH 条目：去除首尾空白与尾部反斜杠，
// 但盘符根路径（如 C:\）会保留根反斜杠。
func normalizePathEntry(entry string) string {
	entry = strings.TrimSpace(entry)
	trimmed := strings.TrimRight(entry, `\`)
	if trimmed != entry && len(trimmed) == 2 && trimmed[1] == ':' {
		return trimmed + `\`
	}
	return trimmed
}

// RegistryPath 读取指定根键与子键下的 "Path" 值，并展开 REG_EXPAND_SZ
// 类型的变量（如 %SystemRoot%），使结果可直接用于 PATH。键或值不存在时
// 返回空字符串。
func RegistryPath(root registry.Key, subkey string) string {
	k, err := registry.OpenKey(root, subkey, registry.QUERY_VALUE)
	if err != nil {
		return ""
	}
	defer k.Close()
	v, typ, err := k.GetStringValue("Path")
	if err != nil {
		return ""
	}
	if typ == registry.EXPAND_SZ {
		if expanded, err := registry.ExpandString(v); err == nil {
			v = expanded
		}
	}
	return v
}

// SystemPath 从注册表返回系统级 PATH。
func SystemPath() string {
	return RegistryPath(registry.LOCAL_MACHINE, systemEnvPathKey)
}

// UserPath 从注册表返回用户级 PATH。
func UserPath() string {
	return RegistryPath(registry.CURRENT_USER, userEnvPathKey)
}

// NodeExtraDirs 返回安装 Node.js 后需要加入 PATH 的目录。
func NodeExtraDirs() []string {
	progFiles := os.Getenv("ProgramFiles")
	if progFiles == "" {
		progFiles = `C:\Program Files`
	}
	dirs := []string{filepath.Join(progFiles, "nodejs")}
	if appData := os.Getenv("APPDATA"); appData != "" {
		dirs = append(dirs, filepath.Join(appData, "npm"))
	}
	return dirs
}

// RefreshPathEnv 用注册表与继承环境重建当前进程的 PATH，并追加 extraDirs，
// 使刚安装的工具无需重启控制台即可被解析。
func RefreshPathEnv(extraDirs []string) {
	merged := MergePath(os.Getenv("PATH"), SystemPath(), UserPath(), strings.Join(extraDirs, ";"))
	_ = os.Setenv("PATH", merged)
}
