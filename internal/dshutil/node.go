package dshutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// nodeInstallTimeout 限制等待 msiexec 完成安装的总时长。
	nodeInstallTimeout = 5 * time.Minute
	// msiFailGrace 是 msiexec 非零退出后仍继续等待的宽限时长（UAC 提权子进程可能仍在写盘）。
	msiFailGrace = 20 * time.Second
	// pollInterval 是等待安装完成时的轮询间隔。
	pollInterval = 2 * time.Second
	// probeTimeout 限制单次 node --version 探测时长。
	probeTimeout = 5 * time.Second
	// downloadTimeout 是下载安装包时 HTTP 客户端的总超时。
	downloadTimeout = 15 * time.Minute
)

// NodeInstalled 报告 node 是否可在当前 PATH 中被解析。
func NodeInstalled() bool {
	_, err := exec.LookPath("node")
	return err == nil
}

// NodeVersionString 返回已安装的 node 版本（如 "v24.20.0"）；不可用时返回
// 空字符串。
func NodeVersionString() string {
	out, code := OutputCmd("node", "--version")
	if code != 0 {
		return ""
	}
	return out
}

// nodeMSIURL 构造指定版本的 Node.js 安装包下载地址。
func nodeMSIURL(version string) string {
	return fmt.Sprintf("%s/v%s/node-v%s-x64.msi", NodeMirrorBase, version, version)
}

// InstallNode 下载并静默安装 Node.js 的每机器 MSI 安装包，随后轮询等待
// node 可用。调用方应确保进程已以管理员权限运行。msiexec 可能在安装尚未
// 真正完成时就返回，因此以"node 可用"作为安装完成的真正信号。
func InstallNode() error {
	msiPath := filepath.Join(os.TempDir(), fmt.Sprintf("node-v%s-x64.msi", NodeVersion))
	if err := downloadFile(nodeMSIURL(NodeVersion), msiPath, Stdout); err != nil {
		return err
	}
	defer os.Remove(msiPath)
	if err := validateMSI(msiPath); err != nil {
		return err
	}

	// /l*v 保留 msiexec 详细安装日志，安装再次失败时可据此定位根因。
	msiLog := filepath.Join(os.TempDir(), fmt.Sprintf("node-v%s-install.log", NodeVersion))
	Info("正在静默安装 Node.js %s，一般需要 1-2 分钟...", NodeVersion)
	code := RunCmd("msiexec", "/i", msiPath, "/qn", "/norestart", "/l*v", msiLog)

	deadline := time.Now().Add(nodeInstallTimeout)
	// msiexec 非零退出时给一个短窗口，之后仍不可用则快速失败，而不是空等满
	// 5 分钟。
	msiFailAt := time.Now().Add(msiFailGrace)
	for {
		// 注册表在安装完成那一刻才更新，因此每次轮询都从注册表重建 PATH，
		// 只刷新一次会读到旧值。
		RefreshPathEnv(NodeExtraDirs())
		if nodeReady() {
			_ = os.Remove(msiLog)
			return nil
		}
		if code != 0 && time.Now().After(msiFailAt) {
			return fmt.Errorf("Node.js 安装失败（msiexec 退出码 %d），详细日志：%s", code, msiLog)
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("等待 Node.js 安装完成超时，请手动检查安装，日志：%s", msiLog)
		}
		time.Sleep(pollInterval)
	}
}

// nodeReady 报告当前是否有可用的 node。
func nodeReady() bool {
	if !NodeInstalled() {
		return false
	}
	// 限制单次探测时间，避免 node.exe 半途损坏或被安全软件扫描阻塞而击穿
	// 外层超时。
	_, code := OutputCmdTimeout(probeTimeout, "node", "--version")
	return code == 0
}

// msiMagic 是所有 Windows Installer（.msi）文件开头的 OLE2 复合文档魔数，
// 用于识别镜像站返回的 HTML 错误页。
var msiMagic = [8]byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}

// validateMSI 校验 path 是非空的、带 OLE2 魔数的 MSI 文件。
func validateMSI(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("无法打开安装包：%w", err)
	}
	defer f.Close()
	head := make([]byte, len(msiMagic))
	n, err := io.ReadFull(f, head)
	if err != nil || n != len(msiMagic) {
		return fmt.Errorf("安装包不完整或损坏，请重新下载")
	}
	if !bytes.Equal(head, msiMagic[:]) {
		return fmt.Errorf("下载到的文件不是有效的 MSI 安装包，请重新下载")
	}
	return nil
}

// progressWriter 包装目标文件，在写入时按进度向 progress 上报百分比。
type progressWriter struct {
	w        io.Writer
	progress io.Writer
	total    int64
	written  int64
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n, err := p.w.Write(b)
	if n > 0 {
		p.written += int64(n)
		if p.total > 0 {
			fmt.Fprintf(p.progress, "\r  下载进度：%3d%%  (%d/%d MB)",
				p.written*100/p.total, p.written/1024/1024, p.total/1024/1024)
		}
	}
	return n, err
}

// downloadFile 将 url 下载到 dst，并通过 progress 上报进度（传 io.Discard
// 可关闭进度显示）。目标文件若已存在会被覆盖。
func downloadFile(url, dst string, progress io.Writer) error {
	_ = os.Remove(dst)

	client := &http.Client{Timeout: downloadTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("下载失败：%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败：HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建临时文件：%w", err)
	}
	defer file.Close()

	pw := &progressWriter{w: file, progress: progress, total: resp.ContentLength}
	if _, err := io.Copy(pw, resp.Body); err != nil {
		return fmt.Errorf("下载中断：%w", err)
	}
	if pw.total > 0 && pw.written != pw.total {
		return fmt.Errorf("下载不完整：已接收 %d/%d 字节", pw.written, pw.total)
	}
	fmt.Fprintf(progress, "\r  下载完成。%s\n", strings.Repeat(" ", 30))
	return nil
}
