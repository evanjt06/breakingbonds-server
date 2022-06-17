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
							return
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

						claims := server.ExtractJwtClaims(c)
						uid := int64((claims["uid"]).(float64))
						log.Println(uid)

						quizNumber := c.Param("number")
						quizNumberInt,err := strconv.Atoi(quizNumber)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						r1 := c.PostForm("Response1")
						r2 := c.PostForm("Response2")
						r3 := c.PostForm("Response3")
						k1 := c.PostForm("Key1")
						k2 := c.PostForm("Key2")
						k3 := c.PostForm("Key3")
						elapsedTime := c.PostForm("ElapsedTime")

						if r1 == "" || r2 == "" || r3 == "" || k1 == "" || k2 == "" || k3 == "" || elapsedTime == "" {
							c.JSON(500, "invalid input (missing)")
							return
						}

						et, err := time.Parse("2006-01-02 15:04", elapsedTime)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						qr := internal.QuizResponses{
							QuizID:      int64(quizNumberInt),
							UserID:      uid,
							Response1:   r1,
							Response2:   r2,
							Response3:   r3,
							Key1:        k1,
							Key2:        k2,
							Key3:        k3,
							ElapsedTime: et,
						}
						qr.UseDBWriterPreferred()
						err = qr.Set()
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						c.JSON(200, "")
					},
				},

				{
					RelativePath: "/quiz",
					Method: ginhttpmethod.POST,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						claims := server.ExtractJwtClaims(c)
						uid := int64((claims["uid"]).(float64))
						log.Println(uid)

						pn := c.PostForm("PacketNumber")
						unit := c.PostForm("Unit")
						d := c.PostForm("Difficulty") // 1easy, 2medium, 3hard
						pdflink := c.PostForm("PDFLink")
						timer := c.PostForm("Timer")

						pn_int, err := strconv.Atoi(pn)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}
						unit_int, err := strconv.Atoi(unit)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}
						d_int, err := strconv.Atoi(d)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						if pn == "" || unit == "" || d == "" || pdflink == "" || timer == "" {
							c.JSON(500, "invalid input (missing)")
							return
						}

						t, err := time.Parse("2006-01-02 15:04", timer)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						q := internal.Quiz{
							PacketNumber: pn_int,
							UnitNumber:   unit_int,
							Difficulty:   d_int,
							PDFLink:      pdflink,
							Timer:        t,
							AdminID:      1,
						}
						q.UseDBWriterPreferred()
						err = q.Set()
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						c.JSON(200, "")
					},
				},

				{
					RelativePath: "/quiz/:number",
					Method:       ginhttpmethod.GET,
					Handler: func(c *gin.Context, bindingInputPtr interface{}) {

						claims := server.ExtractJwtClaims(c)
						uid := int64((claims["uid"]).(float64))
						log.Println(uid)

						quizNumber := c.Param("number")
						quizNumberInt, err := strconv.Atoi(quizNumber)
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						q := internal.Quiz{}
						q.UseDBWriterPreferred()
						notFound, err := q.GetByPacketNumber(quizNumberInt)
						if notFound {
							c.JSON(500, "not found")
							return
						}
						if err != nil {
							c.JSON(500, err.Error())
							return
						}

						c.JSON(200, gin.H{
							"quiz": q,
						})

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
