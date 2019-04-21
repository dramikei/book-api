package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
)

type Book struct {
	gorm.Model
	Name   string `gorm:"type:varchar(255)"`
	Author string `gorm:"type:varchar(255)"`
	Qty    uint32 `gorm:"type:INT"`
}

type Env struct {
	db *gorm.DB
}

func (this *Env) setupDB() {
	db, err := gorm.Open("mysql", "root:raghav@tcp(127.0.0.1:3306)/Library?parseTime=True")
	this.db = db
	db.AutoMigrate(&Book{})

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connected to Database")
	}
}

func (this *Env) getBooks(c echo.Context) (err error) {
	var books []Book
	if err := this.db.Find(&books).Error; err != nil {
		return handle404Error(c, err)
	}
	return c.JSON(http.StatusOK, &books)
}

func (this *Env) getBook(c echo.Context) (err error) {
	idInt, err := strconv.Atoi(c.Param("id"))
	id := uint(idInt)
	if err != nil {
		return handleInternalError(c, err)
	}
	var response Book
	if err := this.db.First(&response, id).Error; err != nil {
		return handle404Error(c, err)
	}
	return c.JSON(http.StatusOK, response)

}

func (this *Env) addBook(c echo.Context) (err error) {
	book := new(Book)
	if err := c.Bind(book); err != nil {
		return handleInternalError(c, err)
	}
	this.db.Create(&book)
	return c.JSON(http.StatusOK, book)
}

func (this *Env) queryBook(c echo.Context) (err error) {
	var books []Book
	name := c.QueryParam("name")
	authorName := c.QueryParam("author")
	this.db.Where("name = ? OR author = ?", name, authorName).Find(&books)
	return c.JSON(http.StatusOK, &books)
}

func (this *Env) editBook(c echo.Context) (err error) {
	var book Book
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleInternalError(c, err)
	}

	if err := this.db.Where("id = ?", id).First(&book).Error; err != nil {
		return handle404Error(c, err)
	}

	if err := c.Bind(&book); err != nil {
		fmt.Println(book)
		return handleInternalError(c, err)
	}
	this.db.Save(&book)
	return c.JSON(http.StatusOK, book)
}

func (this *Env) deleteBook(c echo.Context) (err error) {
	idInt, err := strconv.Atoi(c.Param("id"))
	id := uint(idInt)
	if err != nil {
		return handleInternalError(c, err)
	}
	var book Book
	if err := this.db.Unscoped().Where("id=?", id).Delete(&book).Error; err != nil {
		return handle404Error(c, err)
	}
	return c.String(http.StatusOK, "Deleted.")
}

func handleInternalError(c echo.Context, e error) error {
	fmt.Println(e)
	return c.String(http.StatusInternalServerError, e.Error())
}

func handle404Error(c echo.Context, e error) error {
	fmt.Println(e)
	return c.String(http.StatusNotFound, e.Error())
}

func main() {

	env := new(Env)
	env.setupDB()

	defer env.db.Close()
	e := echo.New()

	e.GET("/books/", env.getBooks)
	e.GET("/books/:id", env.getBook)
	e.GET("/books/", env.queryBook)
	e.POST("/books/", env.addBook)
	e.PUT("/books/:id", env.editBook)
	e.DELETE("/books/:id", env.deleteBook)

	e.Logger.Fatal(e.Start(":1323"))
}
