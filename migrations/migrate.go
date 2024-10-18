package migrations

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Migration struct {
	Name string
	Up   func(db *gorm.DB)
	Down func(db *gorm.DB)
}

var migrations = []Migration{
	{
		Name: "001_create_users_table",
		Up:   CreateUsersTableUp,
		Down: CreateUsersTableDown,
	},

	// Add other migrations here
}

// CreateMigrationsTable ensures that the migrations table exists before applying migrations.
func CreateMigrationsTable(db *gorm.DB) {
	err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) UNIQUE NOT NULL,
            applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
        );
    `).Error

	if err != nil {
		log.Fatalf("Error creating migrations table: %v", err)
	} else {
		log.Println("Migrations table is ready.")
	}
}

func RunMigrations(db *gorm.DB) {
	for _, migration := range migrations {
		var count int64
		db.Table("migrations").Where("name = ?", migration.Name).Count(&count)
		if count == 0 {
			log.Printf("Applying migration: %s", migration.Name)
			migration.Up(db)

			// Record the migration in the database
			err := db.Exec("INSERT INTO migrations (name, applied_at) VALUES (?, ?)", migration.Name, time.Now()).Error
			if err != nil {
				log.Fatalf("Failed to record migration: %v", err)
			} else {
				log.Printf("Migration applied: %s", migration.Name)
			}
		} else {
			log.Printf("Migration already applied: %s", migration.Name)
		}
	}
}
func RollbackMigrations(db *gorm.DB, steps int) {
	if steps > len(migrations) {
		steps = len(migrations)
	}

	for i := len(migrations) - 1; i >= 0 && steps > 0; i-- {
		var count int64
		db.Table("migrations").Where("name = ?", migrations[i].Name).Count(&count)
		if count > 0 {
			log.Printf("Rolling back migration: %s", migrations[i].Name)
			migrations[i].Down(db)
			db.Exec("DELETE FROM migrations WHERE name = ?", migrations[i].Name)
			steps--
		}
	}
}
