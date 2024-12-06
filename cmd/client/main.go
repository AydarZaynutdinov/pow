package main

import (
	"io"
	"log"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/logger"
	"github.com/AydarZaynutdinov/pow/internal/service"
)

func main() {
	clientCfg, err := config.ParseClient()
	if err != nil {
		log.Fatal("Failed to parse client config:", err)
	}

	logger.SetLogger(clientCfg.LogLevel)
	ticker := time.NewTicker(time.Second / time.Duration(clientCfg.RPS))

	wg := &sync.WaitGroup{}
	for i := 0; i < clientCfg.Total; i++ {
		<-ticker.C

		wg.Add(1)
		go func() {
			defer wg.Done()

			conn, err := net.Dial("tcp", clientCfg.Address)
			if err != nil {
				slog.Error("Failed to connect to server:", err)
			}

			challenge := make([]byte, service.ChallengeLen)
			if _, err := conn.Read(challenge); err != nil {
				slog.Error("Failed to read challenge from server:", err)
			}
			slog.Debug("challenge:", challenge)

			complexity := make([]byte, 1)
			if _, err := conn.Read(complexity); err != nil {
				slog.Error("Failed to read  complexity from server:", err)
			}
			slog.Debug("complexity:", complexity)

			pow := service.NewPoW(complexity[0])
			solution := pow.Solve(challenge)
			slog.Debug("solution:", solution)

			if _, err := conn.Write(solution); err != nil {
				slog.Error("Failed to write solution:", err)
			}

			response, err := io.ReadAll(conn)
			if err != nil {
				slog.Error("Failed to read response:", err)
			}

			slog.Debug("Response:", string(response))
			slog.Debug("___")
		}()
	}
}
