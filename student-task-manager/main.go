package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
    _ "github.com/lib/pq"
)

func main() {
    // Database connection string
    dbURL := "postgresql://postgres:postgres@localhost:5432/student_task_manager?sslmode=disable"
    
    // Connect to PostgreSQL
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Create database if it doesn't exist (run this once)
    _, err = db.Exec("CREATE DATABASE student_task_manager")
    if err != nil {
        fmt.Println("Database may already exist:", err)
    }
    
    // Reconnect to the specific database
    db.Close()
    db, err = sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create migration driver
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize migrations
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "postgres", driver)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run migrations
    fmt.Println("Running migrations...")
    err = m.Up()
    if err != nil && err != migrate.ErrNoChange {
        log.Fatal("Migration failed:", err)
    }
    
    fmt.Println("✅ Migrations completed successfully!")
    
    // Verify tables
    verifyTables(db)
}

func verifyTables(db *sql.DB) {
    tables := []string{"users", "assignments", "tasks", "study_sessions", "reminders"}
    
    fmt.Println("\nVerifying tables:")
    for _, table := range tables {
        var exists bool
        err := db.QueryRow(`
            SELECT EXISTS (
                SELECT FROM information_schema.tables 
                WHERE table_name = $1
            )`, table).Scan(&exists)
        
        if err != nil {
            fmt.Printf("❌ Error checking %s: %v\n", table, err)
            continue
        }
        
        if exists {
            fmt.Printf("✅ Table '%s' created\n", table)
            
            // Count rows
            var count int
            db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
            fmt.Printf("   Contains %d rows\n", count)
        } else {
            fmt.Printf("❌ Table '%s' not found\n", table)
        }
    }
}