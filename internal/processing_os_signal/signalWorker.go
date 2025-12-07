package processing_os_signal

import (
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type WorkerOSSignal struct {
	isStoped     atomic.Bool
	signals      chan os.Signal
	exitChan     chan struct{}
	countRequest atomic.Int64
	isOff        atomic.Bool
}

func MakeOSSignalWorker() *WorkerOSSignal {
	res := &WorkerOSSignal{}
	res.isStoped.Store(false)
	res.countRequest.Store(0)
	res.isOff.Store(false)
	res.signals = make(chan os.Signal)
	res.exitChan = make(chan struct{})
	signal.Notify(res.signals, syscall.SIGINT, syscall.SIGTERM)
	return res
}

func (w *WorkerOSSignal) ExitChan() chan struct{} {
	return w.exitChan
}
func (w *WorkerOSSignal) Signals() chan os.Signal {
	return w.signals
}
func (w *WorkerOSSignal) IsStoped() bool {
	return w.isStoped.Load()
}

func (w *WorkerOSSignal) Start() {
	go func() {
		for s := range w.signals {
			w.isStoped.Store(!w.IsStoped())
			log.Printf("Server stop status: %v\n", w.IsStoped())
			if s == syscall.SIGTERM {
				defer w.isOff.Store(true)
				if w.countRequest.Load() == 0 {
					w.exitChan <- struct{}{}
					return
				}
				time.Sleep(3 * time.Second)
				if w.countRequest.Load() == 0 {
					w.exitChan <- struct{}{}
					return
				}
				log.Printf("Lost requests: %d\n", w.countRequest.Load())
				w.exitChan <- struct{}{}
				return
			}
		}
	}()
}

func (w *WorkerOSSignal) AddRequest() {
	w.countRequest.Add(1)
}

func (w *WorkerOSSignal) DoneRequest() {
	w.countRequest.Add(-1)
}

func (w *WorkerOSSignal) IsOff() bool {
	return w.isOff.Load()
}

func (w *WorkerOSSignal) PauseUnpauseServerTesting() bool {
	w.isStoped.Store(!w.IsStoped())
	return w.IsStoped()
}
