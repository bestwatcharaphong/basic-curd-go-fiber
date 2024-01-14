package main

func createProduct(product *Product) error {

	_, err := db.Exec(
		"INSERT INTO public.products(name, price) VALUES ($1, $2);",
		product.Name,
		product.Price,
	)

	return err
}

func getProduct(id int) (Product, error) {
	var p Product
	row := db.QueryRow(
		"SELECT id,name,price FROM products WHERE id=$1;",
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	if err != nil {
		return Product{}, err
	}

	return p, nil
}

func updateProduct(id int, product *Product) (Product, error) {
	var p Product

	row := db.QueryRow(
		"UPDATE products SET name=$1, price=$2 WHERE id=$3 RETURNING id,name,price;",
		product.Name,
		product.Price,
		id,
	)

	err := row.Scan(&p.ID, &p.Name, &p.Price)

	return p, err
}

func deleteProduct(id int) error {
	_, err := db.Exec(
		"DELETE FROM products WHERE id=$1;",
		id,
	)
	return err
}

func getProducts() ([]Product, error) {
	row, err := db.Query("SELECT id, name, price FROM products") //query data

	if err != nil { // check err form query
		return nil, err
	}

	var products []Product // declare products resive data is arrary

	for row.Next() { // loop all data
		var p Product
		err := row.Scan(&p.ID, &p.Name, &p.Price) // sacan data inseart memerry
		if err != nil {                           // check err from loop
			return nil, err
		}
		products = append(products, p)
	}

	if err = row.Err(); err != nil { //check error from loop
		return nil, err
	}

	return products, nil // return products
}

type ProductWithSupplier struct {
	ProductID    int
	ProductName  string
	Price        int
	SupplierName string
}

func getProductsAndSuppliers() ([]ProductWithSupplier, error) {
	// SQL JOIN query
	query := `
		SELECT
			p.id AS product_id,
			p.name AS product_name,
			p.price,
			s.name AS supplier_name
		FROM
			products p
		INNER JOIN suppliers s
			ON p.supplier_id = s.id;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []ProductWithSupplier
	for rows.Next() {
		var p ProductWithSupplier
		err := rows.Scan(&p.ProductID, &p.ProductName, &p.Price, &p.SupplierName)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
