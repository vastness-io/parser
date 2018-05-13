package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/vastness-io/parser-svc"
	"github.com/vastness-io/parser/pkg/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type parserServer struct {
	parserService service.ParserService
	log           *logrus.Entry
}

func NewParserServer(parserService service.ParserService, log *logrus.Entry) parser.ParserServer {
	return &parserServer{
		parserService: parserService,
		log:           log,
	}
}

func (ps *parserServer) Analyse(ctx context.Context, req *parser.ParserRequest) (*parser.ParserResponse, error) {

	repository, err := ParserRequestToFileTypes(req)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return ps.parserService.Parse(repository)
}
