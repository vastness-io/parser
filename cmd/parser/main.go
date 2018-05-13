package main

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"fmt"
	"os"
	"net"
	"strconv"
	toolkit "github.com/vastness-io/toolkit/pkg/grpc"
	"github.com/vastness-io/parser/pkg/vcs"
	"github.com/vastness-io/parser/pkg/parser"
	"github.com/vastness-io/parser/pkg/service"
	"github.com/vastness-io/parser/pkg/server"
	"os/signal"
	"syscall"
	"github.com/vastness-io/parser/pkg/vcs/git"
	svc "github.com/vastness-io/parser-svc"
	"github.com/opentracing/opentracing-go"
)

const (
	name        = "parser"
	description = "Analyses files deeply"
)

var (
	log       = logrus.WithField("component", name)
	commit    string
	version   string
	addr      string
	port      int
	debugMode bool
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Usage = description
	app.Version = fmt.Sprintf("%s (%s)", version, commit)
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "addr,a",
			Usage:       "TCP address to listen on",
			Value:       "127.0.0.1",
			Destination: &addr,
		},
		cli.IntFlag{
			Name:        "port,p",
			Usage:       "Port to listen on",
			Value:       8083,
			Destination: &port,
		},
		cli.BoolFlag{
			Name:        "debug,d",
			Usage:       "Debug mode",
			Destination: &debugMode,
		},
	}
	app.Action = func(_ *cli.Context) { run() }
	app.Run(os.Args)
}

func run() {

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}

	log.Infof("Starting %s", name)

	var (
		address        = net.JoinHostPort(addr, strconv.Itoa(port))
		tracer         = opentracing.GlobalTracer()
		lis, err       = net.Listen("tcp", address)
		srv            = toolkit.NewGRPCServer(tracer, log)
		gitVcs         = git.InMemoryGit{}
		vcsSet         = vcs.NewVcsSet(&gitVcs)
		mavenPomParser = parser.MavenPomParser{}
		typeParserSet  = parser.NewTypeParserSet(&mavenPomParser)
		parserSvc      = service.NewParserService(vcsSet,typeParserSet)
	)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Infof("Listening on %s", address)
		if err := srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	svc.RegisterParserServer(srv, server.NewParserServer(parserSvc, log))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			log.Infof("Exiting %s", name)
			srv.GracefulStop()
			os.Exit(0)
		}
	}
}
