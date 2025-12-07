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

type TestCaseGracefulStop struct {
	Name              string
	EndPoint          string
	Method            string
	StopBeforeRequest bool
	StopInRequestTime bool
	ExpectedCode      int
	Body              any
}

func TestGracefulStop(t *testing.T) {
	server := server.MakeServer(10)
	err := InitStorage(server)
	if err != nil {
		t.Fatalf("Init Storage error - %v", err)
	}

	testCases := []TestCaseGracefulStop{

		//STOP BEFORE REQUEST
		{
			Name:              "stop_before_request_dock_query",
			EndPoint:          "/dock/query?bucketID=1",
			Method:            http.MethodGet,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
		},
		{
			Name:              "stop_before_request_registration_no_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
		},
		{
			Name:              "stop_before_request_registration_with_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
			Body: dto.RegistrationLinksRequest{
				Link: gePtr("http://youtube.com"),
			},
		},
		{
			Name:              "stop_before_request_dock_no_body",
			EndPoint:          "/dock",
			Method:            http.MethodGet,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
		},
		{
			Name:              "stop_before_request_dock_with_body",
			EndPoint:          "/dock",
			Method:            http.MethodGet,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2},
			},
		},
		{
			Name:              "stop_before_long_time_request_registration_with_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopBeforeRequest: true,
			ExpectedCode:      http.StatusServiceUnavailable,
			Body: dto.RegistrationLinksRequest{
				Link: gePtr("1111111111111111111111111111111111111"),
			},
		},

		//STOP AFTER REQUEST

		{
			Name:              "stop_after_request_dock_query",
			EndPoint:          "/dock/query?bucketID=1,2,3,4",
			Method:            http.MethodGet,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusOK,
		},
		{
			Name:              "stop_after_request_registration_no_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusBadRequest,
		},
		{
			Name:              "stop_after_request_registration_with_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusCreated,
			Body: dto.RegistrationLinksRequest{
				Link: gePtr("http://youtube.com"),
			},
		},
		{
			Name:              "stop_after_request_dock_no_body",
			EndPoint:          "/dock",
			Method:            http.MethodGet,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusBadRequest,
		},
		{
			Name:              "stop_after_request_dock_with_body",
			EndPoint:          "/dock",
			Method:            http.MethodGet,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusOK,
			Body: dto.GetBucketsRequest{
				LinksList: []int64{1, 2},
			},
		},
		{
			Name:              "stop_in_long_time_request_registration_with_body",
			EndPoint:          "/registration",
			Method:            http.MethodPost,
			StopInRequestTime: true,
			ExpectedCode:      http.StatusCreated,
			Body: dto.RegistrationLinksRequest{
				Link: gePtr("1111111111111111111111111111111111111"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.StopBeforeRequest {
				st := server.PauseUnpauseServerTesting()
				if st != true {
					t.Fatalf("error in stop server function")
				}
			}
			var req *http.Request
			w := httptest.NewRecorder()
			if tc.Body == nil {
				req = httptest.NewRequest(tc.Method, URL+tc.EndPoint, nil)
			} else {
				b, err := json.Marshal(tc.Body)
				if err != nil {
					t.Fatalf("error json: %v\n", err)
				}
				req = httptest.NewRequest(tc.Method, URL+tc.EndPoint, bytes.NewBuffer(b))
			}
			if tc.StopInRequestTime {
				started := make(chan bool)
				completed := make(chan bool)
				go func(rec *httptest.ResponseRecorder, r *http.Request) {
					started <- true
					server.R.ServeHTTP(rec, r)
					completed <- true
				}(w, req)
				<-started
				server.PauseUnpauseServerTesting()
				<-completed
			} else {
				server.R.ServeHTTP(w, req)
			}
			if w.Code != tc.ExpectedCode {
				t.Errorf("CODE ERROR: got: %d, expect : %d\n", w.Code, tc.ExpectedCode)
			}

			st := server.PauseUnpauseServerTesting()
			if st != false {
				t.Fatalf("error in stop server function")
			}
		})
	}
}
