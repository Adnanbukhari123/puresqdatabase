// migrations/001_create_users_table.go
package migrations

import (
	"gorm.io/gorm"
	"log"
)

func CreateUsersTableUp(db *gorm.DB) {
	err := db.Exec(`
        CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            first_name VARCHAR(50) NOT NULL,
            last_name VARCHAR(50) NOT NULL,
            email VARCHAR(50) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP               
        );
    `).Error

	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	} else {
		log.Println("Successfully created users table.")
	}
}

func CreateUsersTableDown(db *gorm.DB) {
	err := db.Exec(`DROP TABLE IF EXISTS users;`).Error

	if err != nil {
		log.Fatalf("Error dropping users table: %v", err)
	} else {
		log.Println("Successfully dropped users table.")
	}
}
