package main

import (
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var cacheDir string

func initCache() {
	base, err := os.UserCacheDir()
	if err != nil {
		base = os.TempDir()
	}
	cacheDir = filepath.Join(base, "s26mini-lineageos")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		die(t("Não consegui criar a pasta %s: %v", "Could not create %s: %v"), cacheDir, err)
	}
}

func fileSHA256(path string, label string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	st, _ := f.Stat()
	h := sha256.New()
	var r io.Reader = f
	if label != "" && st.Size() > 50<<20 {
		p := &progress{label: label, total: st.Size()}
		r = &countingReader{r: f, p: p}
		defer p.finish()
	}
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

type countingReader struct {
	r io.Reader
	n int64
	p *progress
}

func (c *countingReader) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	c.n += int64(n)
	c.p.update(c.n)
	return n, err
}

var httpClient = &http.Client{Timeout: 0}

// fetch returns the path of d in the cache, downloading and verifying it if needed.
func fetch(d download) string {
	dest := filepath.Join(cacheDir, d.Name)
	if _, err := os.Stat(dest); err == nil {
		got, err := fileSHA256(dest, t("Conferindo  ", "Checking    ")+short(d.Name))
		if err == nil && got == d.SHA256 {
			ok(t("%s (já baixado, conferido)", "%s (already downloaded, verified)"), short(d.Name))
			return dest
		}
		os.Remove(dest)
	}
	var lastErr error
	for _, url := range d.URLs {
		for attempt := 1; attempt <= 3; attempt++ {
			if lastErr = downloadTo(url, dest+".part", d); lastErr == nil {
				break
			}
			warn(t("Falha no download (tentativa %d): %v", "Download failed (attempt %d): %v"), attempt, lastErr)
			time.Sleep(3 * time.Second)
		}
		if lastErr == nil {
			break
		}
	}
	if lastErr != nil {
		die(t("Não consegui baixar %s. Confira a internet e tente de novo.", "Could not download %s. Check your internet and try again."), d.Name)
	}
	got, err := fileSHA256(dest+".part", "")
	if err != nil || got != d.SHA256 {
		os.Remove(dest + ".part")
		die(t("O arquivo %s veio diferente do esperado (SHA-256 não confere). Por segurança, parei aqui.",
			"%s does not match its pinned SHA-256. Stopping for safety."), d.Name)
	}
	if err := os.Rename(dest+".part", dest); err != nil {
		die("%v", err)
	}
	ok(t("%s baixado e conferido", "%s downloaded and verified"), short(d.Name))
	return dest
}

func downloadTo(url, dest string, d download) error {
	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); strings.HasPrefix(ct, "text/html") {
		return fmt.Errorf("got a web page instead of the file (mirror hiccup)")
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	total := d.Size
	if total == 0 {
		total = resp.ContentLength
	}
	p := &progress{label: t("Baixando    ", "Downloading ") + short(d.Name), total: total}
	_, err = io.Copy(f, &countingReader{r: resp.Body, p: p})
	p.finish()
	return err
}

func short(name string) string {
	if len(name) > 34 {
		return name[:31] + "..."
	}
	return name
}

// gunzipImage turns the verified .img.gz into a .img next to it (once).
func gunzipImage(gz string, d download) string {
	img := strings.TrimSuffix(gz, ".gz")
	marker := img + ".ok"
	if b, err := os.ReadFile(marker); err == nil && strings.TrimSpace(string(b)) == d.SHA256 {
		if _, err := os.Stat(img); err == nil {
			ok(t("Imagem já descompactada", "Image already decompressed"))
			return img
		}
	}
	in, err := os.Open(gz)
	if err != nil {
		die("%v", err)
	}
	defer in.Close()
	st, _ := in.Stat()
	zr, err := gzip.NewReader(&countingReader{r: in, p: &progress{label: t("Descompactando", "Decompressing "), total: st.Size()}})
	if err != nil {
		die("%v", err)
	}
	out, err := os.Create(img + ".part")
	if err != nil {
		die("%v", err)
	}
	if _, err := io.Copy(out, zr); err != nil {
		out.Close()
		os.Remove(img + ".part")
		die(t("Falha ao descompactar (disco cheio?): %v", "Decompression failed (disk full?): %v"), err)
	}
	out.Close()
	fmt.Println()
	if err := os.Rename(img+".part", img); err != nil {
		die("%v", err)
	}
	os.WriteFile(marker, []byte(d.SHA256+"\n"), 0o644)
	ok(t("Imagem descompactada", "Image decompressed"))
	return img
}

// setupPlatformTools extracts adb and fastboot and returns the folder that holds them.
func setupPlatformTools() string {
	d, found := platformTools[runtime.GOOS]
	if !found {
		die(t("Sistema não suportado: %s", "Unsupported system: %s"), runtime.GOOS)
	}
	zipPath := fetch(d)
	dir := filepath.Join(cacheDir, "platform-tools-r37.0.1")
	if _, err := os.Stat(filepath.Join(dir, "platform-tools", exe("fastboot"))); err == nil {
		return filepath.Join(dir, "platform-tools")
	}
	zr, err := zip.OpenReader(zipPath)
	if err != nil {
		die("%v", err)
	}
	defer zr.Close()
	for _, f := range zr.File {
		target := filepath.Join(dir, f.Name)
		if !strings.HasPrefix(target, filepath.Clean(dir)+string(os.PathSeparator)) {
			continue // zip-slip guard
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0o755)
			continue
		}
		os.MkdirAll(filepath.Dir(target), 0o755)
		rc, err := f.Open()
		if err != nil {
			die("%v", err)
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode()|0o600)
		if err != nil {
			die("%v", err)
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return filepath.Join(dir, "platform-tools")
}

func exe(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
