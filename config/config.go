package config

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

var GameObjects []Object
var Locations []description
var Overwrites [][3]string
var Answers []string

type Verb struct {
	Name   string
	Func   string
	Single bool
}

var Verbs []Verb
var BoxLen int

// Conf : Struct to read from yaml config files
type conf struct {
	Verbs      []Verb          `yaml:"verbs"`
	Nouns      []string        `yaml:"nouns"`
	Objects    []Object        `yaml:"objects"`
	Answers    []string        `yaml:"answers"`
	Locations  []description   `yaml:"locations"`
	Overwrites []mapOverwrites `yaml:"overwrites"`
}

// Places : Contains long and short description of locations.
type description struct {
	Long    string
	Short   string
	Article string
}

// ObjectProperties : Contains long and short description of locations.
type Object struct {
	ID          int
	Description description
	Area        int
	Value       int
}

type mapOverwrites struct {
	Area    int
	Content [3]string
}

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

// Areas start at 1 in the original game.
// Kept it like this to adapt game more easily.
var Areas = [52][4]int{
	{}, {0, 0, 2, 0}, {8, 0, 3, 1}, {9, 0, 0, 2}, {10, 0, 5, 0},
	{11, 0, 0, 4}, {13, 0, 7, 0}, {0, 8, 0, 0}, {7, 2, 9, 0},
	{15, 3, 10, 8}, {0, 4, 11, 9}, {0, 5, 12, 10}, {0, 0, 0, 11},
	{16, 6, 0, 0}, {17, 0, 15, 0}, {0, 9, 0, 14}, {0, 13, 17, 0},
	{0, 14, 18, 16}, {24, 0, 19, 17}, {25, 0, 20, 18}, {0, 0, 21, 19},
	{26, 0, 0, 20}, {27, 0, 23, 0}, {0, 0, 24, 22}, {29, 18, 0, 23},
	{30, 19, 0, 0}, {32, 21, 0, 0}, {0, 22, 28, 0}, {0, 0, 0, 27},
	{0, 0, 30, 0}, {0, 25, 0, 0}, {0, 0, 0, 0}, {0, 26, 0, 33},
	{0, 34, 32, 0}, {33, 36, 37, 35}, {0, 0, 34, 0}, {34, 0, 5, 0},
	{37, 38, 39, 34}, {37, 38, 38, 38}, {40, 39, 39, 37}, {42, 39, 0, 0},
	{43, 41, 42, 41}, {44, 40, 42, 41}, {47, 41, 43, 43}, {48, 42, 45, 44},
	{49, 45, 45, 44}, {50, 51, 0, 0}, {47, 43, 48, 47}, {48, 44, 48, 47},
	{0, 45, 50, 0}, {0, 46, 0, 49}, {46, 0, 0, 0},
}

/*
var AreaMap = [12][10]int{
	{0, 0, 0, 0, 0, 0, 47, 48, 49, 50},
	{0, 0, 0, 0, 0, 0, 43, 44, 45, 46},
	{0, 0, 0, 0, 0, 0, 41, 42, 0, 51},
	{0, 0, 0, 0, 33, 32, 0, 40, 0, 0},
	{0, 0, 0, 35, 34, 0, 37, 39, 0, 0},
	{27, 28, 29, 30, 36, 0, 38, 0, 0, 0},
	{22, 23, 24, 25, 0, 26, 0, 0, 0, 0},
	{16, 17, 18, 19, 20, 21, 0, 0, 0},
	{13, 14, 15, 0, 0, 0, 0, 0, 0, 0},
	{6, 7, 0, 0, 0, 0, 0, 0, 0, 0},
	{0, 8, 9, 10, 11, 12, 0, 0, 0, 0},
	{1, 2, 3, 4, 5, 0, 0, 0, 0, 0},
}
*/

// AreaCoordinates  Array with the coordinates [x,y] for each area
// in a field [0..11][0..9]
// Again: First area starts with 1
var AreaCoordinates = [52]Coordinates{
	{}, {11, 0}, {11, 1}, {11, 2}, {11, 3}, {11, 4},
	{9, 0}, {9, 1},
	{10, 1}, {10, 2}, {10, 3}, {10, 4}, {10, 5},
	{8, 0}, {8, 1}, {8, 2},
	{7, 0}, {7, 1}, {7, 2}, {7, 3}, {7, 4}, {7, 5},
	{6, 0}, {6, 1}, {6, 2}, {6, 3}, {6, 5},
	{5, 0}, {5, 1}, {5, 2}, {5, 3}, {9, 3}, {3, 5}, {3, 4},
	{4, 4}, {4, 3}, {5, 4}, {4, 6}, {5, 6}, {4, 7}, {3, 7},
	{2, 6}, {2, 7}, {1, 6}, {1, 7}, {1, 8}, {1, 9},
	{0, 6}, {0, 7}, {0, 8}, {0, 9}, {2, 9},
}

