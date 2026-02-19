package main

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pretty66/websocketproxy"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//	"github.com/labstack/echo/v5/middleware"
	//	"github.com/r3labs/sse/v2"
	"github.com/ziflex/lecho/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type (
	Room struct {
		ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
		Name      string    `form:"room_name" json:"room_name,omitempty" validate:"required,alphanumunicode,max=20"`
		VNCAddr   string    `form:"server_addr" json:"server_addr,omitempty" validate:"required,ip|url"`
		Port      uint16    `gorm:"type:uint;default:15901" form:"server_port" json:"server_port,omitempty" validate:"number,gte=0,lte=65535"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	RoomValidator struct {
		validator *validator.Validate
	}
)

var (
	WSProxies map[string]*websocketproxy.WebsocketProxy
)

func createRoom(c echo.Context, db *gorm.DB, ws *echo.Group) error { // change error messages
	r := new(Room)
	if err := c.Bind(r); err != nil {
		data := map[string]any{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	if err := c.Validate(r); err != nil {
		data := map[string]any{
			"message": err.Error(),
		}
		return c.JSON(http.StatusInternalServerError, data)
	}

	room := &Room{
		Name:    r.Name,
		VNCAddr: r.VNCAddr,
	}

	if err := db.Create(&room).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
	}

	wsproxy, err := websocketproxy.NewProxy("ws://"+room.VNCAddr+":"+strconv.FormatUint(uint64(room.Port), 10), func(r *http.Request) error {

		r.Header.Set("Origin", room.VNCAddr+strconv.FormatUint(uint64(room.Port), 10))
		return nil
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{
			"message": err.Error(),
		})
	}

	ws.GET("/roomproxy/"+room.ID.String(), func(c echo.Context) error {
		return echo.WrapHandler(wsproxy)(c)
	})

	return c.JSON(http.StatusOK, map[string]any{
		"room_id": room.ID,
	})

}

func getRoom(c echo.Context, db *gorm.DB) error {
	var r Room
	id := c.QueryParam("uuid")
	//fmt.Println(id)
	if result := db.First(&r, "id = ?", id); result.Error != nil {
		data := map[string]any{
			"message": result.Error,
		}
		return c.JSON(http.StatusNotFound, data)
	}
	return c.JSON(http.StatusOK, map[string]any{
		"server_addr": r.VNCAddr + ":" + strconv.FormatUint(uint64(r.Port), 10),
		"room_name":   r.Name,
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
	WSProxies = make(map[string]*websocketproxy.WebsocketProxy)
	e := echo.New()
	e.Logger = lecho.New(
		os.Stdout,
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("Web"),
	)

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
		Root:       "public",
		Index:      "index.html",
		HTML5:      true,
		Browse:     false,
		IgnoreBase: true,
	}))

	e.GET("/room", func(c echo.Context) error {
		return c.File("index.html")
	})

	api := e.Group("/api")

	e2 := echo.New()
	e2.Logger = lecho.New(
		os.Stdout,
		lecho.WithTimestamp(),
		lecho.WithCaller(),
		lecho.WithPrefix("Proxy"),
	)

	ws := e2.Group("/ws")

	api.POST("/create_room", func(c echo.Context) error {
		return createRoom(c, db, ws)
	})

	api.GET("/get_room", func(c echo.Context) error {
		return getRoom(c, db)
	})

	go e2.Start(":1454")
	e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	//e.Logger.Fatal(e.Start(":1323"))

}
