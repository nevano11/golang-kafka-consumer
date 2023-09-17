package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang-kafka/internal/handler"
	"golang-kafka/internal/repository"
	"golang-kafka/internal/service"
	"net/http"
	"time"
)

// @title           Kafka producer
// @version         1.0
// @description     fio sender

// @host      localhost:8081
// @BasePath  /
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.Info("Starting kafka consumer")

	// Read config
	if err := initConfig(); err != nil {
		logrus.Fatalf("Failed to read config: %s", err.Error())
	}

	// Configure logger
	if err := configureLogger(viper.GetString("logger.log-level")); err != nil {
		logrus.Fatalf("Failed to configure logger: %s", err.Error())
	}

	// Database
	db, err := repository.NewPostgresDb(readDbConfig())
	if err != nil {
		logrus.Fatalf("Failed to create db connection: %s", err.Error())
	}

	// Repository
	rep := repository.NewRepository(db)

	// KafkaService
	kafkaService, err := service.NewKafkaService(
		viper.GetString("kafka.topic"),
		viper.GetString("kafka.topic-on-fail"),
		viper.GetString("kafka.config-path"),
		rep)
	if err != nil {
		logrus.Fatalf("Failed to create service: %s", err.Error())
	}
	defer kafkaService.Shutdown()

	// handler
	han := handler.NewHandler(kafkaService)

	// Routes
	routes := han.InitRoutes()

	// Server
	server := createServer(viper.GetString("server.port"), routes)
	logrus.Infof("Server running on http://localhost%s", server.Addr)
	logrus.Infof("Swagger: http://localhost%s/swagger/index.html", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		logrus.Fatalf("Failed to start server: %s", err.Error())
	}
}

// Configuration
func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func configureLogger(logLevel string) error {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	return nil
}

func createServer(port string, routes *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              ":" + port,
		Handler:           routes,
		ReadHeaderTimeout: 2 << 20,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
}

// Database
func readDbConfig() repository.DbConfig {
	return repository.DbConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		DbName:   viper.GetString("database.dbname"),
		SslMode:  viper.GetString("database.sslmode"),
	}
}
