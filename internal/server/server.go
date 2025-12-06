package server

import (
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/Piccadilly98/linksChecker/internal/handlers"
	linkchecker "github.com/Piccadilly98/linksChecker/internal/linkChecker"
	"github.com/Piccadilly98/linksChecker/internal/midddleware"
	processing_os_signal "github.com/Piccadilly98/linksChecker/internal/processingOSsignal"
	"github.com/Piccadilly98/linksChecker/internal/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	St           *storage.Storage
	Lp           *linkchecker.LinkProcessor
	R            *chi.Mux
	Reg          *handlers.RegistrationHandler
	GetBody      *handlers.GetBucketInfoHandler
	GetQuery     *handlers.GetBucketsInfoQueryHandler
	SignalWorker *processing_os_signal.WorkerOSSignal
}

func MakeServer(maxGourutine int) *Server {
	server := &Server{
		St:           storage.MakeStorage(),
		Lp:           linkchecker.MakeLinkProcessor(maxGourutine),
		R:            chi.NewRouter(),
		SignalWorker: processing_os_signal.MakeOSSignalWorker(),
	}
	server.Reg = handlers.MakeRegistrationHandler(server.St, server.Lp)
	server.GetBody = handlers.MakeGetBucketInfoHandler(server.St)
	server.GetQuery = handlers.MakeGetBucketInfoQueryHandler(server.St)
	server.R.Use(midddleware.MidddlewareCounterRequests(server.SignalWorker))
	server.R.Get("/dock/query", server.GetQuery.Handler)
	server.R.Get("/dock", server.GetBody.Handler)
	server.R.Post("/registration", server.Reg.Handler)
	return server
}

func (s *Server) Start(addresWithPort string) int {
	go func() {
		err := http.ListenAndServe(addresWithPort, s.R)
		if err != nil {
			panic(err)
		}
	}()
	time.Sleep(500 * time.Millisecond)
	log.Printf("server start in: %s\n", addresWithPort)
	s.SignalWorker.Start()
	log.Printf("server signal worker starting")
	return os.Getpid()
}

func (s *Server) Shutdown() {
	s.SignalWorker.Signals() <- syscall.SIGTERM
}
