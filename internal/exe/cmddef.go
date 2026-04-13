package exe

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type CmdDef struct {
	Name string
	Args []string
}

func NewCmdDef(name string, args ...string) CmdDef {
	return CmdDef{Name: name, Args: args}
}

func (c CmdDef) String() string {
	return strings.TrimSpace(c.Name + " " + strings.Join(c.Args, " "))
}

func (c CmdDef) Run() error {
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func DetermineCmds(path string) []CmdDef {
	res := make([]CmdDef, 0)

	// Infer from content first
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		first := ""
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			first = scanner.Text()
			break
		}
		if strings.HasPrefix(first, "#!") {
			if runtime.GOOS != "windows" {
				if s, _ := os.Stat(path); isExecAny(s.Mode()) {
					res = append(res, NewCmdDef(path))
				}
			}
			// Extract from shebang
			if strings.HasPrefix(first, "#!/usr/bin/env ") {
				rem := strings.Split(first[15:], " ")
				if len(rem) > 0 {
					rem = append(rem, path)
					res = append(res, NewCmdDef(rem[0], rem[1:]...))
				}
			} else if first == "#!/bin/sh" || strings.HasPrefix(first, "#!/bin/sh ") {
				res = append(res, NewCmdDef("sh", path))
			}
		}
	}

	// TODO: config file with preference or filters, e.g. ignored_runners = ["kscript"]
	switch filepath.Ext(path) {
	case ".sh":
		res = append(res,
			NewCmdDef("sh", path),
		)
	case ".kt", ".kts":
		res = append(res,
			NewCmdDef("kotlinc", "-script", path),
			NewCmdDef("kscript", path),
		)
	case ".ts":
		res = append(res,
			NewCmdDef("deno", path),
			NewCmdDef("bun", path),
		)
	case ".js", ".cjs", ".mjs":
		res = append(res,
			NewCmdDef("node", path),
			NewCmdDef("deno", path),
			NewCmdDef("bun", path),
		)
	case ".py":
		res = append(res, NewCmdDef("python", path))
	case ".go":
		res = append(res, NewCmdDef("go", "run", path))
	case ".lua":
		res = append(res, NewCmdDef("lua", path))
	}

	return distinct(res)
}

func distinct(defs []CmdDef) []CmdDef {
	if len(defs) == 0 {
		return defs
	}
	lut := make(map[string]bool)
	res := make([]CmdDef, 0)
	for _, def := range defs {
		key := def.String()
		if _, seen := lut[key]; !seen {
			res = append(res, def)
			lut[key] = true
		}
	}
	return res
}

// Source - https://stackoverflow.com/a/60128480
// Posted by icza, modified by community. See post 'Timeline' for change history
// Retrieved 2026-04-13, License - CC BY-SA 4.0
func isExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}
