package linkchecker

import (
	"net"
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

const (
	defaultMaxGoroutine = 50
)

func MakeLinkProcessor(maxGoroutine int) *LinkProcessor {
	if maxGoroutine <= 0 {
		maxGoroutine = defaultMaxGoroutine
	}
	lp := &LinkProcessor{
		ch: make(chan struct{}, maxGoroutine),
	}
	return lp
}

func (lp *LinkProcessor) LinkChecker(links []string) map[string]string {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

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
		go func(link string, cl *http.Client) {
			defer lp.wg.Done()
			defer func() {
				<-lp.ch
			}()
			status, _, _ := Processinglink(link, cl)
			lp.mu.Lock()
			res[link] = status
			lp.mu.Unlock()
		}(l, client)
	}
	lp.wg.Wait()
	return res
}

func Processinglink(link string, client *http.Client) (string, string, string) {
	needProtocol := false
	method := ""
	proto := ""
	if !strings.HasPrefix(link, Http) && !strings.HasPrefix(link, Https) {
		needProtocol = true
	}
	if needProtocol {
		link = Https + link
		proto = Https
	}

	resp, err := client.Head(link)
	if err != nil {
		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return storage.StatusNotAvalible, http.MethodHead, proto
			}
		}
		status := getRequest(link, client)
		if status == storage.StatusNotAvalible && needProtocol {
			link = Http + link
			proto = Http
			resp, err = client.Head(link)
			if err != nil {
				return getRequest(link, client), http.MethodGet, proto
			}
			return processingCode(resp), http.MethodHead, proto
		}
		method = http.MethodGet
	} else {
		method = http.MethodHead
	}
	status := processingCode(resp)
	return status, method, proto
}

func getRequest(link string, client *http.Client) string {
	resp, err := client.Get(link)
	if err != nil {
		return storage.StatusNotAvalible
	}
	return processingCode(resp)
}

func processingCode(resp *http.Response) string {
	if resp == nil {
		return storage.StatusNotAvalible
	}
	if resp.StatusCode >= 500 {
		return storage.StatusNotAvalible
	}
	return storage.StatusAvalible
}
