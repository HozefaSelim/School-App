package Config

import (
	"SchoolProject/Models"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func Connect() {
	// DB connection
	// A Go (golang) port of the Ruby dotenv project (which loads env vars from a .env file).
	godotenv.Load()
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_NAME")

	connection := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbName)
	var db, err = gorm.Open(mysql.Open(connection), &gorm.Config{})

	if err != nil {
		panic("db connection failed")
	}
	DB = db
	fmt.Println("db connection successful")

	// we comment this function calls because we want to call only once
	//AutoMigration(db)

	// DB connection End
}
func AutoMigration(connection *gorm.DB) {
	connection.Debug().AutoMigrate(
		&Models.Class{},
		&Models.School{},
		&Models.Student{},
		&Models.Subject{},
		&Models.Teacher{},
	)
}
