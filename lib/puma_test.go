package mppuma

import "testing"

func TestGraphDefinition(t *testing.T) {
	desired := 2
	var puma PumaPlugin

	graphdef := puma.GraphDefinition()

	if len(graphdef) != desired {
		t.Errorf("GraphDefinition: %d should be %d", len(graphdef), desired)
	}
}

func TestGraphDefinitionWithGC(t *testing.T) {
	desired := 3
	var puma PumaPlugin
	puma.WithGC = true

	graphdef := puma.GraphDefinition()

	if len(graphdef) != desired {
		t.Errorf("GraphDefinition: %d should be %d", len(graphdef), desired)
	}
}
