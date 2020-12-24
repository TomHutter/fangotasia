package setup

import (
	"io/ioutil"
	"log"
	"os"

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
	Func   string
	Single bool
	Name   map[string][]string
}

// Conf : Struct to read from yaml config files
type conf struct {
	Verbs        []Verb                       `yaml:"verbs"`
	Objects      map[int]ObjectProperties     `yaml:"objects"`
	Reactions    map[string]Reaction          `yaml:"reactions"`
	Locations    map[int]AreaProperties       `yaml:"locations"`
	Overwrites   []MapOverwrites              `yaml:"overwrites"`
	Contitions   map[string]Condition         `yaml:"conditions"`
	TextElements map[string]map[string]string `yaml:"elements"`
}

// Long, short and alternative description and the article for the noun
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
	Description map[string]description
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
	Description map[string]description
	Directions  [4]int
	Coordinates Coordinates
}

type Coordinates struct {
	Y int
	X int
}

type MapOverwrites struct {
	Area    int
	Content map[string][]string
}

type Reaction struct {
	Statement map[string][]string
	OK        bool
	KO        bool
	Color     string
}

type Condition map[string]map[string]string

var (
	PathName     string
	Language     string
	TextElements map[string]map[string]string
	GameObjects  map[int]ObjectProperties
	GameAreas    map[int]AreaProperties
	Overwrites   []MapOverwrites
	Reactions    map[string]Reaction
	Conditions   map[string]Condition
	Verbs        []Verb
	Moves        int

	Beads  int
	BoxLen int = 19
	Flags  map[string]bool

	Map [12][10]int
)

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
	yamlFile, err := ioutil.ReadFile(PathName + "/config/" + filename)
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
	c.getConf("map_overwrites.yaml")
	for _, v := range c.Overwrites {
		var o MapOverwrites
		o.Area = v.Area
		overwrites = append(overwrites, o)
		overwrites[len(overwrites)-1].Content = make(map[string][]string)
		o.Content = make(map[string][]string)
		for lang, c := range v.Content {
			var ov []string
			ov = make([]string, 3)
			//o.Content[lang] = [3]string{}
			o.Content[lang] = make([]string, 3)
			overwrites[len(overwrites)-1].Content[lang] = make([]string, 3)
			for i, line := range c {
				ov[i] = line
			}
			copy(o.Content[lang], ov)
			copy(overwrites[len(overwrites)-1].Content[lang], ov)
		}
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
	c.getConf("verbs.yaml")
	for _, v := range c.Verbs {
		if v.Func == "Map" {
			return append(verbs, v)
		}
	}
	return verbs
}

/*
// InitBoxLen : Get min length for boxes to fit all short descriptions of locations
func initBoxLen() {
	BoxLen = 0
	for _, v := range GameAreas {
		lineLen := len([]rune(strings.Split(v.Description[Language].Short, "\n")[0]))
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
*/

func setLang() {
	if _, err := os.Stat(PathName + "/save/fangotasia.lang"); os.IsNotExist(err) {
		Language = "en"
	} else {
		lang, _ := ioutil.ReadFile(PathName + "/save/fangotasia.lang")
		Language = string(lang)
	}
}

func Setup() {
	var c conf
	setLang()
	c.getConf("text_elements.yaml")
	TextElements = c.TextElements
	c.getConf("objects.yaml")
	GameObjects = c.Objects
	c.getConf("locations.yaml")
	GameAreas = c.Locations
	c.getConf("reactions.yaml")
	Reactions = c.Reactions
	c.getConf("verbs.yaml")
	Verbs = removeMapVerb(c.Verbs)
	c.getConf("conditions.yaml")
	Conditions = c.Contitions
	Overwrites = getMapOverwrites()
	//initBoxLen()
	initMap()
	Flags = make(map[string]bool, 7)
	Flags["DoorOpen"] = false
	Flags["BoxOpen"] = false
	Flags["MapMissed"] = false
	Flags["HoodVanished"] = false
	Flags["Swamp"] = false
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

func GetOverwriteByArea(area int) (ok bool, o MapOverwrites) {
	for _, o := range Overwrites {
		if o.Area == area {
			return true, o
		}
	}
	return false, o
}

func GetReactionByName(name string) (r Reaction) {
	r = Reactions[name]
	r.Statement = make(map[string][]string, len(Reactions[name].Statement))
	for lang := range Reactions[name].Statement {
		r.Statement[lang] = make([]string, len(Reactions[name].Statement[lang]))
		copy(r.Statement[lang], Reactions[name].Statement[lang])
	}
	return
}
