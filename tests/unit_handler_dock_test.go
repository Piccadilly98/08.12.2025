package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
	"github.com/Piccadilly98/linksChecker/internal/server"
)

type TestCaseDockHandler struct {
	Name           string
	Body           dto.GetBucketsRequest
	Method         string
	ExpectedFormat string
	ExpectedCode   int
}

const (
	targetDockURL = "http://localhost:8080/dock"
)

func TestDockHandler(t *testing.T) {
	server := server.MakeServer(10)
	err := InitStorage(server)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []TestCaseDockHandler{

		//VALID
		{
			Name: "valid_1_number_in_link",
			Body: dto.GetBucketsRequest{
				LinkList: getPtrInt(1),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_1_number_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_2_numbers_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_2_numbers_in_links_1_in_link",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2},
				LinkList:  getPtrInt(3),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_3_numbers_in_links_1_in_link",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2, 3},
				LinkList:  getPtrInt(4),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_4_numbers_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2, 3, 4},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_4_repeat_numbers_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 1, 1, 1},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name: "valid_2_repeat_numbers_in_links_and_1_in_link",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 1},
				LinkList:  getPtrInt(1),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatPDF,
			ExpectedCode:   http.StatusOK,
		},

		//INVALID METHOD

		{
			Name:         "invalid_method_post",
			Method:       http.MethodPost,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_head",
			Method:       http.MethodHead,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_patch",
			Method:       http.MethodPatch,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_put",
			Method:       http.MethodPut,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_delete",
			Method:       http.MethodDelete,
			ExpectedCode: http.StatusMethodNotAllowed,
		},

		//INVALID

		{
			Name: "invalid_1_number_in_link",
			Body: dto.GetBucketsRequest{
				LinkList: getPtrInt(5),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name: "invalid_1_number_in_link_<0",
			Body: dto.GetBucketsRequest{
				LinkList: getPtrInt(-1),
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name: "invalid_1_number_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{5},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name: "invalid_1_number_in_links_<0",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{-1},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name: "invalid_2_number_in_links",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{5, 6},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name: "invalid_1_number_in_links_and_1_valid",
			Body: dto.GetBucketsRequest{
				LinksList: []int64{5, 1},
			},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_empty_body",
			Body:           dto.GetBucketsRequest{},
			Method:         http.MethodGet,
			ExpectedFormat: HeaderFormatJSON,
			ExpectedCode:   http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			b, err := json.Marshal(tc.Body)
			if err != nil {
				t.Error("Json err")
			}
			req := httptest.NewRequest(tc.Method, targetDockURL, bytes.NewBuffer(b))
			server.R.ServeHTTP(w, req)
			if w.Code != tc.ExpectedCode {
				t.Errorf("CODE: got: %d, expect: %d\n", w.Code, tc.ExpectedCode)
			}
			if w.Header().Get(HeaderName) != tc.ExpectedFormat {
				t.Errorf("FORMAT: got:%s, expect: %s\n", w.Header().Get(HeaderName), tc.ExpectedFormat)
			}
			if tc.ExpectedFormat == HeaderFormatPDF {
				body := w.Body.Bytes()
				if len(body) == 0 {
					t.Error("PDF empty")
				}
				if len(body) < 5 || string(body[:4]) != "%PDF" {
					t.Error("Its not pdf!")
				}
				if !bytes.Contains(body, []byte("%%EOF")) {
					t.Error("PDF missing EOF marker")
				}
			}
		})
	}
}

func getPtrInt(i int) *int64 {
	res := int64(i)
	return &res
}
