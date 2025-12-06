package midddleware

import (
	"net/http"

	processing_os_signal "github.com/Piccadilly98/linksChecker/internal/processingOSsignal"
)

func MidddlewareCounterRequests(signalWorker *processing_os_signal.WorkerOSSignal) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if signalWorker.IsStoped() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			signalWorker.AddRequest()
			defer signalWorker.DoneRequest()
			next.ServeHTTP(w, r)
		})
	}
}
