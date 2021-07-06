package main

import (
	"context"
	poll_service "github.com/Miroshinsv/disko_go/internal/poll-service"
	schedule_service "github.com/Miroshinsv/disko_go/internal/schedule-service"
	web_server "github.com/Miroshinsv/disko_go/internal/web-server"
	configService "github.com/Miroshinsv/disko_go/pkg/config-service"
	db_connector "github.com/Miroshinsv/disko_go/pkg/db-connector"
	loggerService "github.com/Miroshinsv/disko_go/pkg/logger-service"
	webServer "github.com/Miroshinsv/disko_go/pkg/web-server"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	conf, log := bootstrap()

	webConf := webServer.Config{}
	err := conf.Convert(&webConf)
	if err != nil {
		log.Fatal("Error in web config", err, nil)
	}

	if os.Getenv("PORT") != "" {
		webConf.Port, err = strconv.Atoi(os.Getenv("PORT"))
		webConf.Host = ""
	}

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

	registerAutoPolls(log)
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

func registerAutoPolls(log loggerService.ILogger) {
	ticker := time.NewTicker(5 * time.Minute) //@todo: move to env
	log.Info("Start autopolls", nil)
	go func() {
		n := time.Now()
		tm := time.Now().Add(24 * time.Hour)

		for {
			select {
			case <-ticker.C:
				events, err := schedule_service.GetScheduleService().LoadEventsForPeriod(n, tm)
				if err != nil {
					log.Error("error on scheduling polls", err, nil)

					return
				}

				pS := poll_service.GetPollService()
				err = pS.ScheduleAutoPolls(events[n.Format("2006-01-02")], n)
				if err != nil {
					log.Error("error on scheduling polls", err, nil)

					return
				}

				err = pS.ScheduleAutoPolls(events[tm.Format("2006-01-02")], tm)
				if err != nil {
					log.Error("error on scheduling polls", err, nil)

					return
				}
				log.Info("scheduling polls", nil)
			}
		}
	}()
}
