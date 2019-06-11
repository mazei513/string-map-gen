//go:generate go run github.com/mazei513/string-map-gen robot.go robot

package example

type robot string

// const robot_comment = "no!"

const robot_Chicken = "BAWK!"
const robotCow = "MOO!"

var robot_variable = "no!"

const (
	robot_R2D2 = "R2-D2"
	robot_C3P0 = "C3P0"
	foobar     = 2
)

type bar struct {
	robotVar string
}

func (r robot) String() string {
	const robot_local = "no!"
	return string(r)
}

func anotherFunc() string {
	const robot_local = "no!"
	return robot_local
}

func robotFunc() bool {
	return true
}
