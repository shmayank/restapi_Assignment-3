package main

// This is Simplest REST API Book CRUD example in Golang
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Product struct
type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Status      string `json:"status"`
	Image       string `json:"image"`
}

// db is an object of sql connection
var db *sql.DB
var err error

func main() {
	// DB Setup
	db, err = sql.Open("mysql", "root:@(127.0.0.1:3306)/shopping_cart")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected")
	defer db.Close()

	// Init Router, here we are using gorilla mux router
	router := mux.NewRouter().StrictSlash(true)

	// Route Handler which establish endpoints

	router.HandleFunc("/products", createProduct).Methods("POST")
	router.HandleFunc("/products", getallProducts).Methods("GET")
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products/{id}", updateProduct).Methods("PUT")
	router.HandleFunc("/products/{id}", deleteProduct).Methods("DELETE")

	// Run Server
	log.Fatal(http.ListenAndServe(":8000", router))
}

// createProduct create a new product
func createProduct(w http.ResponseWriter, r *http.Request) {

	statement, err := db.Prepare("INSERT INTO products(p_id, p_name, price, p_description, quantity, status, image)VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	dataMap := make(map[string]string)
	json.Unmarshal(body, &dataMap)

	id := dataMap["id"]
	ID, _ := strconv.Atoi(id)
	name := dataMap["name"]
	price, _ := strconv.Atoi(dataMap["price"])
	desc := dataMap["description"]
	quantity, _ := strconv.Atoi(dataMap["quantity"])
	status := dataMap["status"]
	image := dataMap["image"]

	fmt.Println(dataMap)
	_, err = statement.Exec(ID, name, price, desc, quantity, status, image)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New product added")
	_, err = statement.Exec(id, name, price, desc, quantity, status, image)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New product added")
}

// getProduct return a single product based on id which you pass.
func getProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	result, err := db.Query("SELECT * FROM products WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	// product is object of Book struct
	var product Product

	for result.Next() {
		err := result.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Quantity, &product.Status, &product.Image)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(product)
}

// getProducts returns all products.
func getallProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var products []Product

	result, err := db.Query("SELECT * FROM products")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Quantity, &product.Status, &product.Image)
		if err != nil {
			panic(err.Error())
		}
		products = append(products, product)

	}
	json.NewEncoder(w).Encode(products)
}

// updateProduct update the product name based on id passed.
func updateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	statement, err := db.Prepare("UPDATE products SET p_name = ? WHERE p_id = ?")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	dataMap := make(map[string]string)
	json.Unmarshal(body, &dataMap)
	newTitle := dataMap["p_name"]

	_, err = statement.Exec(newTitle, params["p_id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Product %s updated", params["p_id"])

}

// deleteProduct delete a product based on id
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	statement, err := db.Prepare("DELETE FROM products WHERE p_id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = statement.Exec(params["p_id"])
	if err != nil {
		panic(err.Error)
	}
	fmt.Fprintf(w, "Product %s deleted", params["p_id"])
}
