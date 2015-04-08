// +build none

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"labix.org/v2/pipe"
)

type Matcher struct {
	Text string
	M    []string
}

func (m *Matcher) Find(pattern string) bool {
	re := regexp.MustCompile(pattern)
	m.M = re.FindStringSubmatch(m.Text)
	return len(m.M) > 0
}

func (m *Matcher) ReplaceAll(pattern, replacement string) {
	re := regexp.MustCompile(pattern)
	m.Text = re.ReplaceAllString(m.Text, replacement)
}

func publicize(name string) string {
	return strings.ToUpper(string(name[0])) + name[1:]
}

func main() {
	_ = pipe.Write

	class := ""
	replacements := map[string]string{}

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		line := Matcher{Text: scan.Text()}
		// FIXME(akavel): rescue strings and comments
		line.ReplaceAll(`\$`, ``)
		line.ReplaceAll(`->`, `.`)
		line.ReplaceAll(`<<<SVG`, "`")
		line.ReplaceAll(`^SVG;$`, "`")
		line.ReplaceAll(`\.=`, `+=`)
		line.ReplaceAll(`\bnew\s+`, `New`)
		for before, after := range replacements {
			line.ReplaceAll(before, after)
		}

		switch {

		// class X -> type X struct
		case line.Find(`^class\s+(\w+)`):
			class = line.M[1]
			replacements["self::"] = class + "_"
			fmt.Printf("type %s struct%s\n",
				class, line.Text[len(line.M[0]):])

		case line.Find(`^(\s*)public\s+static\s+function\s+(\w+)`):
			oldname := fmt.Sprintf("%s::%s", class, line.M[2])
			newname := fmt.Sprintf("%s_%s", class, publicize(line.M[2]))
			replacements[oldname] = newname
			fmt.Printf("%sfunc %s%s\n",
				line.M[1], newname, line.Text[len(line.M[0]):])

		// print rest of lines intact
		default:
			fmt.Println(line.Text)
		}
	}
}
