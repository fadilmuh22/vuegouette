package product

import "restskuy/cmd/db"

func GetAllProduct() ([]Product, error) {
	var products []Product
	c := db.Connect()

	result, err := c.Query("SELECT * FROM product")
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var product Product
		err := result.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetProduct(id string) (Product, error) {
	var product Product
	c := db.Connect()

	result, err := c.Query("SELECT * FROM product WHERE id = ?", id)
	if err != nil {
		return product, err
	}

	defer result.Close()

	errorScan := result.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if errorScan != nil {
		return product, errorScan
	}

	return product, nil
}

func CreateProduct(product Product) (Product, error) {
	c := db.Connect()

	_, err := c.Exec("INSERT INTO product (name, price, stock) VALUES (?, ?, ?)", product.Name, product.Price, product.Stock)

	if err != nil {
		return product, err
	}

	return product, nil
}

func UpdateProduct(id string, product Product) (Product, error) {
	c := db.Connect()

	_, err := c.Exec("UPDATE product SET name = ?, price = ?, stock = ? WHERE id = ?", product.Name, product.Price, product.Stock, id)

	if err != nil {
		return product, err
	}

	return product, nil
}

func DeleteProduct(id string) error {
	c := db.Connect()

	_, err := c.Exec("DELETE FROM product WHERE id = ?", id)

	if err != nil {
		return err
	}

	return nil
}
