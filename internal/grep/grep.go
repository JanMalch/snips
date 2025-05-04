package grep

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var idStyle = lipgloss.NewStyle().Bold(true)
var matchStyle = lipgloss.NewStyle()

// Source: https://github.com/healeycodes/tools/blob/005aea0a2f82a8f88c0d0a9f8dcfc9a534659713/grup/utils/search.go

type SearchDebug struct {
	Workers int
}

const (
	LITERAL = iota
	REGEX
)

type SearchOptions struct {
	Kind   int
	Lines  bool
	Regex  *regexp.Regexp
	Finder *stringFinder
}

type searchJob struct {
	path string
	opts *SearchOptions
}

func Search(path string, opts *SearchOptions, debug *SearchDebug, stdout io.Writer) {
	searchJobs := make(chan *searchJob)

	var wg sync.WaitGroup
	for w := 0; w < debug.Workers; w++ {
		go searchWorker(searchJobs, &wg, path, stdout)
	}
	dirTraversal(path, opts, searchJobs, &wg)
	wg.Wait()
}

func dirTraversal(path string, opts *SearchOptions, searchJobs chan *searchJob, wg *sync.WaitGroup) {
	info, err := os.Lstat(path)
	if err != nil {
		cobra.CheckErr(fmt.Sprintf("couldn't lstat path %s: %s\n", path, err))
	}

	if !info.IsDir() {
		wg.Add(1)
		searchJobs <- &searchJob{
			path,
			opts,
		}
		return
	}

	f, err := os.Open(path)
	if err != nil {
		cobra.CheckErr(fmt.Sprintf("couldn't open path %s: %s\n", path, err))
	}
	dirNames, err := f.Readdirnames(-1)
	if err != nil {
		cobra.CheckErr(fmt.Sprintf("couldn't read dir names for path %s: %s\n", path, err))
	}

	for _, deeperPath := range dirNames {
		dirTraversal(filepath.Join(path, deeperPath), opts, searchJobs, wg)
	}
}

func searchWorker(jobs chan *searchJob, wg *sync.WaitGroup, basepath string, stdout io.Writer) {
	for job := range jobs {
		f, err := os.Open(job.path)
		if err != nil {
			cobra.CheckErr(fmt.Sprintf("couldn't open path %s: %s\n", job.path, err))
		}

		scanner := bufio.NewScanner(f)
		isBinary := false

		line := 1
		for scanner.Scan() {
			text := scanner.Bytes()

			// Check the first buffer for NUL
			if line == 1 {
				isBinary = bytes.IndexByte(text, 0) != -1
			}

			id, err := filepath.Rel(basepath, job.path)
			id = idStyle.Render(id)
			match := matchStyle.Render(string(text))
			if job.opts.Kind == LITERAL {
				cobra.CheckErr(err)
				if job.opts.Finder.next(string(text)) != -1 {
					if isBinary {
						fmt.Fprintf(stdout, "Binary file %s matches\n", id)
						break
					} else if job.opts.Lines {
						fmt.Fprintf(stdout, "%s:%d %s\n", id, line, match)
					} else {
						fmt.Fprintf(stdout, "%s %s\n", id, match)
					}
				}
			} else if job.opts.Kind == REGEX {
				if job.opts.Regex.Find(scanner.Bytes()) != nil {
					if isBinary {
						fmt.Fprintf(stdout, "Binary file %s matches\n", id)
						break
					} else if job.opts.Lines {
						fmt.Fprintf(stdout, "%s:%d %s\n", id, line, match)
					} else {
						fmt.Fprintf(stdout, "%s %s\n", id, match)
					}
				}
			}
			line++
		}
		wg.Done()
	}
}
