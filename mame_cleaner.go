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
	RomList []Rom `xml:"game"`
}

type Header struct {
	Name string  `xml:"name"`
	Version string  `xml:"version"`
}

type Rom struct {
	Name string `xml:"name,attr"`
	CloneOf string `xml:"cloneof,attr"`
	RomOf string `xml:"romof,attr"`
	IsBios string `xml:"isbios,attr"`
	Description string `xml:"description"`
	Region string
	Date int
}

type Game struct {
	Parent Rom
	Children []Rom
}

func extractGames(romList []Rom) map[string]Game {

	games := make(map[string]Game)

	for _, rom := range romList[20:80] {
	
		fmt.Printf("%v", rom)

		if rom.CloneOf == "" {
			fmt.Println("ajout de " + rom.Name)
			game := games[rom.Name]
			game.Parent = rom
			games[rom.Name] = game
		} else {
			game := games[rom.CloneOf]
			fmt.Println("   ajout enfant  pour " + rom.CloneOf + " -> " + rom.Name)
			game.Children = append(game.Children, rom)
			games[rom.CloneOf] = game
		}
		

	}

	return games

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


func main() {

	flag.Parse()
	q := parseXml()

	fmt.Println("ROM DIR: " + *romDir)
	fmt.Println("DAT FILE: " + *datFile)
	
	games := extractGames(q.RomList)

	for _, game := range games {
		fmt.Println("    Game " + game.Parent.Name)		
		fmt.Println(game.Parent)		
		fmt.Println(game.Children)		
	}
}

