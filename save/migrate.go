package main

import (
	"fangotasia/setup"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	GameObjects    map[int]ObjectProperties
	NewGameObjects map[int]NewObjectProperties
	AreaID         int
)

// ObjectProperties : Contain long and short description of locations.
type ObjectProperties struct {
	Description description
	Area        int
	Value       int
}

// ObjectProperties : Contain long and short description of locations.
type NewObjectProperties struct {
	Description map[string]description
	Area        int
	Value       int
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
	Properties NewObjectProperties
}

// init() will be called before main() by go convention
func init() {
	pathName, _ := os.Getwd()
	setup.PathName = pathName + "/../"
}

func main() {
	setup.Setup()
	Load()
	Migrate()
	Save()
}

func Migrate() {
	NewGameObjects = make(map[int]NewObjectProperties, len(GameObjects))

	for i, o := range GameObjects {
		var obj Object
		obj.ID = i
		obj.Properties.Area = o.Area
		obj.Properties.Value = o.Value
		obj.Properties.Description = make(map[string]description)
		gameObject := setup.GetObjectByID(i)
		for lang := range gameObject.Properties.Description {
			obj.Properties.Description[lang] = description(gameObject.Properties.Description[lang])
		}
		// object has not default long description
		if o.Description.Long == "[#ff69b4::b]<IMKE>[blue:black:-] den pink Diamanten" {
			for lang := range gameObject.Properties.Description {
				var desc description
				switch lang {
				case "de":
					desc.Long = "[#ff69b4::b]<IMKE>[blue:black:-] den pink Diamanten"
				case "en":
					desc.Long = "[#ff69b4::b]<IMKE>[blue:black:-] the pink diamond"

				}
				desc.Short = gameObject.Properties.Description[lang].Short
				desc.Alt = gameObject.Properties.Description[lang].Alt
				desc.Article = gameObject.Properties.Description[lang].Article
				obj.Properties.Description[lang] = desc
				continue
			}
		}
		if strings.Contains(o.Description.Long, "besudelt") {
			for lang := range gameObject.Properties.Description {
				var desc description
				switch lang {
				case "de":
					desc.Long = o.Description.Long
				case "en":
					article := gameObject.Properties.Description[lang].Article
					parts := strings.Split(gameObject.Properties.Description[lang].Long, " ")[1:]
					long := strings.Join(parts, " ")
					desc.Long = fmt.Sprintf(setup.Conditions["fango"][lang][article], long)
				}
				desc.Short = gameObject.Properties.Description[lang].Short
				desc.Alt = gameObject.Properties.Description[lang].Alt
				desc.Article = gameObject.Properties.Description[lang].Article
				obj.Properties.Description[lang] = desc
				continue
			}
		}
		if o.Description.Long != gameObject.Properties.Description["de"].Long {
		label:
			for _, condObj := range setup.Conditions {
				for _, cond := range condObj {
					if cond["de"] == o.Description.Long {
						for lang := range cond {
							var desc description
							desc.Long = cond[lang]
							desc.Short = gameObject.Properties.Description[lang].Short
							desc.Alt = gameObject.Properties.Description[lang].Alt
							desc.Article = gameObject.Properties.Description[lang].Article
							obj.Properties.Description[lang] = desc
						}
						break label
					}
				}
			}
		}
		NewGameObjects[i] = obj.Properties
	}

}

func Save() {
	m := make(map[interface{}]interface{})
	m["area"] = AreaID
	m["map"] = setup.Map
	m["objects"] = NewGameObjects
	m["flags"] = setup.Flags
	m["moves"] = setup.Moves

	os.Rename("./fangotasia.sav", "./fangotasia.sav1")
	file, err := os.Create("./fangotasia.sav")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(m); err != nil {
		fmt.Println(err.Error())
		return
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return
}

func Load() {
	var content struct {
		AreaID  int                      `yaml:"area"`
		AreaMap [12][10]int              `yaml:"map"`
		Objects map[int]ObjectProperties `yaml:"objects"`
		Flags   map[string]bool          `yaml:"flags"`
		Moves   int                      `yaml:"moves"`
	}

	yamlFile, err := ioutil.ReadFile("./fangotasia.sav")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	GameObjects = content.Objects
	setup.Map = content.AreaMap

	setup.Flags = content.Flags
	setup.Moves = content.Moves

	AreaID = content.AreaID
}
