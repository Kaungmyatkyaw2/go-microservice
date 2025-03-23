package data

import (
	"testing"
)

func TestCheckValidation(t *testing.T) {
	p := &Product{
		Name:  "Tea",
		Price: 1,
		SKU:   "abs-abs-abs",
	}

	validate := NewValidation()

	errs := validate.Validate(p)

	if len(errs) != 0 {
		t.Fatal(errs.Errors())
	}

}
