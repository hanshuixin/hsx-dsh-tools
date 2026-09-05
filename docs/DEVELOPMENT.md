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
- [图标 / 版本 / 语言资源](#图标--版本--语言资源)

## 环境要求

- Windows 10 / 11（64 位）
- Go 1.24 及以上（本仓库 `go.mod` 声明 `go 1.24.0`，默认使用 vendor 模式，无需联网下载依赖）
- [go-winres](https://github.com/tc-hib/go-winres)（打包 exe 的图标 / 版本 / 语言资源时需要；一次性安装：
  `go install github.com/tc-hib/go-winres@latest`，并确保 `GOPATH\bin` 在 PATH 中）

## 一键构建

在终端（Windows PowerShell / CMD）中，于仓库根目录执行：

```powershell
powershell -ExecutionPolicy Bypass -File build.ps1
```

> 提示：`build.ps1` 双击默认会用记事本打开（Windows 不默认关联 ps1 执行），请改用上面的命令运行。

脚本会依次完成：

1. 依据各 `cmd/*/winres.json` 模板，用 go-winres 生成 exe 的图标 / 版本 / 语言资源（版本取自 `internal/dshutil/VERSION`，见下文「图标 / 版本 / 语言资源」）
2. 构建三个 exe 到 `dist/`（`-trimpath` + `-ldflags "-s -w"` 精简体积）
3. 复制 `LICENSE` 到 `dist/`
4. 用 PowerShell 将三个 exe 与 `LICENSE` 打包为 `hsx-dsh-tools-windows-x64-<版本号>.zip`

> 发行 zip 只包含三个 exe 与 `LICENSE`，不含 README / 源码等其它文件。

发行 zip 的命名规则详见 `build.ps1` 顶部注释。

## 手动构建

```bat
go build -trimpath -ldflags="-s -w" -o dist\安装DeepSeekHarness.exe .\cmd\dsh-install
go build -trimpath -ldflags="-s -w" -o dist\启动DeepSeekHarness.exe .\cmd\dsh-start
go build -trimpath -ldflags="-s -w" -o dist\升级DeepSeekHarness.exe .\cmd\dsh-upgrade
```

> 上面的命令只适合快速验证逻辑：生成的 exe 不含图标 / 版本 / 语言资源。
> 需要完整资源（并确保 exe 属性中的版本与 `VERSION` 一致）时，请用 [一键构建](#一键构建) 的 `build.ps1`。

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
│   │   └── winres.json         # exe 图标/版本/语言资源模板（由 build.ps1 编译）
│   ├── dsh-start/              # 启动程序入口
│   │   └── winres.json         # 同上（三个 exe 各自一份）
│   └── dsh-upgrade/            # 升级程序入口
│       └── winres.json
├── internal/
│   └── dshutil/                # 三个程序共享的公共库
│       ├── VERSION             # 发行版本号（全项目唯一来源，go:embed / build.ps1 共用）
│       ├── meta.go             # 版本号等身份/配置常量（统一集中管理）
│       ├── console.go          # UTF-8 控制台、彩色输出、窗口标题、按任意键
│       ├── elevate.go          # 管理员权限检测、UAC 自提权重启（安装/升级用）
│       ├── runner.go           # 外部命令执行（流式输出、退出码、超时）
│       ├── path.go             # 注册表 PATH 读取与合并
│       ├── node.go             # Node.js 检测 / 下载 / 安装 / 等待完成
│       ├── npm.go              # npm 镜像源与全局包操作
│       ├── dsh.go              # dsh 检测、版本、启动（含 npx 兜底）
│       └── shortcut.go         # 桌面快捷方式创建
├── assets/
│   ├── hsx-dsh-tools-logo.png      # 图标源文件（设计稿）
│   ├── hsx-dsh-tools-logo-256.png  # 缩小版标志（README 展示 / exe 内嵌图标共用）
│   └── hsx-dsh-tools-logo.ico      # 多分辨率图标（历史产物，当前 exe 已改用上述 PNG）
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
- 打包侧：`build.ps1` 读取同一文件，用作发行 zip 文件名中的版本号，并在编译资源时把版本注入 exe 的右键属性（文件 / 产品版本，四段式，如 `0.1.1` → `0.1.1.0`）。

**发版前记得改版本号**：给仓库打 tag 之前，先修改 `internal/dshutil/VERSION` 的内容（只含版本号，如 `0.1.1`），然后重新执行 `build.ps1` 打包。这样 zip 文件名、三个 exe 的 `--version` 与右键属性里的版本号三者保持一致，无需任何手动重生成步骤。

## 图标 / 版本 / 语言资源

exe 右键属性中的"版本 / 产品版本 / 文件说明 / 语言"等，来自 Windows 的 **VERSIONINFO 资源**（连同图标、manifest），由 [go-winres](https://github.com/tc-hib/go-winres) 生成并链接进 exe：

- **模板**：每个程序一份 `cmd/<x>/winres.json`（如 `cmd/dsh-install/winres.json`），其中写明了内嵌图标、VERSIONINFO 各字段（公司名 / 文件说明 / 版权等）、语言，以及 manifest（`execution-level: as invoker`，不做 UAC 提权——提权由 `elevate.go` 在代码层按需触发）。
- **语言**：`winres.json` 里 `info` 的语言键为 `0804`（中文简体），因此 exe 属性显示"语言：中文(简体,中国)"。
- **版本**：`winres.json` 中版本用占位符 `0.0.0.0` 表示，`build.ps1` 构建时把 `VERSION` 的值（补足四段）替换进去后再用 `go-winres make` 编译，因此**改 `VERSION` 重建即可，无需手动改资源**。
- **产物**：编译得到的 `cmd/<x>/rsrc_windows_amd64.syso` 与临时 `cmd/<x>/winres/` 均由构建自动生成，已加入 `.gitignore`，不要手动提交。

日常修改建议：

- 改版本号 → 只改 `internal/dshutil/VERSION`；
- 改图标 / 文件说明 / 公司名 / 版权等展示信息 → 改对应 `cmd/<x>/winres.json`；
- 若想更换内嵌图标，把新 PNG 放进 `assets/` 并更新各 `winres.json` 中 `RT_GROUP_ICON` 指向的文件（路径相对 `winres.json` 所在目录解析）。
