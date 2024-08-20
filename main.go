package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
    "time"

	"github.com/rvk/rvkBooking/models"
	"github.com/rvk/rvkBooking/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

    "github.com/gofiber/template/html/v2"

)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Qrdetails struct {
	Qrcode    	string    `json:"qrcode"`
	Count     	int    		`json:"count"`
	First_scanned_at time.Time `json:"first_scanned_at"`
}

type Repository struct {
	DB *gorm.DB
}


func (r *Repository) FetchDetails(context *fiber.Ctx) error {

    id := context.Params("id")
	bookingModel := &models.Booking{}
    qrCode := &models.Qrdetails{};
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("qrcode = ?", id).First(bookingModel).Error;
   // result := r.DB.Where("qrcode = ?", id)
    result :=r.DB.First(&qrCode, "qrcode = ?", id);
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get the details"})
		return err
	} else {
        fmt.Println("the row is", result.RowsAffected);
        if result.RowsAffected == 0{
            user := Qrdetails{Qrcode: id, Count: 1, First_scanned_at: time.Now()}
            result := r.DB.Create(&user) // pass pointer of data to Create
            fmt.Println("the rocountw is", result.RowsAffected)
            fmt.Println("the erorr is", result.RowsAffected)
            fmt.Println("the roddddddw is", id)
            context.Status(http.StatusOK).Render("index", fiber.Map{
                "AttendeeName": bookingModel.Name,
                "SmartCardNo":bookingModel.Smartcard,
                "Category":bookingModel.Category,
                "Gate":bookingModel.Gateno,
                })
        } else{
            context.Status(http.StatusBadRequest).JSON(
                &fiber.Map{"message": "Already Scanned"})
            return err
        }
	return nil

    }




}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

    api.Get("qr_code/:id", r.FetchDetails);
    http.Handle("/", http.FileServer(http.Dir("./static")))
}



func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}
	//err = models.MigrateBooking(db)
	// if err != nil {
	// 	log.Fatal("could not migrate db")
	// }

	r := Repository{
		DB: db,
	}

    engine := html.New("./views", ".html")


	app := fiber.New(fiber.Config{
		Views: engine,
	})
    app.Static("/", "/views")
	//app := fiber.New()

	r.SetupRoutes(app)

	app.Listen(":8080")
}
