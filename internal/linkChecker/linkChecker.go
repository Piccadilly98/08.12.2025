package linkchecker

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Piccadilly98/linksChecker/internal/storage"
)

const (
	Http  = "http://"
	Https = "https://"
)

type LinkProcessor struct {
	wg sync.WaitGroup
	ch chan struct{}
	mu sync.RWMutex
}

func MakeLinkProcessor(maxGoroutine int) *LinkProcessor {
	if maxGoroutine <= 0 {
		return nil
	}
	lp := &LinkProcessor{
		ch: make(chan struct{}, maxGoroutine),
	}
	return lp
}

func (lp *LinkProcessor) LinkChecker(links []string) map[string]string {
	res := make(map[string]string)
	for _, l := range links {
		lp.mu.RLock()
		_, ok := res[l]
		lp.mu.RUnlock()
		if ok {
			continue
		}
		lp.wg.Add(1)
		lp.ch <- struct{}{}
		go func(link string) {
			defer lp.wg.Done()
			defer func() {
				<-lp.ch
			}()
			status, _ := Processinglink(l)
			lp.mu.Lock()
			res[l] = status
			lp.mu.Unlock()
		}(l)
	}
	lp.wg.Wait()
	return res
}

// Func Processinglink(link string) string
// Func add befor url protocol http/https else needed and check real url-not redirect.
// The second line in the return values ​​is a test value to check the method.
func Processinglink(link string) (string, string) {
	needProtocol := false
	method := ""
	if !strings.HasPrefix(link, Http) && !strings.HasPrefix(link, Https) {
		needProtocol = true
	}
	if needProtocol {
		link = Https + link
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Head(link)
	if err != nil {
		status := getRequest(link, client)
		if status == storage.StatusNotAvalible && needProtocol {
			link = Http + link
			resp, err = client.Head(link)
			if err != nil {
				return getRequest(link, client), http.MethodGet
			}
			return processingCode(resp), http.MethodHead
		}
		method = http.MethodGet
	} else {
		method = http.MethodHead
	}
	status := processingCode(resp)
	return status, method
}

func getRequest(link string, client *http.Client) string {
	resp, err := client.Get(link)
	if err != nil {
		return storage.StatusNotAvalible
	}
	return processingCode(resp)
}

func processingCode(resp *http.Response) string {
	if resp.StatusCode < 400 {
		return storage.StatusAvalible
	}
	return storage.StatusNotAvalible
}
