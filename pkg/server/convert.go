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
		RemoteURL: in.GetRemoteUrl(),
		Version:   in.GetVersion(),
		FileInfo:  make([]*model.FileInfo, 0),
		Type: in.GetType(),
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
		version   = in.GetVersion()
		repoType  = in.GetType()
		remoteUrl = in.GetRemoteUrl()
		fileInfo  = in.GetFileInfo()
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

		for _, fn := range fi.FileNames {
			if fn == "" {
				return errors.New("invalid file name")
			}
		}


	}

	return nil
}
