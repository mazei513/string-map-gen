package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPrefixedNames(t *testing.T) {
	//GIVEN
	f := `
	package example

	type robot string
	
	const (
		robot_R2D2 = "R2-D2"
		robot_C3P0 = "C3P0"
	)
	`
	prefix := "robot"
	//WHEN
	v, err := getPrefixedNames(f, prefix)
	//THEN
	e := []string{"robot_R2D2", "robot_C3P0"}
	assert.NoError(t, err)
	assert.Equal(t, e, v)
}
