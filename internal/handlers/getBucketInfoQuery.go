package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	document_worker "github.com/Piccadilly98/linksChecker/internal/documentWorker"
	"github.com/Piccadilly98/linksChecker/internal/storage"
)

const (
	AtoiErr = "invalid syntax"
)

type GetBucketsInfoQueryHandler struct {
	st *storage.Storage
}

func MakeGetBucketInfoQueryHandler(st *storage.Storage) *GetBucketsInfoQueryHandler {
	return &GetBucketsInfoQueryHandler{
		st: st,
	}
}

func (g *GetBucketsInfoQueryHandler) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ids, err := g.processingQuery(r)
	if err != nil {
		if strings.Contains(err.Error(), AtoiErr) {
			err = fmt.Errorf("invalid bucket id")
		}
		ProcessingError(w, r, err, nil)
		return
	}
	res, err := g.st.GetBucketsInfo(ids...)
	if err != nil {
		ProcessingError(w, r, err, nil)
		return
	}
	b, err := document_worker.CreateDocument(res)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(b)
}

func (g *GetBucketsInfoQueryHandler) processingQuery(r *http.Request) ([]int64, error) {
	ids := r.URL.Query().Get("bucketID")
	if ids == "" {
		return nil, fmt.Errorf("invalid bucket id")
	}
	if strings.Contains(ids, "-") {
		idRange := strings.Split(ids, "-")
		beginNum, finishNum, err := g.processingFormatNums(idRange)
		if err != nil {
			return nil, err
		}

		slice := g.getSliceForRange(beginNum, finishNum)
		return slice, nil
	}
	idSlice := strings.Split(ids, ",")
	res := make([]int64, len(idSlice))
	for i, r := range idSlice {
		num, err := strconv.Atoi(r)
		if err != nil {
			return nil, err
		}
		if num < 0 {
			return nil, fmt.Errorf("invalid bucketID - %d", num)
		}
		res[i] = int64(num)
	}
	return res, nil
}

func (g *GetBucketsInfoQueryHandler) getSliceForRange(startNum, finishNum int) []int64 {
	res := []int64{}
	for i := startNum; i <= finishNum; i++ {
		res = append(res, int64(i))
	}
	return res
}

func (g *GetBucketsInfoQueryHandler) processingFormatNums(idRange []string) (int, int, error) {
	if len(idRange) != 2 {
		return 0, 0, fmt.Errorf("invalid format query")
	}
	beginNum, err := strconv.Atoi(idRange[0])
	if err != nil {
		return 0, 0, err
	}
	finishNum, err := strconv.Atoi(idRange[1])
	if err != nil {
		return 0, 0, err
	}
	if beginNum > finishNum {
		return 0, 0, fmt.Errorf("invalid range")
	}
	if beginNum < 0 || finishNum < 0 {
		return 0, 0, fmt.Errorf("invalid bucketID")
	}
	return beginNum, finishNum, nil
}
