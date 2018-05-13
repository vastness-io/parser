package parser

import (
	"encoding/xml"
)

type MavenPomParser struct{}

func (p *MavenPomParser) Parse(b []byte) (interface{}, error) {
	out := NewMavenPom()

	err := xml.Unmarshal(b, out)

	if err != nil {
		return nil, err
	}

	return out, nil
}
