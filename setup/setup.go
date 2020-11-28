package setup

import (
	"io/ioutil"
	"log"
	"path"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

// Make life more colorful :-)
const (
	// https://en.wikipedia.org/wiki/ANSI_escape_code
	RED     = "\033[01;31m"
	GREEN   = "\033[01;32m"
	YELLOW  = "\033[01;33m"
	CYAN    = "\033[01;96m"
	BLUE    = "\033[01;34m"
	WHITE   = "\033[01;97m"
	NEUTRAL = "\033[0m"
)

// Runes for creating boxes
const (
	BTL = "\u250F"
	BTR = "\u2513"
	BBL = "\u2517"
	BBR = "\u251B"
	HL  = "\u2501"
	VL  = "\u2503"
	LC  = "\u252B"
	RC  = "\u2523"
	TC  = "\u253B"
	BC  = "\u2533"
	AR  = "\u2BC8"
	AL  = "\u2BC7"
	AU  = "\u2BC5"
	AD  = "\u2BC6"
)

const (
	INVENTORY = 1000
	INUSE     = 2000
)

type Verb struct {
	Name   string
	Func   string
	Single bool
}

// Conf : Struct to read from yaml config files
type conf struct {
	Verbs      []Verb                   `yaml:"verbs"`
	Nouns      []string                 `yaml:"nouns"`
	Objects    map[int]ObjectProperties `yaml:"objects"`
	ID         int
	Reactions  map[string]Reaction    `yaml:"reactions"`
	Locations  map[int]AreaProperties `yaml:"locations"`
	Overwrites []MapOverwrites        `yaml:"overwrites"`
}

// Long and short description and the article for the noun
type description struct {
	Long    string
	Short   string
	Article string
}

type Object struct {
	ID         int
	Properties ObjectProperties
}

// ObjectProperties : Contain long and short description of locations.
type ObjectProperties struct {
	Description description
	Area        int
	Value       int
}

type Area struct {
	ID         int
	Properties AreaProperties
}

// Area : long and short description of locations.
//        directions: which area will be reachable in n,s,e,w
//        coordinates: y and x coordinates for area on map
type AreaProperties struct {
	Description description
	Directions  [4]int
	Coordinates Coordinates
	Visited     bool
}

type Coordinates struct {
	Y int
	X int
}

type MapOverwrites struct {
	Area    int
	Content [3]string
}

type Reaction struct {
	Statement string
	OK        bool
	KO        bool
	Color     string
	Sleep     int
}

var GameObjects map[int]ObjectProperties
var GameAreas map[int]AreaProperties
var Overwrites []MapOverwrites
var Reactions map[string]Reaction
var Verbs []Verb

var BoxLen int
var Flags map[string]bool

var Map [12][10]int

func initMap() {
	Map[11][0] = 1
}

func AreaVisible(a int) bool {
	area := GetAreaByID(a)
	return Map[area.Properties.Coordinates.Y][area.Properties.Coordinates.X] != 0
}

// GetConf : Read yaml config files into struct Conf
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

func getMapOverwrites(pathname string) (overwrites []MapOverwrites) {
	var c conf
	c.getConf(pathname + "/../config/map_overwrites.yaml")
	for _, v := range c.Overwrites {
		var o MapOverwrites
		o.Area = v.Area
		for i, line := range v.Content {
			o.Content[i] = line
		}
		overwrites = append(overwrites, o)
	}
	return
}

// InitBoxLen : Get min length for boxes to fit all short descriptions of locations
func initBoxLen() {
	BoxLen = 0
	for _, v := range GameAreas {
		lineLen := len([]rune(strings.Split(v.Description.Short, "\n")[0]))
		if lineLen > BoxLen {
			BoxLen = lineLen
		}
	}
	// make boxLen odd for middle connection piece
	if BoxLen%2 == 0 {
		BoxLen++
	}
	BoxLen = BoxLen + 2 // one blank and border left and right
}

func Init() {
	var c conf
	_, filename, _, _ := runtime.Caller(0)
	pathname := path.Dir(filename)
	c.getConf(pathname + "/../config/objects.yaml")
	GameObjects = c.Objects
	c.getConf(pathname + "/../config/locations.yaml")
	GameAreas = c.Locations
	c.getConf(pathname + "/../config/reactions.yaml")
	Reactions = c.Reactions
	c.getConf(pathname + "/../config/verbs.yaml")
	Verbs = c.Verbs
	Overwrites = getMapOverwrites(pathname)
	initBoxLen()
	initMap()
	Flags = make(map[string]bool, 3)
	Flags["DoorOpen"] = false
	Flags["BoxOpen"] = false
	Flags["MapMissed"] = false
}

func ObjectsInArea(area Area) (objects []Object) {
	for id, prop := range GameObjects {
		if prop.Area == area.ID {
			objects = append(objects, Object{id, prop})
		}
	}
	return
}

func GetObjectByName(name string) (object Object) {
	for id, prop := range GameObjects {
		if strings.ToLower(prop.Description.Short) == strings.ToLower(name) {
			return Object{id, prop}
		}
	}
	return
}

func GetObjectByID(id int) (object Object) {
	return Object{id, GameObjects[id]}
}

func GetAreaByID(id int) (area Area) {
	return Area{id, GameAreas[id]}
}

func GetOverwriteByArea(area int) (o MapOverwrites) {
	for _, o := range Overwrites {
		if o.Area == area {
			return o
		}
	}
	return
}
