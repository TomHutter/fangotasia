package actions

import (
	"fantasia/setup"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

func folderListing() (filename string) {
	/*
		files, err := ioutil.ReadDir(pathname)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Files:")

		for _, f := range files {
			fmt.Println(f.Name())
		}
		//filename = pathname + view.Scanner("prompt: filename > ")
	*/
	filename = setup.PathName + "/save/fantasia.sav"
	return
}

func (obj Object) Save(area setup.Area) (r setup.Reaction) {
	m := make(map[interface{}]interface{})
	m["area"] = area.ID
	m["map"] = setup.Map
	m["objects"] = setup.GameObjects

	filename := folderListing()

	r = setup.Reactions["saved"]
	file, err := os.Create(filename)
	if err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}

	defer file.Close()

	if err = yaml.NewEncoder(file).Encode(m); err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}

	err = file.Close()
	if err != nil {
		r.Statement = err.Error()
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
	}

	filename := folderListing()

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}
	err = yaml.Unmarshal(yamlFile, &content)
	if err != nil {
		r.Statement = err.Error()
		r.OK = false
		return
	}

	setup.GameObjects = content.Objects
	setup.Map = content.AreaMap

	r = setup.Reactions["loaded"]
	areaID = content.AreaID
	return
}
