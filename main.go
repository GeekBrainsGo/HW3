package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	signals := make(chan os.Signal)
	lg := configureLogger()
	lg.Info("Запуск сервера ...")

	serv := NewServer(lg)

	go func() {
		err := serv.Start()
		if err != nil {
			lg.WithError(err).Fatal("Запуск сервера не возможен")
		}
	}()

	signal.Notify(signals, os.Kill, os.Interrupt)
	<-signals
	lg.Info("Сервер остановлен!")
}

// configureLogger - Настраивает логгер
func configureLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)
	return lg
}
