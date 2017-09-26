package mppuma

import "testing"

func TestGraphDefinition(t *testing.T) {
	desired := 3
	var puma PumaPlugin

	graphdef := puma.GraphDefinition()

	if len(graphdef) != desired {
		t.Errorf("GraphDefinition: %d should be %d", len(graphdef), desired)
	}
}
