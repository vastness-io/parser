package server

import (
	"errors"
	"github.com/vastness-io/parser-svc"
	"github.com/vastness-io/parser/pkg/model"
)

func ParserRequestToFileTypes(in *parser.ParserRequest) (*model.Repository, error) {
	if err := validateParserRequest(in); err != nil {
		return nil, err
	}

	out := &model.Repository{
		RemoteURL: in.RemoteUrl,
		Version:   in.Version,
		FileInfo:  make([]*model.FileInfo, 0),
	}
	for _, fi := range in.FileInfo {
		var (
			fileNames = fi.FileNames
			language  = fi.Language
		)

		out.FileInfo = append(out.FileInfo, &model.FileInfo{
			Language:  language,
			FileNames: fileNames,
		})

	}
	return out, nil
}

func validateParserRequest(in *parser.ParserRequest) error {
	var (
		version   = in.Version
		repoType  = in.Type
		remoteUrl = in.RemoteUrl
		fileInfo  = in.FileInfo
	)

	if remoteUrl == "" {
		return errors.New("invalid remoteUrl")
	}

	if version == "" {
		return errors.New("invalid version")
	}

	if repoType == "" {
		return errors.New("invalid type")
	}

	if len(fileInfo) == 0 {
		return errors.New("no file info")
	}

	for _, fi := range in.FileInfo {
		if fi.Language == "" {
			return errors.New("invalid language")
		}

		if len(fi.FileNames) == 0 {
			return errors.New("no files avaliable for language")
		}
	}

	return nil
}
