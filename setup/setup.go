package setup

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
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
	Verbs        []Verb                   `yaml:"verbs"`
	Nouns        []string                 `yaml:"nouns"`
	Objects      map[int]ObjectProperties `yaml:"objects"`
	ID           int
	Reactions    map[string]Reaction    `yaml:"reactions"`
	Locations    map[int]AreaProperties `yaml:"locations"`
	Overwrites   []MapOverwrites        `yaml:"overwrites"`
	Contitions   map[string]Condition   `yaml:"conditions"`
	TextElements map[string]string      `yaml:"elements"`
}

// Long and short description and the article for the noun
type description struct {
	Long    string
	Short   string
	Alt     string
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
	Statement []string
	OK        bool
	KO        bool
	Color     string
}

type Condition map[string]string

var PathName string
var TextElements map[string]string
var GameObjects map[int]ObjectProperties
var GameAreas map[int]AreaProperties
var Overwrites []MapOverwrites
var Reactions map[string]Reaction
var Conditions map[string]Condition
var Verbs []Verb
var Moves int

var Beads int
var BoxLen int
var Flags map[string]bool

var Map [12][10]int

func initMap() {
	for y := 0; y < 12; y++ {
		for x := 0; x < 10; x++ {
			Map[y][x] = 0
		}
	}
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

func getMapOverwrites() (overwrites []MapOverwrites) {
	var c conf
	c.getConf(PathName + "/config/map_overwrites.yaml")
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

func removeMapVerb(verbs []Verb) []Verb {
	for i, v := range verbs {
		if v.Func == "Map" {
			return append(verbs[:i], verbs[i+1:]...)
		}
	}
	return verbs
}

func AddMapVerb(verbs []Verb) []Verb {
	var c conf
	c.getConf(PathName + "/config/verbs.yaml")
	for _, v := range c.Verbs {
		if v.Func == "Map" {
			return append(verbs, v)
		}
	}
	return verbs
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

func Setup() {
	var c conf
	c.getConf(PathName + "/config/text_elements.yaml")
	TextElements = c.TextElements
	c.getConf(PathName + "/config/objects.yaml")
	GameObjects = c.Objects
	c.getConf(PathName + "/config/locations.yaml")
	GameAreas = c.Locations
	c.getConf(PathName + "/config/reactions.yaml")
	Reactions = c.Reactions
	c.getConf(PathName + "/config/verbs.yaml")
	Verbs = removeMapVerb(c.Verbs)
	c.getConf(PathName + "/config/conditions.yaml")
	Conditions = c.Contitions
	Overwrites = getMapOverwrites()
	initBoxLen()
	initMap()
	Flags = make(map[string]bool, 6)
	Flags["DoorOpen"] = false
	Flags["BoxOpen"] = false
	Flags["MapMissed"] = false
	Flags["Moore"] = false
	Flags["Castle"] = false
	Flags["Tree"] = false
}

func ObjectsInArea(area Area) (objects []Object) {
	for id, prop := range GameObjects {
		if prop.Area == area.ID {
			objects = append(objects, Object{id, prop})
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

func GetReactionByName(name string) (r Reaction) {
	r = Reactions[name]
	r.Statement = make([]string, len(Reactions[name].Statement))
	copy(r.Statement, Reactions[name].Statement)
	return
}
