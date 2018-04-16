package main

import (
	"io"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/chzyer/readline"
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
	readline.PcItem("mode",
		readline.PcItem("vi"),
		readline.PcItem("emacs"),
	),
	readline.PcItem("login"),
	readline.PcItem("say",
		readline.PcItemDynamic(listFiles("./"),
			readline.PcItem("with",
				readline.PcItem("following"),
				readline.PcItem("items"),
			),
		),
		readline.PcItem("hello"),
		readline.PcItem("bye"),
	),
	readline.PcItem("setprompt"),
	readline.PcItem("setpassword"),
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

func processLine(l *readline.Instance, line string) int {
	switch {
	case line == "login":
		pswd, err := l.ReadPassword("please enter your password: ")
		if err != nil {
			break
		}
		println("you enter:", strconv.Quote(string(pswd)))
	case line == "help":
		help(l.Stderr())
	case strings.HasPrefix(line, "setprompt"):
		if len(line) <= 10 {
			log.Println("setprompt <prompt>")
			break
		}
		l.SetPrompt(line[10:])
	case line == "bye":
		return RetExit
	case line == "sleep":
		log.Println("sleep 4 second")
		time.Sleep(4 * time.Second)
	case line == "":
	default:
		log.Println("command not found:", strconv.Quote(line))
	}

	return RetContinue
}

const Prompt = "\033[31mcurl-shell>\033[0m "
const HistoryFile = "/tmp/curl-shell.history"

func main() {
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
