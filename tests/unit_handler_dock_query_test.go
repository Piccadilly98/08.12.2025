package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/Piccadilly98/linksChecker/internal/DTO"
	"github.com/Piccadilly98/linksChecker/internal/server"
)

const (
	HeaderName       = "Content-Type"
	HeaderFormatPDF  = "application/pdf"
	HeaderFormatJSON = "application/json"
	URLQuery         = "http://localhost:8080/dock/query"
	BeginQuery       = "?bucketID="
)

type TestCaseDockQueryHandler struct {
	Name           string
	QueryParam     string
	Method         string
	ExpectedFormat string
	ExpectedCode   int
}

func InitStorage(server *server.Server) error {
	method := http.MethodPost
	url := "http://localhost:8080/registration"

	regCase := []dto.RegistrationLinksRequest{
		{
			Links: []string{"vk.com", "yandex.ru", "ya.ru", "mail.ru", "google.com"},
			Link:  gePtr("github.com"),
		},
		{
			Links: []string{"https://www.youtube.com", "https://docs.google.com/document", "https://cs50.com"},
		},
		{
			Link: gePtr("gmail.com"),
		},
		{
			Links: []string{"youtube.com", "https://music.yandex.ru/"},
			Link:  gePtr("https://www.wikipedia.org/"),
		},
	}

	for _, rc := range regCase {
		b, err := json.Marshal(rc)
		if err != nil {
			return err
		}

		req := httptest.NewRequest(method, url, bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		server.R.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			return fmt.Errorf("error!")
		}
	}
	return nil
}

func TestDockQueryHandler(t *testing.T) {
	server := server.MakeServer(10)
	err := InitStorage(server)
	if err != nil {
		t.Error(err)
	}

	testCase := []TestCaseDockQueryHandler{

		//VALID

		{
			Name:           "valid_1_param",
			QueryParam:     BeginQuery + "1",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_2_repeat_param",
			QueryParam:     BeginQuery + "1,1",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_3_repeat_param",
			QueryParam:     BeginQuery + "1,1,1",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_2_in_a_row_param",
			QueryParam:     BeginQuery + "1,2",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_3_in_a_row_param",
			QueryParam:     BeginQuery + "1,2,3",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_4_in_a_row_param",
			QueryParam:     BeginQuery + "1,2,3,4",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_4_param_random",
			QueryParam:     BeginQuery + "2,3,1,4",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_1_param_format\"1-2\"",
			QueryParam:     BeginQuery + "1-1",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_2_param_format\"1-2\"",
			QueryParam:     BeginQuery + "1-2",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_3_param_format\"1-2\"",
			QueryParam:     BeginQuery + "1-3",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_4_param_format\"1-2\"",
			QueryParam:     BeginQuery + "1-4",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},
		{
			Name:           "valid_2_random_param_format\"1-2\"",
			QueryParam:     BeginQuery + "2-3",
			ExpectedFormat: HeaderFormatPDF,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusOK,
		},

		//INVALID METHODS

		{
			Name:         "invalid_method_post",
			QueryParam:   BeginQuery + "2-3",
			Method:       http.MethodPost,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_head",
			QueryParam:   BeginQuery + "2-3",
			Method:       http.MethodHead,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_patch",
			QueryParam:   BeginQuery + "2-3",
			Method:       http.MethodPatch,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_put",
			QueryParam:   BeginQuery + "2-3",
			Method:       http.MethodPut,
			ExpectedCode: http.StatusMethodNotAllowed,
		},
		{
			Name:         "invalid_method_delete",
			QueryParam:   BeginQuery + "2-3",
			Method:       http.MethodDelete,
			ExpectedCode: http.StatusMethodNotAllowed,
		},

		// INVALID QUERY

		{
			Name:           "invalid_query_param_empty",
			QueryParam:     BeginQuery + "",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_1_query_param",
			QueryParam:     BeginQuery + "-",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_1_query_param_<0",
			QueryParam:     BeginQuery + "-1",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_1_query_param_invalid_bucket",
			QueryParam:     BeginQuery + "100",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_2_query_param_valid_and_non_valid",
			QueryParam:     BeginQuery + "1,100",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_2_query_param_invalid_bucket",
			QueryParam:     BeginQuery + "100,12",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"_begin>finish",
			QueryParam:     BeginQuery + "2-1",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"_begin<0",
			QueryParam:     BeginQuery + "-1-1",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"",
			QueryParam:     BeginQuery + "2-3-4",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"_invalid_bucket_finish",
			QueryParam:     BeginQuery + "1-5",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"_invalid_bucket_begin",
			QueryParam:     BeginQuery + "5-5",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
		{
			Name:           "invalid_format_query_param_ \"1-2\"_ultra_big_finish",
			QueryParam:     BeginQuery + "1-1000",
			ExpectedFormat: HeaderFormatJSON,
			Method:         http.MethodGet,
			ExpectedCode:   http.StatusBadRequest,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.Method, URLQuery+tc.QueryParam, http.NoBody)

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
