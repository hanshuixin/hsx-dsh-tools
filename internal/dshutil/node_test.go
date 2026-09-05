package dshutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNodeMSIURL(t *testing.T) {
	got := nodeMSIURL("24.20.0")
	want := "https://npmmirror.com/mirrors/node/v24.20.0/node-v24.20.0-x64.msi"
	if got != want {
		t.Fatalf("nodeMSIURL() = %q, want %q", got, want)
	}
}

func TestDownloadFile(t *testing.T) {
	payload := bytes.Repeat([]byte("dsh-test-"), 1000)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	dst := filepath.Join(t.TempDir(), "out.bin")
	var progress bytes.Buffer
	if err := downloadFile(srv.URL, dst, &progress); err != nil {
		t.Fatalf("downloadFile() error = %v", err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, payload) {
		t.Fatal("downloaded content mismatch")
	}
	if !strings.Contains(progress.String(), "下载完成") {
		t.Fatalf("progress output missing completion marker: %q", progress.String())
	}
}

func TestDownloadFileHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusNotFound)
	}))
	defer srv.Close()

	err := downloadFile(srv.URL, filepath.Join(t.TempDir(), "out.bin"), io.Discard)
	if err == nil {
		t.Fatal("expected error for HTTP 404")
	}
	if !strings.Contains(err.Error(), "HTTP 404") {
		t.Fatalf("error should mention HTTP 404, got: %v", err)
	}
}

func TestDownloadFileTruncated(t *testing.T) {
	payload := bytes.Repeat([]byte("x"), 64)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 声明 1MB 但只发送 64 字节，然后关闭连接。
		w.Header().Set("Content-Length", fmt.Sprint(1024*1024))
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	err := downloadFile(srv.URL, filepath.Join(t.TempDir(), "out.bin"), io.Discard)
	if err == nil {
		t.Fatal("expected error for truncated download")
	}
	// 连接被提前关闭时 Go 客户端报 unexpected EOF（走“下载中断”分支），
	// Content-Length 不匹配分支作为兜底；两者都算截断检测成功。
	if !strings.Contains(err.Error(), "下载") {
		t.Fatalf("error should mention download failure, got: %v", err)
	}
}

func TestValidateMSI(t *testing.T) {
	dir := t.TempDir()

	good := filepath.Join(dir, "good.msi")
	head := append(append([]byte{}, msiMagic[:]...), bytes.Repeat([]byte{0}, 24)...)
	if err := os.WriteFile(good, head, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateMSI(good); err != nil {
		t.Fatalf("valid MSI rejected: %v", err)
	}

	bad := filepath.Join(dir, "bad.msi")
	if err := os.WriteFile(bad, []byte("<html>not found</html>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := validateMSI(bad); err == nil {
		t.Fatal("invalid MSI accepted")
	}
}
