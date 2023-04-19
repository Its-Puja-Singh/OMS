package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
	_ "github.com/go-sql-driver/mysql"
    "testing"
    "net/http/httptest"
    "bytes"
    "os"
    // "fmt"
    "github.com/google/go-cmp/cmp"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
    // Run the tests
	code := m.Run()

	// Teardown
	os.Exit(code)

}

func TestCreateOrder(t *testing.T) {
    // Create a test order
    jsonOrder := []byte(`{
        "id": "abcdef-123458",
        "status": "PENDING_INVOICE",
        "items": [{
            "id": "123457 ",
            "description": "a product description",
            "price": 12.40,
            "quantity": 1
            }],
        "total": 12.40,
        "currencyUnit": "USD"
    }`)
    
    // Create a new request with the test order as the body
    req, err := http.NewRequest("POST", "/createOrder", bytes.NewBuffer(jsonOrder))
    if err != nil {
        t.Fatal(err)
    }

    // Create a new response recorder to record the response
    rr := httptest.NewRecorder()

    // Call the createOrder function with the new request and response recorder
    handler := http.HandlerFunc(createOrder)
    handler.ServeHTTP(rr, req)

    // Check the status code of the response
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("createOrder returned wrong status code: got %v want %v", status, http.StatusCreated)
    }

    // Check the response body
    expectedBody := "Order created successfully"
    if rr.Body.String() != expectedBody {
        t.Errorf("createOrder returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
    }
}

func TestGetOrder(t *testing.T) {
    // Create a test order
    testOrder := Order{
        ID: "test-id",
        Status: "PENDING_PAYMENT",
        Items: []Item{
            {
                ID: "test-item-id",
                Description: "test item description",
                Price: 9.99,
                Quantity: 2,
            },
        },
        Total: 19.98,
        CurrencyUnit: "USD",
    }
    jsonItems, err := json.Marshal(testOrder.Items)
    if err != nil {
        t.Fatal(err)
    }
    // Insert the test order into the database
    _, err = db.Exec("INSERT INTO orders (id, status, items, total, currency_unit) VALUES (?, ?, ?, ?, ?)",
        testOrder.ID, testOrder.Status, jsonItems, testOrder.Total, testOrder.CurrencyUnit)
    if err != nil {
        t.Fatal(err)
    }
    // data := url.Values{}
    // data.Set("id", "test-id")

    // Create a new request to get the test order
    // req, err := http.NewRequest("GET", "/getOrdersById", strings.NewReader(data.Encode()))
    req, err := http.NewRequest("GET", "/getOrdersById?id=test-id", nil)

    if err != nil {
        t.Fatal(err)
    }

    // Create a new response recorder to record the response
    rr := httptest.NewRecorder()

    // Call the getOrder function with the new request and response recorder
    handler := http.HandlerFunc(getOrdersById)
    handler.ServeHTTP(rr, req)

    // Check the status code of the response
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("getOrder returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    expectedBodyJSON := `{"id":"test-id","status":"PENDING_PAYMENT","items":[{"id":"test-item-id","description":"test item description","price":9.99,"quantity":2}],"total":19.98,"currencyUnit":"USD"}`
    var expectedBodyMap, actualBodyMap map[string]interface{}
    err = json.Unmarshal([]byte(expectedBodyJSON), &expectedBodyMap)
    if err != nil {
        t.Fatal(err)
    }
    err = json.Unmarshal(rr.Body.Bytes(), &actualBodyMap)
    if err != nil {
        t.Fatal(err)
    }

    // Compare the expected and actual response bodies using go-cmp
    if diff := cmp.Diff(expectedBodyMap, actualBodyMap); diff != "" {
        t.Errorf("getOrder returned unexpected body (-want +got):\n%s", diff)
    }
}