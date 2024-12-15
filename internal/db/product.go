package db
import (
    "fmt"
    "github.com/lib/pq" // Import pq for PostgreSQL-specific operations
)


type Product struct {
    ID                     int      `json:"id"`
    UserID                 int      `json:"user_id"`
    ProductName            string   `json:"product_name"`
    ProductDescription     string   `json:"product_description"`
    ProductImages          []string `json:"product_images"` // Maps to PostgreSQL ARRAY
    CompressedProductImages []string `json:"compressed_product_images"`
    ProductPrice           float64  `json:"product_price"`
    CreatedAt              string   `json:"created_at"` // Can use time.Time for stricter type
    UpdatedAt              string   `json:"updated_at"` // Same here
    ImageURL               string   `json:"image_url"`
}


func CreateProductInDB(product *Product) error {
    if DB == nil {
        return fmt.Errorf("database connection is not initialized")
    }

    query := `
        INSERT INTO products (user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at, updated_at, imageurl)
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW(), $7)
        RETURNING id
    `
    err := DB.QueryRow(query, product.UserID, product.ProductName, product.ProductDescription, 
        pq.Array(product.ProductImages), pq.Array(product.CompressedProductImages), product.ProductPrice,product.ImageURL).Scan(&product.ID)
    if err != nil {
        fmt.Println("Error executing query:", err) // Log the error for debugging
        return fmt.Errorf("failed to create product: %w", err)
    }
    return nil
}



func GetProductByIDFromDB(id int) (*Product, error) {
    product := &Product{}
    query := `
        SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at, updated_at
        FROM products
        WHERE id = $1
    `
    err := DB.QueryRow(query, id).Scan(&product.ID, &product.UserID, &product.ProductName, 
        &product.ProductDescription, pq.Array(&product.ProductImages), pq.Array(&product.CompressedProductImages),
        &product.ProductPrice, &product.CreatedAt, &product.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return product, nil
}

func GetProductsByUserFromDB(userID int) ([]Product, error) {
    query := `
        SELECT id, user_id, product_name, product_description, product_images, compressed_product_images, product_price, created_at, updated_at
        FROM products
        WHERE user_id = $1
    `
    rows, err := DB.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []Product
    for rows.Next() {
        var product Product
        if err := rows.Scan(&product.ID, &product.UserID, &product.ProductName, &product.ProductDescription, 
            pq.Array(&product.ProductImages), pq.Array(&product.CompressedProductImages), &product.ProductPrice, 
            &product.CreatedAt, &product.UpdatedAt); err != nil {
            return nil, err
        }
        products = append(products, product)
    }
    return products, nil
}

