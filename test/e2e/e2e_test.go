package e2e

import (
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/dungnh3/guardrails-assignment/pkg/migration"

	"github.com/dungnh3/guardrails-assignment/config"
	"github.com/dungnh3/guardrails-assignment/internal/apps/server"
	"github.com/go-logr/logr"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

const migrationPath = "file://../../db/migrations"

var (
	cfg    *config.Config
	logger logr.Logger
)

type e2eTestSuite struct {
	suite.Suite
	db  *gorm.DB
	svc *server.Service
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, &e2eTestSuite{})
}

func (s *e2eTestSuite) SetupSuite() {
	var err error
	cfg, err = config.Load()
	s.Require().NoError(err)

	zapLogger, err := cfg.Logger.Build()
	s.Require().NoError(err)

	grpcMiddlewareUnary := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(zapLogger),
		grpc_validator.UnaryServerInterceptor(),
	}

	svr := server.NewServer(cfg.Server,
		grpc_middleware.WithUnaryServerChain(grpcMiddlewareUnary...),
	)

	s.db = cfg.PostgreSQL.ConnectDatabase()

	sqlDB, err := s.db.DB()
	s.Require().NoError(err)

	err = migration.Up(sqlDB, migrationPath)
	s.Require().NoError(err)

	svc, err := server.NewService(cfg, s.db)
	s.Require().NoError(err)

	err = svr.Register(svc)
	s.Require().NoError(err)

	go svr.Serve()
	time.Sleep(2 * time.Second)
}

func (s *e2eTestSuite) TearDownSuite() {
	var err error
	err = s.svc.Close()
	s.Require().NoError(err)

	// err = migration.Down(s.sqlDB, migrationPath)
	// s.Require().NoError(err)

	p, _ := os.FindProcess(syscall.Getpid())
	p.Signal(syscall.SIGINT)
}

func uri(path string) string {
	host := fmt.Sprintf("http://localhost:%v", cfg.Server.HTTP.Port)
	return fmt.Sprintf("%v%v", host, path)
}
