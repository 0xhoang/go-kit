package must

import (
	"context"
	"fmt"
	"github.com/0xhoang/go-kit/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServiceServer interface {
	RegisterGrpcServer(*grpc.Server)
	RegisterHandler(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

func NewServer(
	ctx context.Context,
	cfg *config.Config,
	opt []grpc.ServerOption,
	serviceServer ...ServiceServer,
) *grpc.Server {
	grpcPort := fmt.Sprintf(":%d", cfg.GrpcPort)

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer(
		opt...,
	)
	// Attach the Greeter service to the server

	for _, server := range serviceServer {
		server.RegisterGrpcServer(s)
	}

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0" + grpcPort)

	go func() {
		log.Fatalln(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		ctx,
		"0.0.0.0"+grpcPort,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	gwmux := runtime.NewServeMux(
		runtime.WithRoutingErrorHandler(handleRoutingError),
	)

	for _, server := range serviceServer {
		err = server.RegisterHandler(ctx, gwmux, conn)
		if err != nil {
			log.Fatalln("Failed to register gateway:", err)
		}
	}

	httpPort := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Handler: gwmux,
		Addr:    httpPort,
	}

	go func() {
		log.Printf("Serving gRPC-Gateway on %s\n", httpPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("srv.ListenAndServe", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatal("Server forced to shutdown", zap.Error(err))
	}

	log.Println("Server exiting")

	return nil
}
