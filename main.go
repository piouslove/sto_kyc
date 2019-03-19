package main

import (
	"Sto_kyc/config"
	"Sto_kyc/controllers"
	"Sto_kyc/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	defer models.DB_mysql.Close()

	router := gin.Default()

	cf := cors.DefaultConfig()
	cf.AllowOrigins = []string{"*"}
	cf.AddAllowHeaders("Origin")
	cf.AllowCredentials = true
	// cf.AllowAllOrigins = true
	router.Use(cors.New(cf))

	router.GET("/KYCItems", controllers.GetKycItems)
	router.POST("/apply", controllers.Apply)
	router.POST("/query", controllers.Query)
	router.POST("/getDataToCheck", controllers.GetCheckData)
	router.POST("/certify", controllers.Certify)
	router.POST("/reject", controllers.Reject)

	/*
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "./front/index.html", nil)
		})

		router.Static("/index", "./front/index.html")
	*/
	router.Static("/passportImages/", config.V.ImagesDir)
	router.Static("/css/", "./front/css")
	router.Static("/js/", "./front/js")
	router.Static("/front/", "./front")

	router.StaticFile("/", "./front/index.html")

	port := ":" + config.V.Port
	router.Run(port)
}
