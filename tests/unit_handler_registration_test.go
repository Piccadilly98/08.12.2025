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
	URL     = "http://localhost:8080"
	Address = "localhost:8080"
)

type TestCaseRegistrationHandler struct {
	Name              string
	EndPoint          string
	Body              dto.RegistrationLinks
	Method            string
	ExpectedCode      int
	CheckResponse     bool
	ExpectedLinkInMap int
}

func TestRegistrationHandler(t *testing.T) {
	testCases := []TestCaseRegistrationHandler{

		//VALID
		{
			Name:     "valid_request_1_link",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Link: gePtr("vk.com"),
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 1,
		},
		{
			Name:     "valid_request_1_link_in_array",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Links: []string{"vk.com"},
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 1,
		},
		{
			Name:     "valid_request_2_link_in_array",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Links: []string{"vk.com", "ok.ru"},
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 2,
		},
		{
			Name:     "valid_request_more_identical_link_in_array",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Links: []string{"vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru",
					"vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru",
					"vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru", "vk.com", "ok.ru"},
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 2,
		},
		{
			Name:     "valid_request_more_unique_link_in_array",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Links: []string{"vk.com", "ok.ru", "https://youtube.com", "gg.g", "http://wik.gg"},
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 5,
		},
		{
			Name:     "valid_request_more_link_in_array_and_link",
			EndPoint: "/registration",
			Body: dto.RegistrationLinks{
				Links: []string{"vk.com", "ok.ru", "https://youtube.com", "gg.g", "http://wik.gg"},
				Link:  gePtr("s.com"),
			},
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 6,
		},
		{
			Name:     "valid_request_100_unique_links",
			EndPoint: "/registration",
			Body: func() dto.RegistrationLinks {
				links := make([]string, 100)
				for i := 0; i < 100; i++ {
					links[i] = fmt.Sprintf("example%d.com", i)
				}
				return dto.RegistrationLinks{Links: links}
			}(),
			Method:            http.MethodPost,
			ExpectedCode:      http.StatusCreated,
			CheckResponse:     true,
			ExpectedLinkInMap: 100,
		},

		//INVALID METHOD

		{
			Name:          "invalid_method_get",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodGet,
			ExpectedCode:  http.StatusMethodNotAllowed,
			CheckResponse: false,
		},
		{
			Name:          "invalid_method_head",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodHead,
			ExpectedCode:  http.StatusMethodNotAllowed,
			CheckResponse: false,
		},
		{
			Name:          "invalid_method_delete",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodDelete,
			ExpectedCode:  http.StatusMethodNotAllowed,
			CheckResponse: false,
		},
		{
			Name:          "invalid_method_patch",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodPatch,
			ExpectedCode:  http.StatusMethodNotAllowed,
			CheckResponse: false,
		},
		{
			Name:          "invalid_method_options",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodOptions,
			ExpectedCode:  http.StatusMethodNotAllowed,
			CheckResponse: false,
		},

		//INVALID BODY

		{
			Name:          "invalid_no_body",
			EndPoint:      "/registration",
			Method:        http.MethodPost,
			ExpectedCode:  http.StatusBadRequest,
			CheckResponse: false,
		},
		{
			Name:          "invalid_empty_array",
			EndPoint:      "/registration",
			Body:          dto.RegistrationLinks{},
			Method:        http.MethodPost,
			ExpectedCode:  http.StatusBadRequest,
			CheckResponse: false,
		},
	}
	lastBucketNum := 0
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			server := server.MakeServer(10)
			body, err := json.Marshal(tc.Body)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v\n", err)
			}
			req := httptest.NewRequest(tc.Method, URL+tc.EndPoint, bytes.NewBuffer(body))
			w := httptest.NewRecorder()
			server.R.ServeHTTP(w, req)
			if tc.ExpectedCode != w.Code {
				t.Errorf("CODE: got: %d, expect: %d\n", w.Code, tc.ExpectedCode)
			}
			if tc.CheckResponse {
				response := dto.GetInfoBucketDTO{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse JSON: %v\nResponse: %s", err, w.Body.String())
				}
				if response.NumBucket <= int64(lastBucketNum) {
					t.Errorf("invalid bucket id\ngot: %d, expect: %d\n", response.NumBucket, lastBucketNum+1)
				}
				if len(response.Links) != tc.ExpectedLinkInMap {
					t.Error("Error: not create bucket")
				}
			}
		})
	}
}

func gePtr(str string) *string {
	s := &str
	return s
}
