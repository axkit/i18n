package i18n

import (
	"bufio"
	"bytes"
	"strings"
)

// HintSeparator holds a separator between value and hint in a .t18n file.
var HintSeparator = "//"

type DefaultParser struct{}

func (p *DefaultParser) ParseFileContent(data []byte) ([]Item, error) {

	var res []Item
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for {
		if !scanner.Scan() {
			break
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

		line := scanner.Text()
		line = strings.TrimLeft(line, " ")
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		item := p.parseLine(line)
		if item != nil {
			res = append(res, *item)
		}
	}

	return res, nil
}

func (p *DefaultParser) parseLine(line string) *Item {
	var res Item

	vx := strings.Index(line, "=")
	if vx == -1 {
		return nil
	}

	res.Key = strings.TrimSpace(line[0:vx])
	val := strings.TrimSpace(line[vx+1:])
	hx := strings.Index(val, HintSeparator)
	if hx != -1 {
		res.Hint = strings.TrimSpace(val[hx+len(HintSeparator):])
		res.Value = strings.TrimSpace(val[0:hx])
	} else {
		res.Value = val
	}

	return &res
}
