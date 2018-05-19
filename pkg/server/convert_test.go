package server

import (
	"testing"
	"errors"
	"github.com/vastness-io/parser-svc"
	"reflect"
)

func TestValidateParserRequest(t *testing.T) {
	tests := []struct{
		version   string
		repoType  string
		remoteUrl string
		fileInfo  []*parser.FileInfo
		err error
	}{
		{
			version: "",
			repoType: "a",
			remoteUrl: "a",
			fileInfo: []*parser.FileInfo {
				{
					Language: "not_blank",
					FileNames: []string {
						"file",
					},
				},
			},
			err: errors.New("invalid version"),
		},
		{
			version: "a",
			repoType: "",
			remoteUrl: "a",
			fileInfo: []*parser.FileInfo {
				{
					Language: "not_blank",
					FileNames: []string {
						"file",
					},
				},
			},
			err: errors.New("invalid type"),
		},
		{
			version: "a",
			repoType: "a",
			remoteUrl: "",
			fileInfo: []*parser.FileInfo {
				{
					Language: "not_blank",
					FileNames: []string {
						"file",
					},
				},
			},
			err: errors.New("invalid remoteUrl"),
		},
		{
			version: "a",
			repoType: "a",
			remoteUrl: "a",
			fileInfo: []*parser.FileInfo {
				{
					Language: "",
					FileNames: []string {
						"file",
					},
				},
			},
			err: errors.New("invalid language"),
		},
		{
			version: "a",
			repoType: "a",
			remoteUrl: "a",
			fileInfo: []*parser.FileInfo {
				{
					Language: "not_blank",
					FileNames: []string {
						"",
					},
				},
			},
			err: errors.New("invalid file name"),
		},
		{
			version: "a",
			repoType: "a",
			remoteUrl: "a",
			fileInfo: []*parser.FileInfo {
				{
					Language: "not_blank",
					FileNames: []string {
						"some_name",
					},
				},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		err := validateParserRequest(&parser.ParserRequest{
			Version: test.version,
			RemoteUrl: test.remoteUrl,
			Type: test.repoType,
			FileInfo: test.fileInfo,
		})

		if !reflect.DeepEqual(err,test.err) {
			t.Fatalf("expected %v, got %v", test.err, err)
		}
	}
}