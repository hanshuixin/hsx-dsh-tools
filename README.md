<p align="center">
  <img src="assets/hsx-dsh-tools-logo-256.png" alt="hsx-dsh-tools 标志" width="200">
</p>

# hsx-dsh-tools — 寒水馨 DeepSeek Harness 工具集

<p align="center">
  <img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-green?style=flat-square">
  <img alt="Go" src="https://img.shields.io/badge/language-Go-00ADD8?style=flat-square">
  <img alt="Windows" src="https://img.shields.io/badge/platform-Windows-0078D6?style=flat-square">
  <img alt="Website" src="https://img.shields.io/badge/website-hanshuixin.org-blue?style=flat-square">
</p>

> **hsx-dsh-tools** 是一键安装、启动与升级 [DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness)（`dsh`）的工具集。

> ⚠️ **当前仅支持 Windows**（Windows 10 / 11，x86_64）。macOS 与 Linux 版本暂未提供。

## 简介

[DeepSeek Harness](https://github.com/deepseek-ai/deepseek-harness)（命令行工具名 `dsh`）是深度求索（DeepSeek）官方推出的开源智能体（Agent）开发框架，通过 npm 包 `@deepseek-ai/dsh` 分发。它通常需要在命令行里一步步配置 Node.js、安装并启动服务，对新手不太友好。

**hsx-dsh-tools** 把这三步封装成三个双击即可的 exe，让没有任何 Node.js 经验的用户也能一键安装、一键启动、一键升级。

## 特性

- **自动安装 Node.js**：检测到系统缺少 Node.js 时，自动从国内镜像下载 LTS 版（Node 24.x）并静默安装，无需手动配置。
- **国内镜像加速**：npm 自动切换到国内镜像源，下载与安装更快更稳。
- **全中文界面**：控制台彩色输出（`[OK]` 绿 / `[警告]` 黄 / `[失败]` 红），安装过程中的 UAC 提示会明确提醒点击"是"。
- **桌面快捷方式**：安装完成后自动在桌面创建「启动 DeepSeek Harness」快捷方式，双击即可使用。
- **安装可靠**：自动处理 UAC 提权后的等待、校验下载文件完整性，避免装到一半失败。
- **免安装、自带图标**：三个 exe 均为独立程序，解压即可用，并内嵌了应用图标。

## 快速开始

### 1. 下载 hsx-dsh-tools

Windows 发行包下载地址：

<https://hanshuixin.org/resource/software_integrated_package/Windows/hsx-dsh-tools>

下载后会得到一个 zip 包，包含 3 个 exe 可执行文件。

### 2. 安装 hsx-dsh-tools（绿色免安装包）

本工具集为**绿色免安装包**：不需要安装程序、不写注册表。把 zip 解压后得到的整个文件夹放到任意你喜欢的位置即可，例如：

```
C:\Program Files\hsx-dsh-tools\
```

日后如要卸载，直接删除该文件夹即可。

### 3. 使用 hsx-dsh-tools

文件夹内包含三个可执行程序，按需双击即可（安装与升级需要联网）。三个 exe 的简要说明如下：

| 可执行文件 | 说明 |
|-----------|------|
| `安装DeepSeekHarness.exe` | 首次使用先运行：装好 Node.js 与 `@deepseek-ai/dsh`，并创建桌面快捷方式 |
| `启动DeepSeekHarness.exe` | 日常使用：启动 DeepSeek Harness 的 Web 界面 |
| `升级DeepSeekHarness.exe` | 有新版时运行：将 `@deepseek-ai/dsh` 升级到最新版 |

> 三个程序均为控制台程序：运行后请保持窗口打开，**关闭窗口即停止服务**。

#### 3.1 安装 DeepSeek Harness

双击 `安装DeepSeekHarness.exe` 

定位：首次部署，一键准备好运行环境。

运行后会自动完成以下工作：

1. 检测系统是否已安装 Node.js，缺失时自动从国内镜像下载 LTS 版并静默安装；
2. 将 npm 切换到国内镜像源；
3. 全局安装 `@deepseek-ai/dsh`；
4. 在桌面创建「启动 DeepSeek Harness」快捷方式。

说明：全程需要联网；安装 Node.js 时如弹出 UAC 提示，请点击"是"；完成一次即可，如运行环境异常可再次运行修复。

#### 3.2 启动 DeepSeek Harness

双击 `启动DeepSeekHarness.exe`

定位：日常使用，启动 DeepSeek Harness 的 Web 界面。

运行后检查 `dsh` 是否已安装并显示版本，然后启动 Web 界面。

启动后会自动打开 DeepSeek Harness 的 Web 页面。

使用前：你需要先在 [DeepSeek 开放平台](https://platform.deepseek.com/api_keys) 创建一个 API key。

在页面的弹窗中填入刚才创建的 API key 即可使用。

保持窗口打开即服务运行中，关闭窗口即停止。

#### 3.3 升级 DeepSeek Harness

双击 `升级DeepSeekHarness.exe`

定位：将 dsh 保持为最新版。

运行后自动将 `@deepseek-ai/dsh` 升级到 npm 上的最新版并验证，全程需要联网。

## 常见问题

- **双击启动后提示"未检测到 @deepseek-ai/dsh"**：先运行 `安装DeepSeekHarness.exe` 完成安装。
- **安装 Node.js 时卡住或报错**：多为网络问题，请检查网络后重试，或访问 [npmmirror Node 下载页](https://npmmirror.com/mirrors/node/) 手动下载对应版本的 `.msi` 安装。
- **Web 界面打不开**：确认启动窗口仍然开着（关闭窗口 = 停止服务），并确认访问地址是 `http://127.0.0.1:3080`。
- **网页提示没有 API key**：在 [DeepSeek 开放平台](https://platform.deepseek.com/api_keys) 创建 API key，并在 Web 界面的模型设置中填入。

## 反馈与贡献

本项目由作者**个人维护**，详情见 [CONTRIBUTING.md](./CONTRIBUTING.md)：

- **欢迎**：通过 **Issues** 提建议、报告问题（发布后会提供链接）。
- **不接受外部 Pull Request**：请勿直接提交 PR。

## 开发者

从源码构建、目录结构、测试、发版与版本号管理等，见 [docs/DEVELOPMENT.md](./docs/DEVELOPMENT.md)。

## 开源许可

本项目依据 **MIT 许可证** 授权，作者 **寒水馨**（[关于作者 / About](https://hanshuixin.org/about)），网站首页 [hanshuixin.org](https://hanshuixin.org)。完整许可证文本见 [LICENSE](./LICENSE)。

> 本项目所安装的 DeepSeek Harness 软件版权归属于 DeepSeek（深度求索），其本身依据 MIT 许可证授权。
