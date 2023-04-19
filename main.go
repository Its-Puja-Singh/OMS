package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
	_ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "os"
    "fmt"
)

type Order struct {
    ID           string    `json:"id"`
    Status       string    `json:"status"`
    Items       []Item    `json:"items"`
    Total        float32   `json:"total"`
    CurrencyUnit string    `json:"currencyUnit"`
}

type Item struct {
    ID          string  `json:"id"`
    Description string  `json:"description"`
    Price       float32 `json:"price"`
    Quantity    int     `json:"quantity"`
}

var db *sql.DB

var router *mux.Router

func main() {

    // Create the router
    router = mux.NewRouter()

    // Define the API endpoints
    router.HandleFunc("/createOrder", createOrder).Methods("POST")
    router.HandleFunc("/updateOrder", updateOrder).Methods("PUT")
    router.HandleFunc("/getOrders", getOrders).Methods("GET")
    router.HandleFunc("/getOrdersById", getOrdersById).Methods("GET")
    router.HandleFunc("/getOrdersByParams", getOrdersByStatusAndCurrency).Methods("GET")

    // Start the server
    http.ListenAndServe(":8000", router)

}

func createOrder(w http.ResponseWriter, r *http.Request) {
    // Parse the request body
    var order Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Encode the items field as a JSON string
    itemsJSON, err := json.Marshal(order.Items)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Insert the order into the database
    _, err = db.Exec("INSERT INTO orders (id, status, items, total, currency_unit) VALUES (?, ?, ?, ?, ?)",
        order.ID, order.Status, string(itemsJSON), order.Total, order.CurrencyUnit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Return a success message
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Order created successfully"))
}

func updateOrder(w http.ResponseWriter, r *http.Request) {

    // Parse the request body
    var order Order
    err := json.NewDecoder(r.Body).Decode(&order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
	    return
	}

	// Update the order in the database
	_, err = db.Exec("UPDATE orders SET status = ? WHERE id = ?", order.Status, order.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order updated successfully"))
}

func getOrders(w http.ResponseWriter, r *http.Request) {

    // Build the SQL query
    sql := "SELECT id, status, items, total, currency_unit FROM orders"

    // Execute the query
    rows, err := db.Query(sql)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Read the rows and convert them to orders
    var orders []Order
    for rows.Next() {
        var order Order
        var itemsJSON []byte
        err := rows.Scan(&order.ID, &order.Status, &itemsJSON, &order.Total, &order.CurrencyUnit)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Decode the items field from the database as a JSON string and unmarshal it into a slice of Item structs
        err = json.Unmarshal(itemsJSON, &order.Items)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        orders = append(orders, order)
    }

    // Write the orders as JSON to the response
    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(orders)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func getOrdersById(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
    // Build the SQL query
    sql := "SELECT id, status, items, total, currency_unit FROM orders where id = '" +  id + "'"

    // Execute the query
    rows, err := db.Query(sql)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Read the rows and convert them to orders
    var order Order
    for rows.Next() {
        var itemsJSON []byte
        err := rows.Scan(&order.ID, &order.Status, &itemsJSON, &order.Total, &order.CurrencyUnit)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Decode the items field from the database as a JSON string and unmarshal it into a slice of Item structs
        err = json.Unmarshal(itemsJSON, &order.Items)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Write the orders as JSON to the response
    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func getOrdersByStatusAndCurrency(w http.ResponseWriter, r *http.Request) {
    // Parse the query parameters
    status := r.URL.Query().Get("status")
    currencyUnit := r.URL.Query().Get("currency_unit")
    sortBy := r.URL.Query().Get("sort_by")

    // Build the SQL query
    query := "SELECT * FROM orders"
    if status != "" {
        query += " WHERE status = '" + status + "'"
    }
    if currencyUnit != "" {
        if status == "" {
            query += " WHERE"
        } else {
            query += " AND"
        }
        query += " currency_unit = '" + currencyUnit + "'"
    }
    if sortBy != "" {
        query += " ORDER BY " + sortBy
    }

    // Execute the SQL query
    rows, err := db.Query(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    // Parse the rows into a slice of orders
    orders := []Order{}
    for rows.Next() {
        var order Order
        var itemsStr string
        err := rows.Scan(&order.ID, &order.Status, &itemsStr, &order.Total, &order.CurrencyUnit)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = json.Unmarshal([]byte(itemsStr), &order.Items)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        orders = append(orders, order)
    }

    // Convert the orders to JSON and return them
    jsonOrders, err := json.Marshal(orders)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonOrders)
}

func init() {
    godotenv.Load()
    // Connect to the database
    var err error
    dbConStr := ""
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    dbConStr += dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" +  dbPort + ")/" + dbName
    db, err = sql.Open("mysql", dbConStr)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to Database")
}