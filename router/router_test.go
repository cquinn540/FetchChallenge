package router

import (
	"FetchChallenge/store"
	"FetchChallenge/types"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router

func TestMain(m *testing.M) {
	// Create the router once before all tests
	router = CreateRouter()
	m.Run()
}

func TestValidRequests(t *testing.T) {
	tt := []struct {
		name           string
		expected       string
		expectedStatus int
		requestBody    string
	}{
		{
			name:           "Target",
			expected:       `{"points":"28"}`,
			expectedStatus: http.StatusCreated,
			requestBody: `
			{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{
						"shortDescription": "Mountain Dew 12PK",
						"price": "6.49"
					},{
						"shortDescription": "Emils Cheese Pizza",
						"price": "12.25"
					},{
						"shortDescription": "Knorr Creamy Chicken",
						"price": "1.26"
					},{
						"shortDescription": "Doritos Nacho Cheese",
						"price": "3.35"
					},{
						"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
						"price": "12.00"
					}
				],
				"total": "35.35"
			}
			`,
		},
		{
			name:     "M&M Corner Market",
			expected: `{"points":"109"}`,
			requestBody: `
			{
				"retailer": "M&M Corner Market",
				"purchaseDate": "2022-03-20",
				"purchaseTime": "14:33",
				"items": [
					{
						"shortDescription": "Gatorade",
						"price": "2.25"
					},{
						"shortDescription": "Gatorade",
						"price": "2.25"
					},{
						"shortDescription": "Gatorade",
						"price": "2.25"
					},{
						"shortDescription": "Gatorade",
						"price": "2.25"
					}
				],
				"total": "9.00"
			}
			`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new store before each test for test isolation
			store.CreateReceiptStore()

			postReceiptReq := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewBuffer([]byte(tc.requestBody)))
			postReceiptRecorder := httptest.NewRecorder()

			router.ServeHTTP(postReceiptRecorder, postReceiptReq)

			if postReceiptRecorder.Code != http.StatusCreated {
				t.Errorf("POST receipt expected status: '%d', actual: '%d'", http.StatusCreated, postReceiptRecorder.Code)
			}

			var idResponse types.ReceiptCreated
			if err := json.Unmarshal(postReceiptRecorder.Body.Bytes(), &idResponse); err != nil {
				t.Errorf("POST receipt json parse error: '%s'", err.Error())
			}

			getPointsUrl := "/receipts/" + idResponse.Id + "/points"
			getPointsReq := httptest.NewRequest(http.MethodGet, getPointsUrl, nil)
			getPointsRecorder := httptest.NewRecorder()

			router.ServeHTTP(getPointsRecorder, getPointsReq)

			if getPointsRecorder.Code != http.StatusOK {
				t.Errorf("GET points expected status: '%d', actual: '%d'", http.StatusOK, postReceiptRecorder.Code)
			}

			points := getPointsRecorder.Body.String()
			if points != tc.expected {
				t.Errorf("GET points expected status: '%s', actual: '%s'", tc.expected, points)
			}
		})
	}
}

func TestInvalidRequests(t *testing.T) {
	tt := []struct {
		name           string
		route          string
		method         string
		expectedStatus int
		requestBody    string
	}{
		{
			name:           "/ returns 404",
			route:          "/",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			requestBody:    "",
		},
		{
			name:           "/abc returns 404",
			route:          "/abc",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			requestBody:    "",
		},
		{
			name:           "GET /receipts/process returns 404",
			route:          "/receipts/process",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			requestBody:    "",
		},
		{
			name:           "POST /receipts/{id}/points returns 404",
			route:          "/receipts/9c5b94b1-35ad-49bb-b118-8e8fc24abf80/points",
			method:         http.MethodPost,
			expectedStatus: http.StatusNotFound,
			requestBody:    "",
		},
		{
			name:           "GET /receipts/{id}/points returns 404 when id is not in store",
			route:          "/receipts/9c5b94b1-35ad-49bb-b118-8e8fc24abf80/points",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
			requestBody:    "",
		},
		{
			name:           "POST /receipts/process returns 400 with empty body",
			route:          "/receipts/process",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
			requestBody:    "",
		},
		{
			name:           "POST /receipts/process returns 400 with invalid item price",
			route:          "/receipts/process",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
			requestBody: `
			{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{
						"shortDescription": "Mountain Dew 12PK",
						"price": "6.49"
					},{
						"shortDescription": "Emils Cheese Pizza",
						"price": "abcd"
					},{
						"shortDescription": "Knorr Creamy Chicken",
						"price": "1.26"
					},{
						"shortDescription": "Doritos Nacho Cheese",
						"price": "3.35"
					},{
						"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
						"price": "12.00"
					}
				],
				"total": "35.35"
			}
			`,
		},
		{
			name:           "POST /receipts/process returns 400 with invalid purchaseDate",
			route:          "/receipts/process",
			method:         http.MethodPost,
			expectedStatus: http.StatusBadRequest,
			requestBody: `
			{
				"retailer": "Target",
				"purchaseDate": "abc",
				"purchaseTime": "13:01",
				"items": [
					{
						"shortDescription": "Mountain Dew 12PK",
						"price": "6.49"
					},{
						"shortDescription": "Emils Cheese Pizza",
						"price": "6.49"
					},{
						"shortDescription": "Knorr Creamy Chicken",
						"price": "1.26"
					},{
						"shortDescription": "Doritos Nacho Cheese",
						"price": "3.35"
					},{
						"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
						"price": "12.00"
					}
				],
				"total": "35.35"
			}
			`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new store before each test for test isolation
			store.CreateReceiptStore()

			request := httptest.NewRequest(tc.method, tc.route, bytes.NewBuffer([]byte(tc.requestBody)))
			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			status := recorder.Code
			if status != tc.expectedStatus {
				t.Errorf("'%s %s' expected status: '%d', actual: '%d'", tc.method, tc.route, tc.expectedStatus, status)
			}
		})
	}
}
