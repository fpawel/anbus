package work

import (
	"github.com/pkg/errors"
	"strings"
)

type txtCmd struct {
	source string
}

func (x txtCmd) tokens() []string {
	xs := append([]string{}, strings.Split(x.source, " ")...)
	for i := range xs {
		xs[i] = strings.ToUpper(xs[i])
	}
	return xs
}

func (x txtCmd) name() string {
	xs := x.tokens()
	if len(xs) == 0 {
		return ""
	}
	return xs[0]
}

func parseTxtCmd(sourceStr string) (txtCmd, error) {
	x := txtCmd{
		source: sourceStr,
	}
	if len(x.tokens()) == 0 {
		return x, errors.New("команда не задана")
	}
	return x, nil
}
