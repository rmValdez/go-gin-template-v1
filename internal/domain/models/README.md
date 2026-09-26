# Database Models

This directory contains the data structures that map directly to our PostgreSQL database tables using GORM. 

Think of this folder as the Go equivalent of your `schema.prisma` file, but heavily modularized.

## 🏛️ Architecture Rules

1. **Group by Domain Module:** We keep models modular by grouping related tables into domain-specific files (Domain-Driven Design). For example, `auth_schema.go` contains `User`, `Account`, and `Session` because they all belong to the identity module. If you add e-commerce tables, put them in `commerce_schema.go`.
2. **GORM Tags are your Schema:** We use Go struct tags (e.g., ``gorm:"uniqueIndex;not null"``) to define database constraints, relationships, and primary keys.
3. **No Business Logic:** Models should **only** contain struct definitions. Do not write SQL queries or API validation logic here. 
   - Database queries belong in `internal/repositories/`
   - API validation belongs in `internal/handlers/` (DTOs)
4. **JSON Serialization:** Use `json` struct tags to define how the model looks if it is serialized. Use `json:"-"` to completely hide sensitive fields (like passwords) from ever accidentally leaking in an API response.

---

## 🛠️ How to Add or Modify a Model

### 1. Create or Edit the Struct
Create a new file (e.g., `product.go`) or edit an existing one. Define your fields and their GORM constraints:

```go
package models

import "time"

type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null;index" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

### 2. Register it for Auto-Migration
If you created a **new** model, you must tell GORM to build the table in PostgreSQL! 

Open `cmd/api/main.go` and add your new struct to the `AutoMigrate` function:

```go
database.DB.AutoMigrate(
    &models.User{},
    &models.Account{},
    &models.Session{},
    &models.Product{}, // <-- Add your new model here!
)
```

The next time you run `docker-compose up`, GORM will automatically create or update the table schema to match your struct exactly like `npx prisma db push`!
