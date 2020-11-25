package data

import "testing"

func TestCheckValidation(t *testing.T) {
	p := &Post{}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
