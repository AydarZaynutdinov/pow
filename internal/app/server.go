package app

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/logger"
	"github.com/AydarZaynutdinov/pow/internal/service"
)

type QuoteRepository interface {
	GetQuote() []byte
}

type PoWService interface {
	GenerateChallenge() []byte
	VerifySolution([]byte, []byte) bool
}

type PoWServer struct {
	ctx             context.Context
	cfg             *config.ServerConfig
	connChan        chan net.Conn
	wg              *sync.WaitGroup
	quoteRepository QuoteRepository
	powService      PoWService
}

func NewServer(ctx context.Context, cfg *config.ServerConfig, quoteRepository QuoteRepository, powService PoWService) *PoWServer {
	logger.SetLogger(cfg.Server.LogLevel)

	return &PoWServer{
		ctx:             ctx,
		cfg:             cfg,
		connChan:        make(chan net.Conn, 1000),
		wg:              &sync.WaitGroup{},
		quoteRepository: quoteRepository,
		powService:      powService,
	}
}

func (s *PoWServer) Start() error {
	listener, err := net.Listen("tcp", s.cfg.Server.Address)
	if err != nil {
		slog.Error("Error listening on port %s: %v", s.cfg.Server.Address, err)
		return err
	}
	defer func() {
		_ = listener.Close()
	}()
	slog.Info("Listening on port 8080")

	// start workers
	for i := 0; i < s.cfg.PoW.WorkersCount; i++ {
		s.wg.Add(1)
		go s.workerRun()
	}

	for {
		select {
		case <-s.ctx.Done():
			slog.Info("Server context is canceled. Stop accepting new connections")
			return nil
		default:
			conn, connErr := listener.Accept()
			if connErr != nil {
				slog.Warn("Error accepting connection: %v", err)
			}
			s.connChan <- conn
		}
	}
}

func (s *PoWServer) Shutdown() {
	slog.Info("Shutting down PoW server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.ShutDownTimeout)
	defer cancel()

	// close channel with new connections
	close(s.connChan)

	done := make(chan struct{})
	defer close(done)
	go func() {
		s.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		slog.Info("All workers shutdown complete")
	case <-shutdownCtx.Done():
		slog.Info("Server shutdown timeout")
	}
}

func (s *PoWServer) workerRun() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case conn, ok := <-s.connChan:
			if !ok {
				slog.Debug("Connection channel closed")
				return
			}
			s.handleConn(conn)
		}
	}
}

func (s *PoWServer) handleConn(conn net.Conn) {
	defer func() {
		_ = conn.Close()
	}()

	challenge := s.powService.GenerateChallenge()
	slog.Debug("challenge:", challenge)
	if _, err := conn.Write(challenge); err != nil {
		slog.Error("Error sending challenge: %v", err)
		return
	}

	slog.Debug("complexity:", s.cfg.PoW.Complexity)
	if _, err := conn.Write([]byte{s.cfg.PoW.Complexity}); err != nil {
		slog.Error("Error sending complexity: %v", err)
		return
	}

	done := make(chan struct{})
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		solution := make([]byte, service.SolutionLen)
		if _, err := conn.Read(solution); err != nil {
			slog.Error("Error reading solution: %v", err)
			return
		}
		slog.Debug("solution:", solution)

		if !s.powService.VerifySolution(challenge, solution) {
			slog.Debug("Solution verification failed")

			if _, err := conn.Write([]byte("Incorrect solution\n")); err != nil {
				slog.Error("Error sending solution check: %v", err)
				return
			}
			return
		}

		slog.Debug("Solution verified")
		quote := s.quoteRepository.GetQuote()
		if _, err := conn.Write(quote); err != nil {
			slog.Error("Error sending quote: %v", err)
			return
		}
	}()

	select {
	case <-done:
		slog.Debug("Connection handler done")
	case <-time.After(s.cfg.PoW.HandlerTimeout):
		slog.Debug("Handler timeout")
		if _, err := conn.Write([]byte("Handler timeout\n")); err != nil {
			slog.Error("Error sending timeout message: %v", err)
		}
	}
	slog.Debug("___")
}
