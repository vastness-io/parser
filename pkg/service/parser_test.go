package service

import (
	"testing"
	internal_parser "github.com/vastness-io/parser/pkg/parser"
	"github.com/vastness-io/parser/pkg/vcs"
	"github.com/vastness-io/parser/pkg/vcs/git"
	"reflect"
	"github.com/vastness-io/parser/pkg/model"
	"github.com/vastness-io/parser-svc"
	"path/filepath"
)

func TestParserService_IsParsable(t *testing.T) {

	tests := []struct {
		language string
		vcsSet vcs.VcsSet
		typeParserSet internal_parser.TypeParserSet
		expectedParser internal_parser.TypeParser
		expectedErr error
	}{
		{
			language: "Maven POM",
			vcsSet: vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			expectedParser: &internal_parser.MavenPomParser{},
			expectedErr: nil,
		},
		{
			language: "Go",
			vcsSet: vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			expectedParser: nil,
			expectedErr: NoParserAvailable,
		},
	}

	for _, test := range tests {

		svc := &parserService{
			vcsSet: test.vcsSet,
			typeParserSet: test.typeParserSet,
		}

		typeParser, err := svc.IsParsable(test.language)

		if err != test.expectedErr {
			t.Fatalf("expected %v, got %v", test.expectedErr, err)
		}

		if !reflect.DeepEqual(test.expectedParser, typeParser) {
			t.Fatalf("expected %v, got %v", test.expectedParser, typeParser)
		}

	}
}

func TestParserService_Parse(t *testing.T) {

	tests := []struct {
		vcsSet vcs.VcsSet
		typeParserSet internal_parser.TypeParserSet
		repository *model.Repository
		expectedResponse *parser.ParserResponse
		expectedErr error
	}{
		{
			vcsSet: vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			repository: &model.Repository{
				RemoteURL: getAbsolutePath("../../test-helpers"),
				Version: "with-pom",
				FileInfo: []*model.FileInfo {
					{
						Language: "Maven POM",
						FileNames: []string {
							"pom.xml",
							"b/pom.xml",
						},
					},
				},
				Type: "GITHUB",
			},
			expectedResponse: &parser.ParserResponse{
				MavenResponse: []*parser.MavenResponse{
					{
						GroupId: "io.vastness",
						ArtifactId: "a",
						Version: "1.0",
						Parent: new(parser.MavenResponse_Parent),
						Modules: []string {
							"b",
						},
						Properties: &parser.MavenResponse_Properties{
							JavaVersion: "1.8",
						},
					},
					{
						GroupId: "io.vastness",
						ArtifactId: "b",
						Parent: &parser.MavenResponse_Parent{
							GroupId: "io.vastness",
							ArtifactId: "a",
							Version: "1.0",
						},
						Modules: nil,
						Properties: &parser.MavenResponse_Properties{
							JavaVersion: "1.8",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			vcsSet: vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			repository: &model.Repository{
				RemoteURL: getAbsolutePath("../../test-helpers"),
				Version: "with-pom",
				FileInfo: []*model.FileInfo {
					{
						Language: "Go",
						FileNames: []string {
							"main.go",
						},
					},
				},
				Type: "GITHUB",
			},
			expectedResponse: nil,
			expectedErr: NoParserAvailable,
		},

	}

	for _, test := range tests {

		svc := NewParserService(test.vcsSet,test.typeParserSet)

		res, err := svc.Parse(test.repository)

		if err != test.expectedErr {
			t.Fatalf("expected %v, got %v", test.expectedErr, err)
		}

		if !reflect.DeepEqual(test.expectedResponse, res) {
			t.Fatalf("expected %v, got %v", test.expectedResponse, res)
		}

	}
}

func getAbsolutePath(relative string) string {
	abs, err := filepath.Abs(relative)

	if err != nil {
		return ""
	}

	return abs
}
