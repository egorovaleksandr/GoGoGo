//go:build !solution

package ciletters

import (
	"strings"
	"text/template"

	_ "embed"
)

//go:embed t.txt
var str string

func MakeLetter(n *Notification) (string, error) {
	var sb strings.Builder
	s := strings.ReplaceAll(str, "\r", "")
	v, err := template.New("email").Funcs(template.FuncMap{"cmdLog": func(s string) []string {
		l := strings.Split(s, "\n")
		if len(l) > 9 {
			return l[9:]
		}
		return l
	}}).Parse(s)
	if err != nil {
		return "", err
	}
	if err = v.Execute(&sb, n); err != nil {
		return "", err
	}
	return sb.String(), nil
}
