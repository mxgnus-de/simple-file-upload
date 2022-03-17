package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	envKey := os.Getenv("UPLOAD_KEY")
	envHost := os.Getenv("HOST_NAME")
	envScheme := os.Getenv("HOST_SCHEME")
	envPort := os.Getenv("PORT")

	if envKey == "" {
		panic("UPLOAD_KEY is not set. Please set it in .env file")
	}

	if envHost == "" {
		panic("HOST_NAME is not set. Please set it in .env file")
	}

	if envPort == "" {
		panic("PORT is not set. Please set it in .env file")
	}

	if envScheme == "" {
		panic("HOST_SCHEME is not set. Please set it in .env file")
	}

	_, err := os.Stat("uploads")

	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir("uploads", 0777)
		}
	}

	app := fiber.New()

	app.Use(logger.New())
	api := app.Group("/api")

	app.Static("/", "react-app/build")
	app.Static("/", "uploads")

	api.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		key := c.FormValue("key")

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "key is required",
			})
		}

		if key != envKey {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "key is invalid",
			})
		}

		err = c.SaveFile(file, "./uploads/"+file.Filename)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "File uploaded successfully",
			"path":    envScheme + "://" + envHost + "/" + file.Filename,
		})
	})

	app.Listen(":" + envPort)
}
