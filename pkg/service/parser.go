package service

import (
	"errors"
	"github.com/vastness-io/parser-svc"
	"github.com/vastness-io/parser/pkg/model"
	internal "github.com/vastness-io/parser/pkg/parser"
	"github.com/vastness-io/parser/pkg/vcs"
	"io/ioutil"
)

var (
	NoParserAvailable = errors.New("no parser for type specified")
)
type ParserService interface {
	Parse(*model.Repository) (*parser.ParserResponse, error)
}

type parserService struct {
	vcsSet        vcs.VcsSet
	typeParserSet internal.TypeParserSet
}

func NewParserService(set vcs.VcsSet, parserSet internal.TypeParserSet) ParserService {
	return &parserService{
		vcsSet: set,
		typeParserSet: parserSet,
	}
}

func (ps *parserService) Parse(repository *model.Repository) (*parser.ParserResponse, error) {

	v, err := vcs.DetectVcs(ps.vcsSet, repository.Type)

	if err != nil {
		return nil, err
	}

	defer v.Cleanup()

	err = v.Clone(repository.RemoteURL)

	if err != nil {
		return nil, err
	}

	err = v.Checkout(repository.Version)

	if err != nil {
		return nil, err
	}

	out := parser.ParserResponse{}

	for _, fi := range repository.FileInfo {

		p, err := ps.IsParsable(fi.Language)

		if err != nil {
			return nil, err
		}

		for _, fileName := range fi.FileNames {

			file, err := v.Open(fileName)

			if err != nil {
				return nil, err
			}

			b, err := ioutil.ReadAll(file)

			if err != nil {
				return nil, err
			}

			t, err := p.Parse(b)

			if err != nil {
				return nil, err
			}

			switch con := t.(type) {

			case *internal.MavenPom:
				out.MavenResponse = append(out.MavenResponse, &parser.MavenResponse{
					GroupId:    con.GroupID,
					ArtifactId: con.ArtifactID,
					Version:    con.Version,
					Parent: &parser.MavenResponse_Parent{
						GroupId: con.Parent.GroupID,
						ArtifactId: con.Parent.ArtifactID,
						Version: con.Parent.Version,
						RelativePath: con.Parent.RelativePath,
					},
					Modules: con.Modules,
					Properties: &parser.MavenResponse_Properties{
						JavaVersion: con.Properties.JavaVersion,
					},
				})
			}

		}

	}

	return &out, nil
}

func (ps *parserService) IsParsable(language string) (p internal.TypeParser, err error) {
	switch language {
	case "Maven POM":
		p = ps.typeParserSet.Maven()
	default:
		err = NoParserAvailable
	}
	return
}
