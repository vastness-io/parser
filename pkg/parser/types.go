package parser

type TypeParser interface {
	Parse([]byte) (interface{}, error)
}

type TypeParserSet interface {
	Maven() TypeParser
}

type typeParserSet struct {
	mavenPomParser TypeParser
}

func (tps *typeParserSet) Maven() TypeParser {
	return tps.mavenPomParser
}

func NewTypeParserSet(mavenParser TypeParser) TypeParserSet {
	return &typeParserSet{
		mavenPomParser: mavenParser,
	}
}
