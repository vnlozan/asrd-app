package main

import (
	"context"
	"errors"
	"listener/internal/config"
	"listener/internal/infra/loggerclient"
	"listener/internal/infra/rabbitmq"
	"listener/internal/service/event"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg := config.NewConfig()

	// 1. Dial + defer закрытие
	conn, err := rabbitmq.Dial(cfg.RabbitMQConfig.ConnectionURL)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	// 2. Consumer ←→ handler
	cons := rabbitmq.NewConsumer(conn)
	if err := cons.SetupChannel(); err != nil {
		log.Fatalf("declare exchange: %v", err)
	}

	logger := loggerclient.New(cfg.LoggerConfig.ConnectionURL)
	h := event.NewLogHandler(logger)

	// 3. Контекст, который закроется по Ctrl-C / docker stop
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Запускаем слушателя в отдельной горутине
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := cons.Listen(ctx,
			[]string{"log.INFO", "log.WARNING", "log.ERROR"},
			h.Handle,
		); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("listener error: %v", err)
		}
	}()

	// 5. Ждём сигнала
	<-ctx.Done()
	log.Println("shutdown: context canceled")

	// 6. Дожидаемся, пока consumer завершится
	wg.Wait()
	log.Println("graceful shutdown complete")
}
