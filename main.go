package main

import (
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"github.com/chzyer/readline"
	"os/exec"
	"fmt"
)

func help(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("login"),
	readline.PcItem("bye"),
	readline.PcItem("help"),
	readline.PcItem("go",
		readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
		readline.PcItem("install",
			readline.PcItem("-v"),
			readline.PcItem("-vv"),
			readline.PcItem("-vvv"),
		),
		readline.PcItem("test"),
	),
	readline.PcItem("sleep"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

const RetContinue = 0
const RetExit = 1

func parseLine(line string) (string, string, string){
	words := []string{"","",""}
	words = strings.Split(line, " ")
	if len(words) == 1 {
		return words[0], "", ""
	} else if len(words) == 2 {
		return words[0], words[1], ""
	} else if len(words) == 3 {
		return words[0],words[1],words[2]
	} else {
		return words[0],words[1],words[2]
	}
}

func execCurl(method string, path string, body string) string {
	url := s.BaseURL + path
	fmt.Printf("curl -X %s %s\n", method, url)
	var cmd *exec.Cmd
	var args []string

	if body != "" {
		args = []string{"-X", method, "-H", "Content-type: application/json", "-d", body}
	} else {
		args = []string{"-X", method}
	}

	for k,v := range s.Headers {
		h := fmt.Sprintf("%s: %s", k, v)
		args = append(args, "-H", h)
	}

	args = append(args, url)
	cmd = exec.Command("curl", args...)
	byts, err := cmd.Output()
	if err != nil {
		fmt.Printf(err.Error())
	}
	out := string(byts)
	fmt.Printf(out)
	return out
}

func cmdGet(path string) {
	execCurl("GET", path,"")
}

func cmdDelete(path string) {
	execCurl("DELETE", path, "")
}

func cmdPost(path string, body string) {
	execCurl("POST", path, body)
}

func cmdPut(path string, body string) {
	execCurl("PUT", path, body)
}

type state struct {
	BaseURL string
	Headers map[string]string
}

var s state

func cmdHeader(key string, value string) {
	if value != "" {
		s.Headers[key] = value
	} else {
		println(s.Headers[key])
	}
}

func cmdStatus() {
	println("base-url: " + s.BaseURL)
	for k,v := range s.Headers {
		fmt.Printf("%s => %s\n", k, v)
	}
}

func processLine(l *readline.Instance, line string) int {
	cmd, arg1, arg2 := parseLine(line)
	log.Printf("%s:%s:%s", cmd,arg1,arg2)
	switch {
	case cmd == "get":
		cmdGet(arg1)
	case cmd == "delete":
		cmdDelete(arg1)
	case cmd == "post":
		cmdPost(arg1, arg2)
	case cmd == "put":
		cmdPut(arg1, arg2)
	case cmd == "base-url":
		if arg1 != "" {
			s.BaseURL = arg1
		} else {
			println(s.BaseURL)
		}
	case cmd == "header":
		cmdHeader(arg1, arg2)
	case cmd == "status":
		cmdStatus()
	case cmd == "help":
		help(l.Stderr())
	case cmd == "bye":
		return RetExit
	case cmd == "":
	default:
		log.Println("command not found:", strconv.Quote(line))
	}

	return RetContinue
}

const Prompt = "\033[31mcurl-shell>\033[0m "
const HistoryFile = "/tmp/curl-shell.history"

func main() {
	s.Headers = make(map[string]string)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          Prompt,
		HistoryFile:     HistoryFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	log.SetOutput(l.Stderr())
	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		ret := processLine(l, line)
		if ret == RetExit {
			return
		}
	}
}
