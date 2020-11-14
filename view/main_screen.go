package view

import (
	"fantasia/config"
	"fmt"
	"strings"
	"time"
)

func AppendText(block *[]string, newText string, color ...string) {
	text := *block
	if color == nil {
		*block = append(text, newText)
	}
	*block = append(text, fmt.Sprintf("%s%s%s", color[0], newText, config.WHITE))
}

func Flash(text []string, err string) {
	flashText := make([]string, len(text))
	copy(flashText, text)
	flashText = append(text, "")
	flashText = append(text, fmt.Sprintf("%s%s%s", config.RED, err, config.NEUTRAL))
	PrintScreen(flashText)
	time.Sleep(2 * time.Second)
	PrintScreen(text)
}

func PrintScreen(text []string) {
	// clear screen
	fmt.Print("\033[H\033[2J")
	block := strings.Join(text, "\n")
	fmt.Println(block)
}
