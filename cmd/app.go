package main

import (
	"avchem-server/internal"
	"database/sql"
	ginw "github.com/aldelo/common/wrapper/gin"
	"github.com/aldelo/common/wrapper/gin/ginhttpmethod"
	"github.com/aldelo/connector/webserver"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var server *webserver.WebServer

func main() {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// db init
	err = internal.SetWriterDBInfo(os.Getenv("HOST"), port, os.Getenv("DBNAME"), os.Getenv("USERNAME"), os.Getenv("PASSWORD"))
	if err != nil {
		log.Println(err)
	}

	err = internal.ConnectToWriterDB()

	if err != nil {
		log.Println(err)
	}

	// ping test
	isReady := internal.IsWriterDBReady()

	// ping 200
	if !isReady {
		log.Println("db not ready")
	}

	// last func
	defer internal.DisconnectFromWriterDB()

	server = webserver.NewWebServer("BreakingBonds", "config", "")

	server.LoginRequestDataPtr = &internal.Credentials{}

	server.LoginResponseHandler = func(c *gin.Context, statusCode int, token string, expires time.Time) {
		c.JSON(statusCode, gin.H{
			"token": token,
			"exp":   expires,
		})
	}
	server.AuthenticateHandler = func(loginRequestDataPtr interface{}) (loggedInCredentialPtr interface{}) {
		if lg, ok := loginRequestDataPtr.(*internal.Credentials); !ok {
			return nil
		} else {
			defer func() {
				lg.Password = ""
				lg.Email = ""
				lg.AdminID = 0
				lg.UserID = 0
				lg = nil
			}()

			// authenticate user
			uid, err := internal.ValidateCredentials(*lg)
			if err != nil {
				log.Println(err.Error())
				return nil
			}

			return &internal.Credentials{
				Email:    lg.Email,
				Password: lg.Password,
				UserID:   uid,
				AdminID:  0,
			}
		}
	}

	server.AddClaimsHandler = func(loggedInCredentialPtr interface{}) (identityKeyValue string, claims map[string]interface{}) {

		ptr, ok := loggedInCredentialPtr.(*internal.Credentials)

		if !ok {
			return "", nil
		}

		if loggedInCredentialPtr != nil {
			return "app", map[string]interface{}{
				"uid": ptr.UserID,
			}
		}

		return "", nil
	}

	server.AuthorizerHandler = func(loggedInCredentialPtr interface{}, c *gin.Context) bool {
		return true
	}

	server.Routes = map[string]*ginw.RouteDefinition{
		"base": {
			Routes: []*ginw.Route{
				{
					RelativePath: "/register",
					Method:       ginhttpmethod.POST,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						email := c.PostForm("email")
						if email == "" {
							c.JSON(500, "invalid input (EMAIL)")
							return
						}

						if !internal.IsEmailValid(email) {
							c.JSON(500, "invalid input (EMAIL)")
							return
						}

						password := c.PostForm("password")
						if password == "" || len(password) < 8 {
							c.JSON(500, "invalid input (PASSWORD)")
							return
						}

						user := internal.User{
							Email:    email,
							Password: password,
							Points: sql.NullInt32{
								Int32: 0,
								Valid: true,
							},
						}
						user.UseDBWriterPreferred()
						err = user.Set()
						if err != nil {
							c.JSON(500, err.Error())
						}

						c.JSON(200, "")
					},
				},
			},
		},
		"auth": {
			Routes: []*ginw.Route{

				{
					RelativePath: "/submitQuiz/:number",
					Method:       ginhttpmethod.POST,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						// todo - submit quiz

					},
				},

				{
					RelativePath: "/packet/:number",
					Method:       ginhttpmethod.GET,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						// todo - get specific packet pdf for display

					},
				},

				{
					RelativePath: "/history/:uid",
					Method:       ginhttpmethod.GET,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						// todo - get specific packet pdf for display

					},
				},
			},
		},
	}

	err = server.Serve()

	if err != nil {
		log.Println(err)
	}

}
