package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"  // or the Docker service name if running in another container
	port     = 5432         // default PostgreSQL port
	user     = "myuser"     // as defined in docker-compose.yml
	password = "mypassword" // as defined in docker-compose.yml
	dbname   = "mydatabase" // as defined in docker-compose.yml
)

var db *sql.DB

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Category string `json:"category"`
}

func main() {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	sdb, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	db = sdb
	defer db.Close()
	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	app := fiber.New()
	app.Get("/products", getproductsHandele)
	app.Get("/product/:id", getproductHandele)
	app.Post("/product/", productHandele)
	app.Put("/product/:id", updateHandele)
	app.Delete("/product/:id", deleteHandele)

	app.Listen(":8000")

	// product, err := getProduct(4)
	// fmt.Println("Get Success full", product)

	// err = createProduct(&Product{Name: "kk", Price: 30})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// product, err := updateProduct(1, &Product{Name: "UU", Price: 450})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Get Success full", product)

	// product, err := getProductsAndSuppliers()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("delete Success full", product)
}

func getproductHandele(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	product, err := getProduct(productId)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(product)
}

func productHandele(c *fiber.Ctx) error {
	p := new(Product)
	if err := c.BodyParser(p); err != nil {
		return err
	}

	err := createProduct(p)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(p)

}

func updateHandele(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	p := new(Product)

	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	product, err := updateProduct(productId, p)

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(product)

}

func deleteHandele(c *fiber.Ctx) error {
	productId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err = deleteProduct(productId)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.SendStatus(fiber.StatusNoContent)

}

func getproductsHandele(c *fiber.Ctx) error {
	product, err := getProducts()
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(product)
}
