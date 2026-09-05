# 开发者文档

本文档面向**开发者 / 维护者**，介绍如何从源码构建、目录结构、测试与发版。
普通用户请阅读 [README.md](../README.md)。

## 目录

- [环境要求](#环境要求)
- [一键构建](#一键构建)
- [手动构建](#手动构建)
- [测试](#测试)
- [目录结构](#目录结构)
- [版本号管理](#版本号管理)
- [图标资源（重新生成）](#图标资源重新生成)

## 环境要求

- Windows 10 / 11（64 位）
- Go 1.24 及以上（本仓库 `go.mod` 声明 `go 1.24.0`，默认使用 vendor 模式，无需联网下载依赖）

## 一键构建

在终端（Windows PowerShell / CMD）中，于仓库根目录执行：

```powershell
powershell -ExecutionPolicy Bypass -File build.ps1
```

> 提示：`build.ps1` 双击默认会用记事本打开（Windows 不默认关联 ps1 执行），请改用上面的命令运行。

脚本会依次完成：

1. 构建三个 exe 到 `dist/`（`-trimpath` + `-ldflags "-s -w"` 精简体积）
2. 复制 `LICENSE` 到 `dist/`
3. 用 PowerShell 将三个 exe 与 `LICENSE` 打包为 `hsx-dsh-tools-windows-x64-0.1.0.zip`

> 发行 zip 只包含三个 exe 与 `LICENSE`，不含 README / 源码等其它文件。

发行 zip 的命名规则详见 `build.ps1` 顶部注释。

## 手动构建

```bat
go build -trimpath -ldflags="-s -w" -o dist\安装DeepSeekHarness.exe .\cmd\dsh-install
go build -trimpath -ldflags="-s -w" -o dist\启动DeepSeekHarness.exe .\cmd\dsh-start
go build -trimpath -ldflags="-s -w" -o dist\升级DeepSeekHarness.exe .\cmd\dsh-upgrade
```

## 测试

```bat
go test ./...
```

对公共库的纯逻辑（PATH 合并、文件下载、URL 构造、MSI 校验、PowerShell 转义）提供单元测试。
安装 / 启动 / 升级涉及真实下载、UAC 提权与外部进程，建议按 README 的"下载与使用"一节在真实环境中手动验证。

## 目录结构

```
hsx-dsh-tools/
├── cmd/
│   ├── dsh-install/            # 安装程序入口
│   ├── dsh-start/              # 启动程序入口
│   └── dsh-upgrade/            # 升级程序入口
├── internal/
│   └── dshutil/                # 三个程序共享的公共库
│       ├── VERSION             # 发行版本号（全项目唯一来源，go:embed / build.ps1 共用）
│       ├── meta.go             # 版本号等身份/配置常量（统一集中管理）
│       ├── console.go          # UTF-8 控制台、彩色输出、窗口标题、按任意键
│       ├── runner.go           # 外部命令执行（流式输出、退出码、超时）
│       ├── path.go             # 注册表 PATH 读取与合并
│       ├── node.go             # Node.js 检测 / 下载 / 安装 / 等待完成
│       ├── npm.go              # npm 镜像源与全局包操作
│       ├── dsh.go              # dsh 检测、版本、启动（含 npx 兜底）
│       └── shortcut.go         # 桌面快捷方式创建
├── assets/
│   ├── hsx-dsh-tools-logo.png      # 图标源文件（设计稿）
│   ├── hsx-dsh-tools-logo-256.png  # 缩小版标志（README 展示用）
│   └── hsx-dsh-tools-logo.ico      # 多分辨率图标（由 PNG 生成，用于嵌入 exe）
├── build.ps1                   # 一键构建脚本
├── go.mod / go.sum
├── vendor/                     # 依赖（vendor 模式）
├── docs/
│   └── DEVELOPMENT.md          # 本文档
├── README.md
├── CONTRIBUTING.md             # 贡献指南（欢迎建议、暂不接受 PR）
├── .gitattributes              # Git 行尾/文本属性规范
└── LICENSE
```

## 版本号管理

版本号的**唯一来源**是 `internal/dshutil/VERSION` 文件，全项目共用：

- Go 侧：`internal/dshutil/meta.go` 用 `go:embed` 嵌入该文件，导出为 `dshutil.Version`，三个 exe 的 `--version` 都取自它；
- 打包侧：`build.ps1` 读取同一文件，用作发行 zip 文件名中的版本号。

**发版前记得改版本号**：给仓库打 tag 之前，先修改 `internal/dshutil/VERSION` 的内容（只含版本号，如 `0.1.0`），然后重新执行 `build.ps1` 打包，确保 zip 名与 exe 的 `--version` 一致。

> 补充：exe 图标资源里的文件/产品版本（`0.1.0.0`）是单独由 go-winres 生成的，发版如需同步，见下文「图标资源（重新生成）」。

## 图标资源（重新生成）

`.syso` 图标资源文件已生成并提交到仓库，日常构建无需额外步骤。若更换图标或调整版本信息，可用 [go-winres](https://github.com/tc-hib/go-winres) 重新生成：

```bat
go-winres simply --arch amd64 --icon assets\hsx-dsh-tools-logo.ico ^
  --product-version 0.1.0.0 --file-version 0.1.0.0 ^
  --original-filename 安装DeepSeekHarness.exe ...
```

分别在 `cmd/dsh-install`、`cmd/dsh-start`、`cmd/dsh-upgrade` 三个目录下执行。
