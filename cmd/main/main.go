package main

import (
	"net/http"

	"github.com/Piccadilly98/linksChecker/internal/handlers"
	linkchecker "github.com/Piccadilly98/linksChecker/internal/linkChecker"
	"github.com/Piccadilly98/linksChecker/internal/storage"
	"github.com/go-chi/chi/v5"
)

func main() {
	// m := make(map[int64]map[string]string)
	// m[1] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[2] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[3] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[4] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[5] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[6] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[7] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// m[8] = map[string]string{
	// 	"youtube.com": storage.StatusAvalible,
	// 	"vk.com":      storage.StatusNotAvalible,
	// 	"ok.ru":       storage.StatusAvalible,
	// 	"grok.com":    storage.StatusNotAvalible,
	// 	"f.gg":        storage.StatusNotAvalible,
	// }
	// document_worker.CreateDocument(m)
	// return
	r := chi.NewRouter()
	st := storage.MakeStorage()
	lp := linkchecker.MakeLinkProcessor(50)
	if lp == nil {
		panic("invalid max goroutine num")
	}
	reg := handlers.MakeRegistrationHandler(st, lp)
	if reg == nil {
		panic("reg is null")
	}

	get := handlers.MakeGetBucketInfoHandler(st)
	r.Get("/dock", get.Hadler)
	r.Post("/registration", reg.Hadler)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		panic(err)
	}
}
