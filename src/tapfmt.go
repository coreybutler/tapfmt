package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"tapfmt/json"
	"tapfmt/spec"

	"github.com/coreybutler/go-timer"
	color "github.com/logrusorgru/aurora/v3"
	tap "github.com/mpontillo/tap13"
)

type Formatter interface {
	Format(results *tap.Results)
	Summary()
}

func main() {
	// Identify formatter type
	format := flag.String("f", "spec", "format")
	flag.Parse()

	var formatter Formatter
	switch strings.ToLower(strings.TrimSpace(*format)) {
	case "spec":
		formatter = spec.Formatter()
	case "json":
		formatter = json.Formatter()
	default:
		fmt.Printf("\n\n  %s %v\n  %s\n", color.BrightMagenta("\u26A0 Warning: unrecognized formatter,"), color.Bold(color.Magenta(*format)), color.Italic(color.Faint("  (using default formatter instead)")))
		formatter = spec.Formatter()
	}

	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)

	monitor := timer.SetTimeout(func(args ...interface{}) {
		fmt.Println("No input stream to format.")
		fmt.Println(color.Faint("usage: myprocess | tapfmt [-f spec]"))
		os.Exit(0)
	}, 1000)

	for {
		n, err := r.Read(buf[:cap(buf)])
		monitor.Stop()
		buf = buf[:n]

		if n == 0 {
			if err == nil {
				fmt.Println("yo2")
				continue
			}
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if err != nil && err != io.EOF {
			panic(err)
		}

		results := tap.Parse(strings.Split(fmt.Sprintf("%s", buf), "\n"))
		formatter.Format(results)
	}

	formatter.Summary()
}
