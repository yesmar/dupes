// SPDX-License-Identifier: BSD-2-Clause
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fatih/color"
)

var (
	cmd        = path.Base(os.Args[0])
	release    string
	failMsg    = color.New(color.FgRed).FprintfFunc()
	statusMsg  = color.New(color.FgCyan).FprintfFunc()
	successMsg = color.New(color.FgGreen).FprintfFunc()
)

type runtime struct {
	noColor bool
	verbose bool
	files   map[string]string
}

func hashFile(pathname string, rt runtime) (hash []byte, ok bool) {
	f, err := os.Open(pathname)
	if err != nil {
		msg := strings.Split(err.Error(), ":")
		err = errors.New(pathname + ":" + msg[1])
		failMsg(os.Stderr, "%s\n", err)
		return nil, false
	}
	defer func() {
		if err := f.Close(); err != nil {
			failMsg(os.Stderr, "%s\n", err)
			hash = nil
			ok = false
		}
	}()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		failMsg(os.Stderr, "%s\n", err)
		return nil, false
	}
	return hasher.Sum(nil), true
}

func processFile(fi os.FileInfo, pathname string, rt runtime) bool {
	if fi.Mode().IsRegular() && fi.Size() > 0 {
		h, ok := hashFile(pathname, rt)
		if ok {
			s := hex.EncodeToString(h)
			f, ok := rt.files[s]
			if ok {
				if rt.verbose {
					fmt.Fprintf(os.Stderr, "%s duplicates %s → ", pathname, f) // nolint: gas
					statusMsg(os.Stderr, "%s\n", s)
				} else {
					fmt.Println(pathname)
				}
			} else {
				rt.files[s] = pathname
			}
			return true
		}
	}
	return false
}

func processDirectory(dirname string, rt runtime) (count uint, err error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		msg := strings.Split(err.Error(), ":")
		err = errors.New(dirname + ":" + msg[1])
		failMsg(os.Stderr, "%s\n", err)
	}
	for _, fi := range files {
		pathname := path.Join(dirname, fi.Name())
		if fi.IsDir() {
			// Process directory.
			c, err := processTarget(pathname, rt)
			if err != nil {
				failMsg(os.Stderr, "%s\n", err)
			}
			count += c
		} else {
			if processFile(fi, pathname, rt) {
				count++
			}
		}
	}
	return count, nil
}

func processTarget(pathname string, rt runtime) (count uint, err error) {
	fi, err := os.Stat(pathname)
	if err != nil {
		// This is safe because func (*PathError) Error takes the form "{op} {path}: {error}".
		// I'm doing this because I don't want the {op} part of the error message.
		msg := strings.Split(err.Error(), ":")
		err = errors.New(pathname + ":" + msg[1])
		failMsg(os.Stderr, "%s\n", err)
		return 0, err
	}

	// Process directory.
	if fi.IsDir() {
		count, err = processDirectory(pathname, rt)
		if err != nil {
			failMsg(os.Stderr, "%s\n", err)
		}
		return count, nil
	}

	// Process file.
	if processFile(fi, pathname, rt) {
		count++
	}

	return count, nil
}

func main() {
	var flagNoColor = flag.Bool("no-color", false, "Disable color output")
	var flagVerbose = flag.Bool("verbose", false, "Verbose output")
	var flagVersion = flag.Bool("version", false, "Display version information")

	flag.Parse()

	if *flagVersion {
		statusMsg(os.Stderr, "%s %s\n", cmd, release)
		os.Exit(0)
	}

	var rt runtime
	rt.noColor = *flagNoColor
	rt.verbose = *flagVerbose
	rt.files = make(map[string]string)

	if rt.noColor {
		color.NoColor = true
	}

	for _, target := range flag.Args() {
		count, err := processTarget(target, rt)
		if err != nil {
			os.Exit(1) // Error message has already been displayed.
		}

		if rt.verbose {
			successMsg(os.Stderr, "%s: processed %d file", cmd, count)
			if count > 1 {
				successMsg(os.Stderr, "s")
			}
			successMsg(os.Stderr, " in %s\n", target)
		}
	}

	os.Exit(0)
}
