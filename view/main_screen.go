package view

/*
https://github.com/jerilseb/gush
extern void disableRawMode();
extern void enableRawMode();
*/
//import "C"

import (
	"fangotasia/setup"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"
)

var Notice struct {
	Message string
	Color   string
	Sleep   int
}

func Highlight(s string, c string) (h string) {
	re := regexp.MustCompile("::(.+)::")

	h = string(re.ReplaceAll([]byte(s), []byte("[white::b]"+"$1"+c)))
	return
}

func AddFlashNotice(message string, sleep int, color string) {
	Notice.Message = message
	Notice.Color = color
	Notice.Sleep = sleep
}

func FlashNotice() bool {
	if len(Notice.Message) == 0 {
		return false
	}
	fmt.Printf("\n%s%s%s\n", Notice.Color, Notice.Message, "[-:black:-]")
	if Notice.Sleep < 0 {
		fmt.Printf("\nWeiter \u23CE\n")
		Scanner("once: true")
	} else {
		time.Sleep(time.Duration(Notice.Sleep) * time.Second)
	}
	Notice.Message = ""
	Notice.Color = ""
	Notice.Sleep = 0
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
	verbs := setup.Verbs
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

func Surroundings(area setup.Area) (text []string) {
	desc := strings.Split(area.Properties.Description.Long, "\\n")
	desc0, desc := desc[0], desc[1:]
	text = append(text, fmt.Sprintf("%sIch bin %s", "[yellow]", desc0))
	for _, v := range desc {
		if strings.Contains(v, "++") {
			v = strings.ReplaceAll(v, "++", "")
			text = append(text, fmt.Sprintf("%s%s", "[cyan]", v))
		} else {
			text = append(text, fmt.Sprintf("%s%s", "[yellow]", v))
		}
	}
	text = append(text, "[-:-:-]")
	var items []string
	for _, object := range setup.ObjectsInArea(area) {
		item := Highlight(object.Properties.Description.Long, "[blue:black:-]")
		items = append(items, fmt.Sprintf("%s  - %s", "[blue:black]", item))
	}
	if len(items) > 0 {
		text = append(text, fmt.Sprintf("%sIch sehe:", "[blue:black]"))
		for _, item := range items {
			text = append(text, item)
		}
		text = append(text, "[-:-:-]")
	}
	var directions []string
	for d := 0; d < 4; d++ {
		if area.Properties.Directions[d] != 0 {
			switch d {
			case 0: // N
				directions = append(directions, "Norden")
			case 1: // S
				directions = append(directions, "Süden")
			case 2: // O
				directions = append(directions, "Osten")
			case 3: // W
				directions = append(directions, "Westen")
			}
		}
	}
	text = append(text, fmt.Sprintf("%sRichtungen: %s", "[white:black:b]", strings.Join(directions, ", ")))
	return
}
