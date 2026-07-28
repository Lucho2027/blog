package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	inPath := flag.String("in", "", "input TypeScript file (sqlc generated)")
	outPath := flag.String("out", "", "output TypeScript file (types-only)")
	keepTypes := flag.Bool("keepTypeAliases", true, "also keep `export type ...` blocks (best-effort)")
	flag.Parse()

	if *inPath == "" || *outPath == "" {
		fmt.Fprintln(os.Stderr, "Usage: extract_ts_types -in <generated.ts> -out <types-only.ts>")
		os.Exit(2)
	}

	in, err := os.Open(*inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open input: %v\n", err)
		os.Exit(1)
	}

	out, err := os.Create(*outPath)
	if err != nil {
		_ = in.Close()
		fmt.Fprintf(os.Stderr, "create output: %v\n", err)
		os.Exit(1)
	}

	w := bufio.NewWriter(out)

	// Header
	fmt.Fprintln(w, "/*")
	fmt.Fprintln(w, " * AUTO-GENERATED (types-only) from sqlc output.")
	fmt.Fprintln(w, " * Do not edit by hand.")
	fmt.Fprintln(w, " */")
	fmt.Fprintln(w)

	sc := bufio.NewScanner(in)

	inBlock := false
	blockBraceDepth := 0

	shouldStartBlock := func(line string) bool {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "export interface ") {
			return true
		}
		if *keepTypes && strings.HasPrefix(line, "export type ") {
			return true
		}
		return false
	}

	for sc.Scan() {
		line := sc.Text()
		trim := strings.TrimSpace(line)

		if strings.HasPrefix(trim, "import ") && strings.Contains(trim, `"postgres"`) {
			continue
		}
		if !inBlock && strings.HasPrefix(trim, "import ") {
			continue
		}

		if !inBlock {
			if shouldStartBlock(line) {
				inBlock = true
				blockBraceDepth = strings.Count(line, "{") - strings.Count(line, "}")
				fmt.Fprintln(w, line)

				if blockBraceDepth <= 0 && strings.HasPrefix(strings.TrimSpace(line), "export type ") {
					inBlock = false
					blockBraceDepth = 0
					fmt.Fprintln(w)
				}
			}
			continue
		}

		fmt.Fprintln(w, line)
		blockBraceDepth += strings.Count(line, "{") - strings.Count(line, "}")

		if blockBraceDepth <= 0 && strings.Contains(trim, "}") {
			inBlock = false
			blockBraceDepth = 0
			fmt.Fprintln(w)
		}
	}

	// 1) check scan errors first
	if err := sc.Err(); err != nil {
		_ = in.Close()
		_ = w.Flush()
		_ = out.Close()
		fmt.Fprintf(os.Stderr, "scan input: %v\n", err)
		os.Exit(1)
	}

	// 2) flush + close outputs so file is fully written
	if err := w.Flush(); err != nil {
		_ = in.Close()
		_ = out.Close()
		fmt.Fprintf(os.Stderr, "flush output: %v\n", err)
		os.Exit(1)
	}
	if err := out.Close(); err != nil {
		_ = in.Close()
		fmt.Fprintf(os.Stderr, "close output: %v\n", err)
		os.Exit(1)
	}

	// 3) close input before removing (Windows-safe)
	if err := in.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "close input: %v\n", err)
		os.Exit(1)
	}

	// 4) now remove the sqlc-generated file
	if err := os.Remove(*inPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(os.Stderr, "remove input: %v\n", err)
			os.Exit(1)
		}
	}
}
