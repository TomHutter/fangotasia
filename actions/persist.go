package actions

import (
	"fangotasia/setup"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func folderListing() (filename string) {
	filename = setup.PathName + "/save/fangotasia.sav"
	return
}

func (obj Object) Save(area setup.Area) (r setup.Reaction) {
	m := make(map[interface{}]interface{})
	m["area"] = area.ID
	m["map"] = setup.Map
	m["objects"] = setup.GameObjects
	m["flags"] = setup.Flags
	m["moves"] = setup.Moves

	filename := folderListing()
	if _, err := os.Stat(setup.PathName + "/save/"); os.IsNotExist(err) {
		os.Mkdir(setup.PathName+"/save/", os.FileMode(0755))
	}

	r = setup.Reactions["saved"]
	file, err := os.Create(filename)
	if err != nil {
		r.Statement[0] = err.Error()
		r.OK = false
		return
	}

	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(m); err != nil {
		r.Statement[0] = err.Error()
		r.OK = false
		return
	}

	err = file.Close()
	if err != nil {
		r.Statement[0] = err.Error()
		r.OK = false
		return
	}
	return
}

func (obj Object) Load(area setup.Area) (r setup.Reaction, areaID int) {
	var content struct {
		AreaID  int                            `yaml:"area"`
		AreaMap [12][10]int                    `yaml:"map"`
		Objects map[int]setup.ObjectProperties `yaml:"objects"`
		Flags   map[string]bool                `yaml:"flags"`
		Moves   int                            `yaml:"moves"`
	}

	filename := folderListing()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		r = setup.Reactions["noSaveFile"]
		areaID = area.ID
		return
	}

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		r.Statement[0] = err.Error()
		r.OK = false
		return
	}
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		r.Statement[0] = err.Error()
		r.OK = false
		return
	}

	setup.GameObjects = content.Objects
	setup.Map = content.AreaMap
	if setup.GetObjectByID(47).Properties.Area == setup.INUSE {
		setup.Verbs = setup.AddMapVerb(setup.Verbs)
	}

	setup.Flags = content.Flags
	setup.Moves = content.Moves

	r = setup.Reactions["loaded"]
	areaID = content.AreaID
	return
}
