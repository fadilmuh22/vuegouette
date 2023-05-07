package product

import (
	"github.com/fadilmuh22/restskuy/cmd/db"
	"github.com/fadilmuh22/restskuy/cmd/model"
)

func GetAllProduct() ([]model.Product, error) {
	var products []model.Product
	c := db.Connect()

	result, err := c.Query("SELECT id, name, price, description, stock FROM product")
	if err != nil {
		return nil, err
	}

	defer result.Close()

	for result.Next() {
		var product model.Product
		err := result.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.Stock)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

func GetProduct(id string) (model.Product, error) {
	var product model.Product
	c := db.Connect()

	result, err := c.Query("SELECT id, name, price, description, stock FROM product WHERE id = ?", id)
	if err != nil {
		return product, err
	}

	defer result.Close()

	for result.Next() {
		err := result.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return product, err
		}
	}

	return product, nil
}

// create product with db and return product from db
func CreateProduct(product model.Product) (model.Product, error) {
	c := db.Connect()

	// insert product to db using sql query store in result
	result, err := c.Exec("INSERT INTO product (name, price, description, stock) VALUES (?, ?, ?, ?)", product.Name, product.Price, product.Description, product.Stock)
	if err != nil {
		return product, err
	}

	// get last insert id from result
	id, err := result.LastInsertId()
	if err != nil {
		return product, err
	}

	product.ID = int(id)

	return product, nil
}

func UpdateProduct(id string, product model.Product) (model.Product, error) {
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
