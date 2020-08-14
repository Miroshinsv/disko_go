package main

import (
	"context"
	web_server "github.com/Miroshinsv/disko_go/internal/web-server"
	configService "github.com/Miroshinsv/disko_go/pkg/config-service"
	db_connector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	webServer "github.com/Miroshinsv/disko_go/pkg/web-server"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	conf, log := bootstrap()

	webConf := webServer.Config{}
	err := conf.Convert(&webConf)
	if err != nil {
		log.Fatal("Error in web config", err, nil)
	}
	webConf.Port, err = strconv.Atoi(os.Getenv("PORT"))
	webConf.Host = ""

	if err != nil {
		log.Fatal("Port no defined", err, nil)
	}

	web := webServer.MustNewWebServer(&webConf, log)
	web_server.RegisterHandlers()
	web.RegisterRoutes(web_server.WebRouter)
	go func() {
		err := web.ListenAndServe(context.Background())
		if err != nil {
			log.Fatal("Error on serving web", err, nil)
		}
	}()

	registerShutdown(context.Background(), log, web)
}

func bootstrap() (configService.IConfig, loggerService.ILogger) {
	config := configService.GetConfigService()
	err := config.Load(configService.DefaultConfigPath)
	if err != nil {
		panic(err)
	}

	logConfig := loggerService.Config{}
	err = config.Convert(&logConfig)
	if err != nil {
		panic(err)
	}

	log := loggerService.MustNewLogger(logConfig)

	dbConf := db_connector.Config{}
	err = config.Convert(&dbConf)
	if err != nil {
		panic(err)
	}

	db := db_connector.MustNewDBConnection(log, &dbConf)
	err = db.Connect()
	if err != nil {
		log.Fatal("error on connecting to db", err, nil)
	}

	return config, log
}

func registerShutdown(ctx context.Context, log loggerService.ILogger, web webServer.IWebServer) {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)

	for {
		select {
		case <-ch:
			log.Warning("Interrupt signal fetched", nil)
			err := web.Stop(ctx)
			if err != nil {
				log.Error("Very sad to see error on shutdown", err, nil)
			}

			return
		}
	}
}
