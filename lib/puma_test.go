package mppuma

import "testing"

func TestGraphDefinition(t *testing.T) {

	var puma PumaPlugin

	graphdef := puma.GraphDefinition()

	if len(graphdef) != 2 {
		t.Errorf("GraphDefinition: %d should be %d", len(graphdef), 2)
	}
}

func TestGraphDefinitionWithGC(t *testing.T) {

	var puma PumaPlugin
	puma.WithGC = true

	graphdef := puma.GraphDefinition()

	if len(graphdef) != 3 {
		t.Errorf("GraphDefinitionWithGC: %d should be %d", len(graphdef), 3)
	}
}
