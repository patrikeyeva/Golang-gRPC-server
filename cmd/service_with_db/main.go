package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"homework8/internal/app/reader"
	"homework8/internal/app/server"
	"homework8/internal/infrastructure/kafka"
	"homework8/internal/pkg/db"
	"homework8/internal/pkg/repository/postgresql"
	pb "homework8/pkg/grpcserver"
	"homework8/pkg/logger"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Configuration{
		ServiceName: "grpc-server",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Println("Error create new tracer")
		log.Fatal(err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger.SetGlobal(zapLogger.With(zap.String("component", "grpc-gateway")))
	l := logger.FromContext(ctx)
	ctx = logger.ToContext(ctx, l.With())

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatal(err)
	}

	clientConsumer, err := sarama.NewConsumerGroup(brokers, "service-with-db", kafka.CreateConfigForConsumerGroup())
	if err != nil {
		log.Fatal("Error creating consumer group client:", err)
	}
	defer func() {
		if err := clientConsumer.Close(); err != nil {
			log.Panicf("Error closing client: %v", err)
		}
	}()

	database, err := db.NewDBWithDSN(ctx, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")))
	if err != nil {
		log.Fatal(err)
	}
	defer database.GetPool(ctx).Close()

	/*  grpc */
	addr := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT_GRPC")
	grpcServer := grpc.NewServer()
	serverImplement := &server.ServerImplement{
		ArticleRepo: postgresql.NewArticles(database),
		CommentRepo: postgresql.NewComments(database),
		Sender:      server.NewKafkaSender(kafkaProducer, "requests"),
	}
	pb.RegisterDbHandlerServer(grpcServer, serverImplement)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	mux := runtime.NewServeMux()
	errRest := pb.RegisterDbHandlerHandlerFromEndpoint(ctx, mux, addr, []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if errRest != nil {
		panic(err)
	}
	addrGateway := os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT_GATEWAY")

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		logger.Infof(ctx, "GRPC server listening on %q", addr)
		if err := grpcServer.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		reader.StartReader(ctx, clientConsumer)
	}()

	go func() {
		defer wg.Done()
		logger.Infof(ctx, "Gateway server listening at on %q", addrGateway)
		if err := http.ListenAndServe(addrGateway, mux); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

}
