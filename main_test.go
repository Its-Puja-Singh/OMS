package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
	_ "github.com/go-sql-driver/mysql"
    "testing"
    "log"
    "net/http/httptest"
    "bytes"
    "os"
)

var testDb *sql.DB

func TestMain(m *testing.M) {
    // Create a test database and initialize the tables
    var err error

    testDbConStr := ""
    testDbHost := os.Getenv("DB_HOST")
    testDbPort := os.Getenv("DB_PORT")
    testDbUser := os.Getenv("DB_USER")
    testDbPassword := os.Getenv("DB_PASSWORD")
    testDbName := os.Getenv("TEST_DB_NAME")

    testDbConStr += testDbUser + ":" + testDbPassword + "@tcp(" + testDbHost + ":" +  testDbPort + ")/" + testDbName
    db, err = sql.Open("mysql", testDbConStr)
    if err != nil {
        log.Fatal(err)
    }

    // Run the tests
	code := m.Run()

	// Teardown
	os.Exit(code)

}

func TestCreateOrder(t *testing.T) {
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

    // Convert the order to JSON
    jsonOrder, _ := json.Marshal(testOrder)
    
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

    // Insert the test order into the database
    _, err := testDb.Exec("INSERT INTO orders (id, status, items, total, currency_unit) VALUES (?, ?, ?, ?, ?)",
        testOrder.ID, testOrder.Status, testOrder.Items, testOrder.Total, testOrder.CurrencyUnit)
    if err != nil {
        t.Fatal(err)
    }

    // Create a new request to get the test order
    req, err := http.NewRequest("GET", "/orders/test-id", nil)
    if err != nil {
        t.Fatal(err)
    }

    // Create a new response recorder to record the response
    rr := httptest.NewRecorder()

    // Call the getOrder function with the new request and response recorder
    handler := http.HandlerFunc(getOrders)
    handler.ServeHTTP(rr, req)

    // Check the status code of the response
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("getOrder returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Check the response body
    expectedBody := `{"id":"test-id","status":"PENDING_PAYMENT","items":[{"id":"test-item-id","description":"test item description","price":9.99,"quantity":2}],"total":19.98,"currencyUnit":"USD"}`
    if rr.Body.String() != expectedBody {
        t.Errorf("getOrder returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
    }
}

