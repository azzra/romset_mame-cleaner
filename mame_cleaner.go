package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var datFile = flag.String("dat_file", "", "The DAT file.")
var dryRun = flag.Bool("dry_run", true, "Print what will be moved.")
var regFav = flag.String("reg_fav", "world,europe,eu,france,french,fr,usa,us", "Favorites region(s).")
var romDir = flag.String("rom_dir", ".", "The directory containing the roms file to process.")
var prefLast = flag.Bool("pref_last", true, "If more than 1 rom match, choose the last one.")

type Query struct {
	Header  Header `xml:"header"`
	RomList []Rom  `xml:"game"`
}

type Header struct {
	Name    string `xml:"name"`
	Version string `xml:"version"`
}

type Rom struct {
	Name        string `xml:"name,attr"`
	CloneOf     string `xml:"cloneof,attr"`
	RomOf       string `xml:"romof,attr"`
	IsBios      string `xml:"isbios,attr"`
	Description string `xml:"description"`
	Region      string
	Date        int
}

type Game struct {
	Parent   Rom
	Children []Rom
}

func extractGames(romList []Rom) map[string]Game {

	games := make(map[string]Game)

	for _, rom := range romList {

		if rom.CloneOf == "" {
			game := games[rom.Name]
			game.Parent = rom
			game.Children = append([]Rom{rom}, game.Children...) // prepend
			games[rom.Name] = game
		} else {
			game := games[rom.CloneOf]
			game.Children = append(game.Children, rom)
			games[rom.CloneOf] = game
		}

	}

	return games
}

func processGames(games map[string]Game) {

	var foundGame *Rom

	for _, game := range games {

		foundGame = nil

		if game.Parent.Name != "" {

			if len(game.Children) > 1 {
				fmt.Println("Choosing for: ", game.Parent.Name, "-- ", game.Parent.Description, "// found", len(game.Children))
				foundGame = findMatchingRom(game)
				fmt.Println("    found in multiple:", foundGame.Description)
			} else {
				foundGame = &game.Parent
				fmt.Println("Only one:", foundGame.Description)
			}

		} else {
			fmt.Println("NO PARENT FOUND:", game.Children[0].Name)
		}
	}

}

func findMatchingRom(game Game) *Rom {

	var found []*Rom
	var length int

	for _, region := range strings.Split(*regFav, ",") {

		//fmt.Println("	looking for", region)
		for i, childrenLength := 0, len(game.Children); i < childrenLength; i++ {

			rom := game.Children[i]

			romRegion, _ := extractAttributes(&rom)

			//fmt.Println("		" + rom.Name + " - region:", romRegion, ", sum:", sum)
			if strings.Index(romRegion, region) != -1 {
				found = append(found, &rom)
			}
		}

		length = len(found)
		if length > 0 {
			break
		}

	}

	if length == 1 || (length > 1 && *prefLast == false) {
		return found[0]
	}

	if length > 0 {
		return found[length-1]
	}

	// assuming game.Children contains at least 2 elements
	if *prefLast == true {
		return &game.Children[len(game.Children)-1]
	}

	return &game.Parent
}

func extractAttributes(rom *Rom) (string, int) {

	desc := rom.Description
	posAttrsBegin, posAttrsEnd := strings.Index(desc, "("), strings.LastIndex(desc, ")")

	if posAttrsEnd <= posAttrsBegin || posAttrsBegin == -1 {
		return "", 0
	}

	attrs := desc[posAttrsBegin+1 : posAttrsEnd]
	romChar := regexp.MustCompile("(?i)[^a-z ]*").ReplaceAllLiteralString(attrs, "")
	romDigit, _ := strconv.Atoi(regexp.MustCompile("(?i)[^0-9]*").ReplaceAllLiteralString(attrs, ""))

	return strings.TrimSpace(strings.ToLower(romChar)), romDigit

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

	processGames(games)
}
