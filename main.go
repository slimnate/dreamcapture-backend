package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slimnate/dreamcapture-backend/data/booking"
)

func InitDB() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s", host, port, user, pass, dbName)

	// open db
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Cannot connect to database server: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Cannot communicate with database server: %s", err.Error())
	}

	appEnv := os.Getenv("APP_ENV")
	log.Printf("Using APP_ENV: %s", appEnv)
	if appEnv == "dev" {
		// dev environment, clear database
		_, err = db.Exec("DROP TABLE IF EXISTS users, organizations, sessions, events")
		if err != nil {
			log.Fatalf("Error dropping tables: %s", err.Error())
		}
	} else if appEnv == "prod" {
		// prod environment - do nothing here currently, tables will be created by migration functions for each repo if needed
	} else {
		// invalid environment
		log.Fatalf("Invalid value supplied for APP_ENV - '%s' - Must be either 'dev' or 'prod'", appEnv)
	}

	return db
}

func InitBooking(db *sql.DB) (*booking.BookingController, *booking.Repository) {
	repo := booking.NewBookingRepository(db)
	controller := booking.NewBookingController(repo)

	if err := repo.Migrate(); err != nil {
		log.Fatal("[bookings] Migration error: ", err)
	}

	return controller, repo
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("can't load env vars")
	}

	// init database and controllers
	db := InitDB()
	bookingController, _ := InitBooking(db)

	// init router
	router := gin.Default()

	// Add routes
	router.GET("/api/bookings", bookingController.List)
	//router.POST("/api/bookings")

	//Start router
	router.Run(":8080")
}
