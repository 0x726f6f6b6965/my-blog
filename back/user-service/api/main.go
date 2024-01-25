package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/0x726f6f6b6965/my-blog/lib/checker"
	"github.com/0x726f6f6b6965/my-blog/lib/logger"
	pbUser "github.com/0x726f6f6b6965/my-blog/protos/user/v1"

	"github.com/0x726f6f6b6965/my-blog/user-service/internal/client"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/config"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/helper"
	"github.com/0x726f6f6b6965/my-blog/user-service/internal/services"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

func main() {
	godotenv.Load()
	path := os.Getenv("CONFIG")
	var cfg config.Config
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read yaml error", err)
		return
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("unmarshal yaml error", err)
		return
	}
	zaplog, cleanup, err := logger.NewLogger(&cfg.Log)
	if err != nil {
		log.Fatal("create log error", err)
		return
	}
	defer cleanup()

	db, dbCleanup, err := client.NewPostgres(&cfg.DB)
	if err != nil {
		zaplog.Error("failed to connect db", zap.Error(err))
		return
	}
	defer dbCleanup()

	rds := client.InitRedisClient(&cfg.Redis)
	defer rds.Close()
	//Init secret
	if secret, _ := client.GetSecret(context.Background(), rds); checker.IsEmpty(secret) {
		err = client.SetSecret(context.Background(), helper.GetRandString(), -1, rds)
		if err != nil {
			zaplog.Error("failed to set secret", zap.Error(err))
			return
		}
	}

	server := services.NewUserService(cfg.Expire, db, rds)
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Grpc.Port))

	if err != nil {
		zaplog.Error("failed to listen", zap.Error(err))
		return
	}

	pbUser.RegisterUserServiceServer(grpcServer, server)
	zaplog.Info("server listening", zap.String("addr", lis.Addr().String()))
	if err := grpcServer.Serve(lis); err != nil {
		zaplog.Error("failed to serve", zap.Error(err))
		return
	}
}
