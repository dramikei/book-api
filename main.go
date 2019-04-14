package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dramikei/book-api/book"

	"github.com/labstack/echo"

	_ "github.com/go-sql-driver/mysql"
)

type Env struct {
	db *sql.DB
}

func (this *Env) setupDB() {
	db, err := sql.Open("mysql", "root:raghav@tcp(127.0.0.1:3306)/Library")

	this.db = db

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connected to Database")
	}
	err = this.db.Ping()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (this *Env) getBook(c echo.Context) (err error) {
	idInt, err := strconv.Atoi(c.Param("id"))
	id := uint32(idInt)
	if err != nil {
		return handleError(c, err)
	}
	var name string
	var author string
	var qty uint32

	get := "SELECT id, name, author, qty FROM BOOKS WHERE id = ?"

	err = this.db.QueryRow(get, id).Scan(&id, &name, &author, &qty)
	if err != nil {
		return handleError(c, err)
	}

	response := book.Book{ID: id, Name: name, Author: author, Qty: qty}
	return c.JSON(http.StatusOK, response)

}

func (this *Env) addBook(c echo.Context) (err error) {
	book := new(book.Book)
	if err := c.Bind(book); err != nil {
		return handleError(c, err)
	}
	sql := "INSERT INTO BOOKS(name, author, qty) VALUES(?, ?, ?)"
	stmt, err := this.db.Prepare(sql)

	if err != nil {
		return handleError(c, err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(book.Name, book.Author, book.Qty)
	if err != nil {
		return handleError(c, err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return handleError(c, err)
	}
	book.ID = uint32(id)
	return c.JSON(http.StatusOK, book)
}

func (this *Env) editBook(c echo.Context) (err error) {
	book := new(book.Book)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, err)
	}
	err1 := c.Bind(book)
	if err1 != nil {
		return handleError(c, err1)
	}
	book.ID = uint32(id)

	update := "UPDATE BOOKS SET name=?, author=?, qty=? WHERE id=?"

	stmt, err := this.db.Prepare(update)

	if err != nil {
		return handleError(c, err)
	}

	defer stmt.Close()
	result, err := stmt.Exec(book.Name, book.Author, book.Qty, id)
	if err != nil {
		return handleError(c, err)
	}
	fmt.Println(result.RowsAffected())
	return c.JSON(http.StatusOK, book)
}

func (this *Env) deleteBook(c echo.Context) (err error) {
	id := c.Param("id")

	delete := "DELETE from BOOKS where id=?"

	stmt, err := this.db.Prepare(delete)
	if err != nil {
		return handleError(c, err)
	}
	result, err := stmt.Exec(id)
	if err != nil {
		return handleError(c, err)
	}
	fmt.Println(result.RowsAffected())
	stmt.Close()
	this.resetAutoIncrement(c)
	return c.String(http.StatusOK, "Deleted.")
}

func (this *Env) resetAutoIncrement(c echo.Context) (err error) {
	maxID := "SELECT MAX(`id`) FROM `Books`"
	var number int
	err = this.db.QueryRow(maxID).Scan(&number)
	if err != nil {
		return handleError(c, err)
	}
	num := number + 1
	alterID := fmt.Sprintf("ALTER TABLE Books AUTO_INCREMENT= %d", num)
	fmt.Println(alterID, number)
	stmt, err := this.db.Prepare(alterID)
	if err != nil {
		return handleError(c, err)
	}
	defer stmt.Close()
	fmt.Println(err)
	result, err := stmt.Exec()
	fmt.Println(result.RowsAffected())
	if err != nil {
		return handleError(c, err)
	}
	return nil
}

func handleError(c echo.Context, e error) error {
	fmt.Println(e)
	return c.String(http.StatusInternalServerError, e.Error())
}

func main() {

	env := new(Env)
	env.setupDB()

	defer env.db.Close()
	e := echo.New()

	e.GET("/books/:id", env.getBook)
	e.POST("/books/", env.addBook)
	e.PUT("/books/:id", env.editBook)
	e.DELETE("/books/:id", env.deleteBook)

	e.Logger.Fatal(e.Start(":1323"))
}
