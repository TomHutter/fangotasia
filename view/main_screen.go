package view

/*
https://github.com/jerilseb/gush
extern void disableRawMode();
extern void enableRawMode();
*/
//import "C"

import (
	"fantasia/config"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

var notice struct {
	Message string
	Color   string
	Sleep   int
}

func AddFlashNotice(message string, sleep int, color string) {
	notice.Message = message
	notice.Color = color
	notice.Sleep = sleep
}

func FlashNotice() bool {
	if len(notice.Message) == 0 {
		return false
	}
	fmt.Printf("\n%s%s%s\n", notice.Color, notice.Message, config.NEUTRAL)
	if notice.Sleep < 0 {
		fmt.Printf("\nWeiter \u23CE\n")
		Scanner("once: true")
	} else {
		time.Sleep(time.Duration(notice.Sleep) * time.Second)
	}
	notice.Message = ""
	notice.Color = ""
	notice.Sleep = 0
	return true
}

func PrintScreen(text []string) {
	// clear screen
	block := strings.Join(text, "\n")
	fmt.Print("\033[H\033[2J")
	fmt.Println(block)
	if FlashNotice() {
		fmt.Print("\033[H\033[2J")
		fmt.Println(block)
	}
}

/*
func Input() {
	verbs := config.Verbs
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("Enter a verb: ").
		SetFieldWidth(30)
	inputField.SetDoneFunc(func(key tcell.Key) {
		fmt.Println(inputField.GetText())
		//app.Stop()
	})
	inputField.SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range verbs {
			if strings.HasPrefix(strings.ToLower(string(word)), strings.ToLower(currentText)) {
				entries = append(entries, string(word))
			}
		}
		if len(entries) <= 1 {
			entries = nil
		}
		return
	})
	if err := app.SetRoot(inputField, true).Run(); err != nil {
		panic(err)
	}
}
*/

func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		os.Exit(0)
	}()
}

func Scanner(params ...string) (line string) {
	var b []byte = make([]byte, 4)
	var once bool
	var prompt string

	if len(params) > 0 {
		for _, v := range params {
			val := strings.Split(v, ": ")
			switch val[0] {
			case "once":
				once = strings.ToLower(val[1]) == "true"
			case "prompt":
				prompt = val[1]
			}
		}
	}
	if len(prompt) > 0 {
		fmt.Printf("\033[36m%s\033[m", prompt)
	}
	for {
		os.Stdin.Read(b)
		r, _ := utf8.DecodeRune(b)
		// once set to true => return directly after one keypress
		if once {
			line = string(r)
			return
		}
		// the enter key was pressed
		if b[0] == 10 {
			//fmt.Println(line)
			line = strings.TrimSpace(line)
			return
		}

		// Special control key was pressed
		if b[0] == 27 {
			continue
		}

		// backspace was pressed
		if b[0] == 127 {
			if len(line) > 0 {
				fmt.Print("\b\033[K")
				_, lastRuneSize := utf8.DecodeLastRuneInString(line)
				line = line[:len(line)-lastRuneSize]
			}
			continue
		}

		// Any normal character
		fmt.Printf("%s", string(r))
		line += string(r)
	}
}

/*
func exit() {
	C.disableRawMode()
	os.Exit(0)
}
*/
