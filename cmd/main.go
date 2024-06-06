package main

import (
	"context"
	"github.com/azicussdu/go-todo-app"
	"github.com/azicussdu/go-todo-app/pkg/handler"
	"github.com/azicussdu/go-todo-app/pkg/repository"
	"github.com/azicussdu/go-todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

//for env files use better libraries: for example: github.com/kelseyhightower/envconfig

func main() {
	// error formatting for logrus lib (will show errors as json object)
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	// using godotenv lib reading from .env file (should store password and sensitive info)
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.user"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.ssl_mode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	// running server on a new goroutine
	go func() {
		if err = srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)

	// SIGINT is the signal sent when we press Ctrl+C (termination of process). (in linux systems)
	// The SIGTERM and SIGQUIT signals are meant to terminate the process.
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit // when we got signal we just read from channel so program(main function) finishes

	logrus.Print("TodoApp Shutting Down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err = db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

// using viper lib getting connection with configs/config.yml file
func initConfig() error {
	viper.AddConfigPath("configs") //folder name
	viper.SetConfigName("config")  //config file name
	return viper.ReadInConfig()
}
