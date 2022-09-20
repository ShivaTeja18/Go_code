package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "golang.org/x/crypto/openpgp/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var DB *gorm.DB
var r = gin.Default()

type EMP struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	City     string `json:"city"`
}

func connect() {
	const DNS = "host = 'localhost' port = '5432' dbname = 'API' user = 'postgres' password = 'Shiva@205101' sslmode = 'prefer'"
	Database, err := gorm.Open(postgres.Open(DNS), &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}
	if err == nil {
		log.Print("Connected")
	}
	_ = Database.AutoMigrate(&EMP{})

	DB = Database
}
func creating(c *gin.Context) {
	var empobj EMP

	if err := c.BindJSON(&empobj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
	}
	if empobj.Name == "" || empobj.Password == "" || empobj.City == "" {
		c.JSON(http.StatusNoContent, &empobj)
	} else {
		c.IndentedJSON(http.StatusOK, &empobj)
	}
	DB.Create(&empobj)
}

func fetching(a *gin.Context) {
	var empobj []EMP
	DB.Find(&empobj)
	a.IndentedJSON(http.StatusOK, &empobj)
}

func fbyid(a *gin.Context) {
	var empobj EMP

	id := a.Param("id")

	if id != id {
		a.JSON(http.StatusNotFound, gin.H{"message": "ID NOT FOUND"})
		return
	}
	if err := DB.Where("id = ?", id).Find(&empobj).Error; err != nil {
		a.JSON(http.StatusOK, &empobj)
		return
	}
}
func main() {
	connect()

	r.GET("/fetch", fetching)
	r.GET("/fetch/:id", fbyid)
	r.POST("/create", creating)
	//r.PUT("/change", updating)
	//r.DELETE("/delete", deleting)

	r.Run(":8080")
}
