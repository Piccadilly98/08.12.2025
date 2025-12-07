package midddleware

import (
	"fmt"
	"net/http"

	"github.com/Piccadilly98/linksChecker/internal/handlers"
	processing_os_signal "github.com/Piccadilly98/linksChecker/internal/processing_os_signal"
)

func MidddlewareCounterRequests(signalWorker *processing_os_signal.WorkerOSSignal) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if signalWorker.IsStoped() {
				w.Header().Set("Content-Type", "application/json")
				handlers.ProcessingError(w, r, fmt.Errorf("server stoped, please repeat you request later"), nil, http.StatusServiceUnavailable)
				return
			}
			signalWorker.AddRequest()
			defer signalWorker.DoneRequest()
			next.ServeHTTP(w, r)
		})
	}
}
