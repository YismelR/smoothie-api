package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/session"
)
type User struct{
	FirstName string `json:"firstName" form:"firstName"`
	LastName string `json:"lastName" form:"lastName"`
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type SignIn struct{
	Email string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func main(){
	users := []User{}

	sessionStore := session.New(session.Config{
		CookieHTTPOnly: true,
		KeyLookup: "cookie:session_id",
	})
	
	app := fiber.New()
	

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowCredentials: true,
	}))

	app.Delete("/logout",func(c *fiber.Ctx) error {
		sess, err := sessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
		}

		err = sess.Destroy()

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Something went wrong")
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ok": true,
				"status": fiber.StatusOK,
		})

	})

	app.Get("/check-auth", func(c *fiber.Ctx) error {
		sess, err := sessionStore.Get(c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong")
		}

		userEmail := sess.Get("user_email")

		if userEmail != nil{
			// { ok: true, status: 200 }
			return c.JSON(fiber.Map{
				"ok": true,
				"status": fiber.StatusOK,
			})
		}

		// { ok: false, status: 401 }
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"ok": false,
				"status": fiber.StatusUnauthorized,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {

		sess, err := sessionStore.Get(c)

		if err != nil {
			return c.Status(500).SendString("Something went wrong")
		}

		logIn := new(SignIn)

		err = c.BodyParser(logIn)

		if err != nil {
			return err
		}

		
		var foundUser User
		
		for _, user := range users{

			if logIn.Email == user.Email{
				foundUser = user
			}
		}


		if foundUser.Password == logIn.Password{
			
			sess.Set("user_email", logIn.Email)

			err := sess.Save()

			if err != nil{
				return c.Status(fiber.StatusInternalServerError).SendString("something went wrong")
			}

			return c.SendString("Success")	
		}
		
		return c.Status(400).SendString("Invalid credentials")

	})

	app.Post("/register", func(c *fiber.Ctx) error {
		user := new(User)

		err := c.BodyParser(user)

		if err != nil{
			return err
		}

		users = append(users, *user)

		return c.SendString("user was created")

	})

	app.Listen(":3000")
}