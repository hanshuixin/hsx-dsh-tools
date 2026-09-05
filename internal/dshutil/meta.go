package dshutil

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

// Version 是三个工具的发行版本号，也用于 exe 的 --version 输出。
// 版本号的唯一来源是同目录的 VERSION 文件：本包在编译时经 go:embed 嵌入，
// build.ps1 打包 zip 时也会读取该文件，因此发版只需修改 VERSION 一处。
var Version = strings.TrimSpace(versionFile)

// 本文件集中管理三个工具除版本号外的“身份 / 配置”常量（可执行文件名、
// 快捷方式、Web 地址、npm 包与镜像等），避免同一份信息散落多处、难以统一
// 修改。发版或调整这些外部依赖时，优先改这里。
const (
	// 三个工具的可执行文件名（start/upgrade 提示、install 建快捷方式时相互引用）。
	ExeInstall = "安装DeepSeekHarness.exe"
	ExeStart   = "启动DeepSeekHarness.exe"
	ExeUpgrade = "升级DeepSeekHarness.exe"

	// ShortcutName 是桌面快捷方式的显示名，ShortcutDescription 是其描述。
	ShortcutName        = "启动 DeepSeek Harness"
	ShortcutDescription = "启动 DeepSeek Harness Web 界面"

	// WebURL 是 dsh web 默认监听的服务地址。
	WebURL = "http://127.0.0.1:3080"

	// DshPackage 是 DeepSeek Harness 的 npm 包名。
	DshPackage = "@deepseek-ai/dsh"

	// NpmRegistry 是国内 npm 镜像源，用于加速下载与安装。
	NpmRegistry = "https://registry.npmmirror.com"

	// NodeVersion 是缺失时安装的 Node.js LTS 版本。
	NodeVersion = "24.20.0"
	// NodeMirrorBase 是 Node.js 安装包的 npmmirror 镜像基址。
	NodeMirrorBase = "https://npmmirror.com/mirrors/node"
)
