package main

import (
	//"errors"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var romDir = flag.String("rom_dir", ".", "The directory containing the roms file to process.")
var dryRun = flag.Bool("dry_run", true, "Print what will be moved.")
var datFile = flag.String("dat_file", "", "The DAT file.")

type Query struct {
	Header Header `xml:"header"`
	GameList []Game `xml:"game"`
}

type Header struct {
	Name string  `xml:"name"`
	Version string  `xml:"version"`
}

type Game struct {
	Name string `xml:"name,attr"`
	CloneOf string `xml:"cloneof,attr"`
	RomOf string `xml:"romof,attr"`
	IsBios string `xml:"isbios,attr"`
	Description string `xml:"description"`
}


func main() {

	flag.Parse()
	q := parseXml()

	fmt.Println("ROM DIR: " + *romDir)
	fmt.Println("DAT FILE: " + *datFile)

	for _, game := range q.GameList {
		if game.CloneOf == "" {
			fmt.Printf("\t%s\n", game)
		}
	}

	fmt.Println(q.Header)
	fmt.Println(len(q.GameList))
}

func parseXml() Query {

	xmlFile, err := os.Open(*datFile)
	var q Query
	
	if err != nil {
		panic(err)
	}

	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(b, &q)

	return q
}
