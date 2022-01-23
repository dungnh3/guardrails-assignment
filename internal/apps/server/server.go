package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dungnh3/guardrails-assignment/api"
	"github.com/dungnh3/guardrails-assignment/pkg/grpc/gateway"
	"github.com/dungnh3/guardrails-assignment/pkg/grpc/health_api"
	"github.com/dungnh3/guardrails-assignment/pkg/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Server structure schema
type Server struct {
	gRPC *grpc.Server
	mux  *runtime.ServeMux
	cfg  server.Config
}

// NewServer return a server instance
func NewServer(cfg server.Config, opt ...grpc.ServerOption) *Server {
	s := &Server{
		gRPC: grpc.NewServer(opt...),
		mux: runtime.NewServeMux(
			runtime.WithMarshalerOption("text/html", &runtime.JSONPb{}),
			gateway.ProtoJSONMarshaler(),
		),
		cfg: cfg,
	}
	return s
}

func (s *Server) Register(svc *Service) error {
	health_api.RegisterHealthCheckServiceServer(s.gRPC, svc)
	if err := health_api.RegisterHealthCheckServiceHandlerFromEndpoint(context.Background(), s.mux,
		fmt.Sprintf(":%v", s.cfg.GRPC.Port), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return err
	}

	api.RegisterGuardRailsServiceServer(s.gRPC, svc)
	if err := api.RegisterGuardRailsServiceHandlerFromEndpoint(context.Background(), s.mux,
		fmt.Sprintf(":%v", s.cfg.GRPC.Port), []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	); err != nil {
		return err
	}
	return nil
}

func (s *Server) Serve() error {
	stop := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", s.mux)

	httpMux.HandleFunc("/swagger.json", serveSwagger)
	fs := http.FileServer(http.Dir("./docs/swagger-ui"))
	httpMux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui", fs))

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.HTTP.Port),
		Handler: httpMux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("Error starting http server, ", err)
			errCh <- err
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPC.Port))
		if err != nil {
			log.Println("Error listening port, ", err)
			errCh <- err
			return
		}
		if err := s.gRPC.Serve(listener); err != nil {
			log.Println("Error starting gRPC server, ", err)
			errCh <- err
		}
	}()

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
		s.gRPC.GracefulStop()
		return nil
	case err := <-errCh:
		return err
	}
}

func serveSwagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "docs/api/api.swagger.json")
}
