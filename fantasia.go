package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"unicode/utf8"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Verbs   []string `yaml:"verbs"`
	Nouns   []string `yaml:"nouns"`
	Objects []string `yaml:"objects"`
	Answers []string `yaml:"answers"`
}

func (c *conf) getConf(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func print_screen(text []string) {
	// clear screen
	fmt.Print("\033[H\033[2J")
	for _, t := range text {
		fmt.Println(t)
	}
}

func scanner() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	var b []byte = make([]byte, 4)
	for {
		os.Stdin.Read(b)
		r, _ := utf8.DecodeRune(b)
		fmt.Println(string(r))
	}
}

func prelude() {
	var text []string
	text = make([]string, 4)
	text = append(text, " fantasia ")
	text = append(text, " - Ein Adventure von Klaus Hartmuth -")
	text = append(text, " - Überarbeitet von Tom Hutter -")
	print_screen(text)
	scanner()
}

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

func main() {
	var c conf
	c.getConf("verbs.yaml")
	verbs := c.Verbs
	c.getConf("nouns.yaml")
	nouns := c.Nouns
	c.getConf("objects.yaml")
	objects := c.Objects

	// Setup our Ctrl+C handler
	setupCloseHandler()
	prelude()

	fmt.Println(verbs)
	fmt.Println(nouns)
	fmt.Println(objects[0])
}
