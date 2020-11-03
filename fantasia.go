package main

import (
	"fmt"
	"io/ioutil"
	"log"

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

func main() {
	var c conf
	c.getConf("verbs.yaml")
	verbs := c.Verbs
	c.getConf("nouns.yaml")
	nouns := c.Nouns
	c.getConf("objects.yaml")
	objects := c.Objects

	fmt.Println(verbs)
	fmt.Println(nouns)
	fmt.Println(objects[0])
}
