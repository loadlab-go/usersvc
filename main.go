package main

import (
	"flag"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	userpb "github.com/loadlab-go/pkg/proto/user"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	flagListenAddr      = flag.String("listen", ":8080", "listen address")
	flagEtcdEndpoints   = flag.String("etcd-endpoints", "localhost:2379", "etcd endpoints")
	flagAdvertiseClient = flag.String("advertise-client", "localhost:8080", "advertise client url")
	flagPostgresDSN     = flag.String("postgres-dsn", "postgres://loadlab:111111@localhost/loadlab?sslmode=disable", "postgres DSN")
)

func main() {
	flag.Parse()

	tracerCloser, err := initTracer()
	if err != nil {
		logger.Panic("init tracer failed", zap.Error(err))
	}
	defer tracerCloser.Close()

	initialize()

	srv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_recovery.StreamServerInterceptor(),
			otgrpc.OpenTracingStreamServerInterceptor(opentracing.GlobalTracer()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
		)))

	us := &userSvc{repo: db}
	userpb.RegisterUserServer(srv, us)

	l, err := net.Listen("tcp", *flagListenAddr)
	if err != nil {
		logger.Fatal("listen failed", zap.Error(err))
	}
	logger.Info("user service server startup", zap.String("listen", l.Addr().String()))

	go signalSet(srv.GracefulStop)

	err = srv.Serve(l)
	if err != nil {
		logger.Panic("user service  server serve failed", zap.Error(err))
	}
	logger.Info("server stopped")
}

func initialize() {
	mustInitDB("postgres", *flagPostgresDSN)

	mustInitEtcdCli(*flagEtcdEndpoints)

	go func() {
		err := registerEndpointWithRetry(*flagAdvertiseClient)
		if err != nil {
			logger.Panic("register endpoint faield", zap.Error(err))
		}
	}()
}

func signalSet(cb func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-sigCh
	logger.Warn("exit signal", zap.String("signal", s.String()))

	cb()
}

func initTracer() (io.Closer, error) {
	cfg := config.Configuration{ServiceName: "usersvc"}
	_, err := cfg.FromEnv()
	if err != nil {
		return nil, err
	}
	tracer, tracerCloser, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}
	opentracing.InitGlobalTracer(tracer)
	return tracerCloser, nil
}
