package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	lang  = "pt"
	stdin = bufio.NewReader(os.Stdin)
)

// t picks the Portuguese or English text.
func t(pt, en string) string {
	if lang == "en" {
		return en
	}
	return pt
}

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	dim    = "\033[2m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	cyan   = "\033[36m"
)

var stepNo, stepTotal int

func step(title string) {
	stepNo++
	fmt.Printf("\n%s%s[%d/%d] %s%s\n", bold, cyan, stepNo, stepTotal, title, reset)
}

func say(format string, a ...any) { fmt.Printf("  "+format+"\n", a...) }
func ok(format string, a ...any) {
	fmt.Printf("  %s✔%s "+format+"\n", append([]any{green, reset}, a...)...)
}
func warn(format string, a ...any) {
	fmt.Printf("  %s⚠ "+format+"%s\n", append(append([]any{yellow}, a...), reset)...)
}

func big(format string, a ...any) {
	fmt.Printf("\n  %s%s"+format+"%s\n", append(append([]any{bold, yellow}, a...), reset)...)
}

// die prints a friendly error and waits for Enter so a double-clicked window doesn't vanish.
func die(format string, a ...any) {
	fmt.Printf("\n  %s%s✖ "+format+"%s\n", append(append([]any{bold, red}, a...), reset)...)
	say(t("Nada foi perdido se o erro aconteceu antes da gravação. Você pode rodar o instalador de novo.",
		"Nothing is lost if this happened before flashing. You can run the installer again."))
	pause()
	os.Exit(1)
}

func pause() {
	fmt.Printf("\n  %s%s%s", dim, t("Aperte Enter para continuar...", "Press Enter to continue..."), reset)
	stdin.ReadString('\n')
}

func readLine() string {
	s, err := stdin.ReadString('\n')
	if err != nil && s == "" {
		die(t("Entrada fechada.", "Input closed."))
	}
	return strings.TrimSpace(s)
}

// choose shows numbered options and returns the index picked (def when Enter alone).
func choose(question string, options []string, def int) int {
	fmt.Printf("\n  %s%s%s\n", bold, question, reset)
	for i, o := range options {
		mark := " "
		if i == def {
			mark = "›"
		}
		fmt.Printf("   %s %d) %s\n", mark, i+1, o)
	}
	for {
		fmt.Printf("  %s [%d]: ", t("Escolha", "Choose"), def+1)
		s := readLine()
		if s == "" {
			return def
		}
		var n int
		if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n >= 1 && n <= len(options) {
			return n - 1
		}
	}
}

func yesNo(question string, def bool) bool {
	hint := t("S/n", "Y/n")
	if !def {
		hint = t("s/N", "y/N")
	}
	for {
		fmt.Printf("\n  %s%s%s [%s]: ", bold, question, reset, hint)
		switch strings.ToLower(readLine()) {
		case "":
			return def
		case "s", "sim", "y", "yes":
			return true
		case "n", "nao", "não", "no":
			return false
		}
	}
}

// confirmWord makes the user type a word (e.g. APAGAR) before something destructive.
func confirmWord(word string) bool {
	fmt.Printf("\n  %s%s%s %s%s%s: ", bold, red, t("Para continuar, digite", "To continue, type"), word, reset, "")
	return strings.EqualFold(readLine(), word)
}

// progress draws a single-line bar; call with done=total at the end.
type progress struct {
	label string
	total int64
	last  time.Time
}

func (p *progress) update(done int64) {
	if time.Since(p.last) < 200*time.Millisecond && done != p.total {
		return
	}
	p.last = time.Now()
	if p.total <= 0 {
		fmt.Printf("\r  %s  %s", p.label, human(done))
		return
	}
	pct := float64(done) / float64(p.total)
	if pct > 1 {
		pct = 1
	}
	const width = 28
	n := int(pct * width)
	fmt.Printf("\r  %s  [%s%s] %3.0f%%  %s / %s", p.label,
		strings.Repeat("█", n), strings.Repeat("░", width-n), pct*100, human(done), human(p.total))
}

func (p *progress) finish() { fmt.Println() }

func human(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.0f MB", float64(b)/(1<<20))
	default:
		return fmt.Sprintf("%d KB", b/1024)
	}
}

// waitWithDots runs check every interval until it returns true or timeout expires.
func waitWithDots(label string, timeout, interval time.Duration, check func() bool) bool {
	fmt.Printf("  %s ", label)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if check() {
			fmt.Printf(" %s✔%s\n", green, reset)
			return true
		}
		fmt.Print(".")
		time.Sleep(interval)
	}
	fmt.Println()
	return false
}
