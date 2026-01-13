package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	//	"github.com/gorilla/websocket"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// func getRoom(c echo.Context) error {
//     return
// }

type (
    Room struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name string `form:"room_name" json:"room_name,omitempty" validate:"required,alphanumunicode"`
	VNCAddr string `form:"server_addr" json:"server_addr,omitempty" validate:"required,ip|url"`
	Port uint16 `gorm:"type:uint;default:15901" form:"server_port" json:"server_port,omitempty" validate:"number,gte=0,lte=65535"`
	CreatedAt time.Time
	UpdatedAt time.Time
    }

    RoomValidator struct {
	validator *validator.Validate
    }
)

func createRoom (c echo.Context, db *gorm.DB) error { // change error messages
    r := new(Room)
    if err := c.Bind(r); err != nil {
	data := map[string]any {
	    "message": err.Error(),
	}
	return c.JSON(http.StatusInternalServerError, data)
    }

    if err := c.Validate(r); err != nil {
	data := map[string]any {
	    "message": err.Error(),
	}
	return c.JSON(http.StatusInternalServerError, data)
    }

    room := &Room{
	Name: r.Name,
	VNCAddr: r.VNCAddr,
    }

    if err := db.Create(&room).Error; err != nil {
	return c.JSON(http.StatusInternalServerError, map[string]any {
	    "message": err.Error(),
	})
    }

    return c.JSON(http.StatusOK, map[string]any {
	"room_id": room.ID,
    })
    
}

func getRoom(c echo.Context, db *gorm.DB) error {
    var r Room
    id := c.QueryParam("uuid")
    fmt.Println(id)
    if result := db.First(&r, "id = ?", id); result.Error != nil {
	data := map[string]any {
	    "message": result.Error,
	}
	return c.JSON(http.StatusNotFound, data)
    }
    return c.JSON(http.StatusOK, map[string]any {
	"server_addr": r.VNCAddr + ":" + strconv.FormatUint(uint64(r.Port), 10),
    })
}

func (rv *RoomValidator) Validate(i any) error {
    if err := rv.validator.Struct(i); err != nil {
	return err
    }
    return nil
}

func main() {
    if err := godotenv.Load(); err != nil {
	log.Fatal("Error loading .env")
    }

    e := echo.New()
    e.Validator = &RoomValidator{validator: validator.New()}

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/New_York",
	os.Getenv("POSTGRES_HOST"),
	os.Getenv("POSTGRES_USER"),
	os.Getenv("POSTGRES_PASSWORD"),
	os.Getenv("POSTGRES_DBNAME"),
	os.Getenv("POSTGRES_PORT"),
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
	e.Logger.Fatal(err)
    }
    err = db.AutoMigrate((&Room{}))
    if err != nil {
	e.Logger.Fatal(err)
    }

    e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	Skipper: nil,
	// Root directory from where the static content is served.
	Root: "static",
	// Index file for serving a directory.
	// Optional. Default value "index.html".
	Index: "index.html",
	// Enable HTML5 mode by forwarding all not-found requests to root so that
	// SPA (single-page application) can handle the routing.
	HTML5:      true,
	Browse:     false,
	IgnoreBase: false,
	Filesystem: nil,
    }))

    api := e.Group("/api")

    api.POST("/create_room", func(c echo.Context) error {
	return createRoom(c, db)
    })

    api.GET("/get_room", func(c echo.Context) error {
	return getRoom(c, db)
    })
   
    e.Use(middleware.RequestLogger())
    e.Use(middleware.Recover())

    // e.GET("/room/:id", getRoom)
    // e.GET("/chat", handleChatServer)
    e.Logger.Fatal(e.Start(":1323"))
}