type Coordinates struct {
	Y int
	X int
	//Visible bool
}

/*
var MapOverwrite = [52]MapSpecials{
	{},{},{},{},{},
	{},{},{{"blah", "fahsel", "Blubb"},{},
	{11, 0, 0, 4}, {13, 0, 7, 0}, {0, 8, 0, 0}, {7, 2, 9, 0},
	{15, 3, 10, 8}, {0, 4, 11, 9}, {0, 5, 12, 10}, {0, 0, 0, 11},
	{16, 6, 0, 0}, {17, 0, 15, 0}, {0, 9, 0, 14}, {0, 13, 17, 0},
	{0, 14, 18, 16}, {24, 0, 19, 17}, {25, 0, 20, 18}, {0, 0, 21, 19},
	{26, 0, 0, 20}, {27, 0, 23, 0}, {0, 0, 24, 22}, {29, 18, 0, 23},
	{30, 19, 0, 0}, {32, 21, 0, 0}, {0, 22, 28, 0}, {0, 0, 0, 27},
	{0, 0, 30, 0}, {0, 25, 0, 0}, {0, 0, 0, 0}, {0, 26, 0, 33},
	{0, 34, 32, 0}, {33, 36, 37, 35}, {0, 0, 34, 0}, {34, 0, 5, 0},
	{37, 38, 39, 34}, {37, 38, 38, 38}, {40, 39, 39, 37}, {42, 39, 0, 0},
	{43, 41, 42, 41}, {44, 40, 42, 41}, {47, 41, 43, 43}, {48, 42, 45, 44},
	{49, 45, 45, 44}, {50, 51, 0, 0}, {47, 43, 48, 47}, {48, 44, 48, 47},
	{0, 45, 50, 0}, {0, 46, 0, 49}, {46, 0, 0, 0},
}

MapOverwrite[7] = {
*/

// ObjectsInArea contains the coordinates for each object on the map.
/*
var ObjectsInArea = [45][2]int{
	{-1, 0}, {}, {}, {}, {}, {}, {}, {}, {},
	{28, 0}, {29, 0}, {8, 0}, {24, 0}, {2, 0}, {26, 0}, {20, 10}, {19, 0}, {19, 22},
	{18, 0}, {18, 20}, {16, 0}, {13, 0}, {14, 0}, {0, 0}, {0, 26}, {15, 5}, {8, 5},
	{9, 0}, {11, 0}, {12, 0}, {1, 0}, {1, 0}, {1, 0}, {3, 7}, {4, 0}, {4, 0},
	{4, 0}, {0, 0}, {0, 0}, {0, 18}, {30, 0}, {33, 10}, {51, 0}, {40, 0}, {51, 47},
}
*/

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

func getMapOverwrites() (overwrites [][3]string) {
	var c conf
	c.getConf("config/map_overwrites.yaml")
	for _, v := range c.Overwrites {
		// overwrites is already large enough to address v.Area
		if v.Area < len(overwrites) {
			var dummy [3]string
			for i, line := range v.Content {
				dummy[i] = line
			}
			overwrites[v.Area] = dummy
		}
		// overwrites needs expansion to address v.Area
		if v.Area > len(overwrites) {
			var dummy = make([][3]string, v.Area)
			copy(dummy, overwrites)
			overwrites = dummy
		}
		if v.Area == len(overwrites) {
			var dummy [3]string
			for i, line := range v.Content {
				dummy[i] = line
			}
			overwrites = append(overwrites, dummy)
		}
	}
	return
}

// InitBoxLen : Get min length for boxes to fit all short descriptions of locations
func initBoxLen() {
	BoxLen = 0
	for _, v := range Locations {
		lineLen := len([]rune(strings.Split(v.Short, "\n")[0]))
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
	c.getConf("config/objects.yaml")
	GameObjects = c.Objects
	c.getConf("config/locations.yaml")
	Locations = c.Locations
	c.getConf("config/answers.yaml")
	Answers = c.Answers
	c.getConf("config/verbs.yaml")
	Verbs = c.Verbs
	Overwrites = getMapOverwrites()
	initBoxLen()
}

func ObjectsInArea(area int) (objects []Object) {
	for _, o := range GameObjects {
		if o.Area == area {
			objects = append(objects, o)
		}
	}
	return
}

func GetObjectByName(name string) (object Object) {
	for _, o := range GameObjects {
		if strings.ToLower(o.Description.Short) == strings.ToLower(name) {
			return o
		}
	}
	return
}

func GetObjectByID(id int) (object Object) {
	for _, o := range GameObjects {
		if o.ID == id {
			return o
		}
	}
	return
}
