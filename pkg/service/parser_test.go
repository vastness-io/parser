package service

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/vastness-io/parser-svc"
	"github.com/vastness-io/parser/pkg/mock/vcs"
	"github.com/vastness-io/parser/pkg/model"
	internal_parser "github.com/vastness-io/parser/pkg/parser"
	shared_test "github.com/vastness-io/parser/pkg/shared/test"
	"github.com/vastness-io/parser/pkg/vcs"
	"github.com/vastness-io/parser/pkg/vcs/git"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParserService_IsParsable(t *testing.T) {

	tests := []struct {
		language       string
		vcsSet         vcs.VcsSet
		typeParserSet  internal_parser.TypeParserSet
		expectedParser internal_parser.TypeParser
		expectedErr    error
	}{
		{
			language:       "Maven POM",
			vcsSet:         vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet:  internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			expectedParser: &internal_parser.MavenPomParser{},
			expectedErr:    nil,
		},
		{
			language:       "Go",
			vcsSet:         vcs.NewVcsSet(&git.InMemoryGit{}),
			typeParserSet:  internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			expectedParser: nil,
			expectedErr:    NoParserAvailable,
		},
	}

	for _, test := range tests {

		svc := &parserService{
			vcsSet:        test.vcsSet,
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
		typeParserSet    internal_parser.TypeParserSet
		setupVcsFunc     func(mockVcs *mock_vcs.MockVcs, repository *model.Repository)
		repository       *model.Repository
		expectedResponse *parser.ParserResponse
		expectedErr      error
	}{
		{
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			setupVcsFunc: func(mockVcs *mock_vcs.MockVcs, repository *model.Repository) {
				mockVcs.EXPECT().Clone(repository.RemoteURL).Return(nil)
				mockVcs.EXPECT().Checkout(repository.Version).Return(nil)
				mockVcs.EXPECT().Open("pom.xml").Return(&shared_test.MockFile{
					FileName: "pom.xml",
					Reader: bytes.NewReader([]byte(`<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>

  <groupId>io.vastness</groupId>
  <artifactId>a</artifactId>
  <version>1.0</version>
  <packaging>pom</packaging>

    <modules>
        <module>b</module>
    </modules>

  <properties>
    <java.version>1.8</java.version>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
  </properties>

</project>`)),
				}, nil)
				mockVcs.EXPECT().Open("b/pom.xml").Return(&shared_test.MockFile{
					FileName: "pom.xml",
					Reader: bytes.NewReader([]byte(`<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>io.vastness</groupId>
    <artifactId>b</artifactId>

    <parent>
        <groupId>io.vastness</groupId>
        <artifactId>a</artifactId>
        <version>1.0</version>
    </parent>

    <properties>
        <java.version>1.8</java.version>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    </properties>

</project>`)),
				}, nil)

				mockVcs.EXPECT().Cleanup().Return(nil)
			},
			repository: &model.Repository{
				RemoteURL: getAbsolutePath("../../test-helpers"),
				Version:   "with-pom",
				FileInfo: []*model.FileInfo{
					{
						Language: "Maven POM",
						FileNames: []string{
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
						GroupId:    "io.vastness",
						ArtifactId: "a",
						Version:    "1.0",
						Parent:     new(parser.MavenResponse_Parent),
						Modules: []string{
							"b",
						},
						Properties: &parser.MavenResponse_Properties{
							JavaVersion: "1.8",
						},
					},
					{
						GroupId:    "io.vastness",
						ArtifactId: "b",
						Parent: &parser.MavenResponse_Parent{
							GroupId:    "io.vastness",
							ArtifactId: "a",
							Version:    "1.0",
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
			typeParserSet: internal_parser.NewTypeParserSet(&internal_parser.MavenPomParser{}),
			setupVcsFunc: func(mockVcs *mock_vcs.MockVcs, repository *model.Repository) {
				mockVcs.EXPECT().Clone(repository.RemoteURL).Return(nil)
				mockVcs.EXPECT().Checkout("with-pom").Return(nil)
				mockVcs.EXPECT().Cleanup()
			},
			repository: &model.Repository{
				RemoteURL: getAbsolutePath("../../test-helpers"),
				Version:   "with-pom",
				FileInfo: []*model.FileInfo{
					{
						Language: "Go",
						FileNames: []string{
							"main.go",
						},
					},
				},
				Type: "GITHUB",
			},
			expectedResponse: nil,
			expectedErr:      NoParserAvailable,
		},
	}

	for _, test := range tests {

		func() {
			var (
				ctrl     = gomock.NewController(t)
				mavenVcs = mock_vcs.NewMockVcs(ctrl)
				vcsSet   = vcs.NewVcsSet(mavenVcs)
				svc      = NewParserService(vcsSet, test.typeParserSet)
			)
			defer ctrl.Finish()

			test.setupVcsFunc(mavenVcs, test.repository)

			res, err := svc.Parse(test.repository)

			if err != test.expectedErr {
				t.Fatalf("expected %v, got %v", test.expectedErr, err)
			}

			if !reflect.DeepEqual(test.expectedResponse, res) {
				t.Fatalf("expected %v, got %v", test.expectedResponse, res)
			}
		}()
	}
}

func getAbsolutePath(relative string) string {
	abs, err := filepath.Abs(relative)

	if err != nil {
		return ""
	}

	return abs
}
