package dockerfile

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

//	https://github.com/moby/buildkit/blob/master/frontend/dockerfile/parser/parser
// https://github.com/moby/buildkit/issues/1561

// Dockerfile is struct that contains parameters to interact with a Dockerfile
type Dockerfile struct {
	File        string
	Instruction string
	Value       string
	DryRun      bool
}

// isPositionKey checks if key use the array position like instruction[0][0].
// where first [0] references the instruction position
// and the  second [0] references the element inside instruction line
func isPositionKeys(key string) bool {
	matched, err := regexp.MatchString(`^(.*)\[[[:digit:]]*\]\[[[:digit:]]*\]$`, key)

	if err != nil {
		fmt.Println(err)
	}
	return matched
}

func getPositionKeys(k string) (
	key string,
	instructionPosition int,
	elementPosition int,
	err error) {

	if isPositionKeys(k) {
		re := regexp.MustCompile(`(.*)*\[[[:digit:]]*\]\[[[:digit:]]*\]$`)
		keys := re.FindStringSubmatch(k)
		key = keys[1]

		re = regexp.MustCompile(`^.*\[([[:digit:]]*)\]\[[[:digit:]]*\]$`)

		positions := re.FindStringSubmatch(k)

		instructionPosition, err = strconv.Atoi(positions[1])

		if err != nil {
			fmt.Println(err)
			return "", -1, -1, err
		}

		re = regexp.MustCompile(`^.*\[[[:digit:]]*\]\[([[:digit:]])*\]$`)

		positions = re.FindStringSubmatch(k)

		elementPosition, err = strconv.Atoi(positions[1])

		if err != nil {
			fmt.Println(err)
			return "", -1, -1, err
		}

	} else {
		key = k
		instructionPosition = 0
		elementPosition = 0
		err = nil
	}

	return key, instructionPosition, elementPosition, err

}

// Search for both a Dockerfile instruction and its value to defined in the Dockerfile
// While the dockerfile instruction not case sensitive, its value is
func (d *Dockerfile) search(node *parser.Node) (bool, error) {
	instruction, instructionPosition, elementPosition, err := getPositionKeys(d.Instruction)

	if err != nil {

		return false, err
	}

	i := 0
	for _, n := range node.Children {

		fmt.Println(i)
		if strings.ToUpper(n.Value) == strings.ToUpper(instruction) && i == instructionPosition {

			fmt.Println("Here")

			if n.Next != nil {
				j := 0
				for nod := n.Next; nod != nil && j <= elementPosition; nod = nod.Next {
					if nod.Value == d.Value && elementPosition == j {

						return true, nil
					}
					j++
				}

			}
		} else if strings.ToUpper(n.Value) == strings.ToUpper(instruction) {
			i++
		}

	}
	return false, nil
}
