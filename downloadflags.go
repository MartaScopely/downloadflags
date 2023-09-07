package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	svgo "github.com/ajstarks/svgo"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/svg"
)

type Country struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type Value struct {
	Currency  string `json:"currency"`
	Name      string `json:"name"`
	Continent string `json:"continent"`
}

func main() {
	jsonFile, err := os.Open("countries.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var countries []Country

	json.Unmarshal(byteValue, &countries)

	err = os.MkdirAll("flags_optimized", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	for _, country := range countries {
		code := strings.ToLower(country.Key)
		url := fmt.Sprintf("https://hatscripts.github.io/circle-flags/flags/%s.svg", code)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Failed to download", url)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Use the minify package to optimize the SVG
		m := minify.New()
		m.AddFunc("image/svg+xml", svg.Minify)
		optimizedData, err := m.Bytes("image/svg+xml", data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		filename := filepath.Join("flags_optimized", code+".svg")
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer file.Close()

		// Use the svgo package to write the optimized SVG to the file
		svgfile := svgo.New(file)
		svgfile.Writer.Write(optimizedData)
		fmt.Println("Downloaded and optimized", filename)
	}
}
