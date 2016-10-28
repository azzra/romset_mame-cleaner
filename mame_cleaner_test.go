package main

import (
	"testing"
)

func TestExtractAttributesValid(t *testing.T) {

	rom := Rom{Name: "foobar", Description: "Foobar (USA) (bootleg 1 990101)"}

	region, sum := extractAttributes(&rom)

	expectedRegion := "usa bootleg"
	if region != expectedRegion {
		t.Errorf("extractAttributes(%v) should return '%v'", rom, expectedRegion)
	}

	expectedSum := 1990101
	if sum != expectedSum {
		t.Errorf("extractAttributes(%v) should return '%v'", rom, expectedSum)
	}

}

func TestExtractAttributesInvalid(t *testing.T) {

	for _, rom := range []Rom{
		Rom{Name: "foobar", Description: "Foobar 123"},
		Rom{Name: "foobar", Description: "Foobar (123"},
		Rom{Name: "foobar", Description: "Foobar )123"},
		Rom{Name: "foobar", Description: "Foobar"}} {

		region, sum := extractAttributes(&rom)
		expectedRegion := ""
		if region != expectedRegion {
			t.Errorf("extractAttributes(%v) should return '%v'", rom, expectedRegion)
		}

		expectedSum := 0
		if sum != expectedSum {
			t.Errorf("extractAttributes(%v) should return '%v'", rom, expectedSum)
		}

	}

}
