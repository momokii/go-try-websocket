package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/momokii/go-try-websocket/internal/ws"
)

func main() {

	// init websocket manager
	manager := ws.NewManager()

	engine := html.New("./web", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).Render("error", fiber.Map{
				"Code":    code,
				"Message": err.Error(),
			})
		},
	})

	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/web", "./web")

	app.Get("/", func(c *fiber.Ctx) error {

		return c.Render("index", fiber.Map{})
	})

	// websocket connection
	app.Get("/ws", adaptor.HTTPHandlerFunc(manager.ServeWS))

	// if using tls for https for accomodate wss in websocket can use below code to using it on local development, because browser will block wss connection if not using https
	// and when deploy on like gcp can just use app.Listen, because gcp will handle tls for us and the wss can still work
	// if err := app.ListenTLS(":3000", "/server.crt", "/server.key"); err != nil {
	// 	log.Println("Error running app: ", err)
	// }

	if err := app.Listen(":3000"); err != nil {
		log.Println("Error running app: ", err)
	}
}
