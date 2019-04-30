package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Book struct {
	Name   string
	Author string
	Qty    uint32
}

func main() {
	getBooks()
	addBook()
	getSecondBook()
	deleteBook()
	editBook()
	queryBook()
	log.Println("All tests passing.. will try testing error handling.")
	getWrongBook()
	wrongQuery()
	editAtWrongIndex()
	wrongDelete()
}

func getBooks() {
	log.Println("Testing GetBooks.. Expecting an array of Books")
	res, err := http.Get("http://localhost:1323/books/")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("Printing Books..")
	result := string(body)
	log.Println("Result: ", result)
	log.Println("Test Passed.")
}

func addBook() {
	log.Println("Testing POST Request: addBook.. Expecting a response 201 (Success)")
	reqBody, err := json.Marshal(Book{"Book1", "Author1", 4})
	if err != nil {
		log.Fatalln(err)
	}
	res, err := http.Post("http://localhost:1323/books/", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	log.Println("Result: ", result)
	log.Println("Test Passed....")
	log.Println("Testing GetBooks again with the new Record added to the array...")
	getBooks()

}

func getSecondBook() {
	log.Println("Testing GetBooks.. Expecting an object of Book")
	res, err := http.Get("http://localhost:1323/books/2")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("Printing Books..")
	result := string(body)
	log.Println("Result: ", result)
	log.Println("Test Passed.")
}

func deleteBook() {
	client := &http.Client{}
	log.Println("Testing DeleteBook.. Expecting a deletion Book at Index 1")
	req, err := http.NewRequest("DELETE", "http://localhost:1323/books/1", nil)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	log.Println("Result: ", result)
	if result == "deleted." {
		log.Println("Test Passed.")
	} else {
		log.Fatalln("Test failing...")
	}
}

func editBook() {
	client := &http.Client{}
	log.Println("Testing editing book...")
	getSecondBook()
	reqBody, err := json.Marshal(Book{"EditedBook", "EditedAuthor2", 45})
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("PUT", "http://localhost:1323/books/2", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	log.Println("Result: ", result)
	log.Println("Test Passed....")
	log.Println("Edited record is....")
	getSecondBook()
}

func addQueryBook() {
	log.Println("Testing POST Request: addBook.. Expecting a response 201 (Success)")
	reqBody, err := json.Marshal(Book{"QueryBook", "Author1", 43})
	if err != nil {
		log.Fatalln(err)
	}
	res, err := http.Post("http://localhost:1323/books/", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Added QueryBook...")
	result := string(body)
	log.Println("Result: ", result)
}

func queryBook() {
	addQueryBook()
	log.Println("Testing QueryBook.. Expecting an array of Books")
	res, err := http.Get("http://localhost:1323/books/q/?name=QueryBook")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("Printing Books..")
	result := string(body)
	log.Println("Result: ", result)
	log.Println("Test Passed.")
}

func getWrongBook() {
	log.Println("Testing book at invalid index.. Expecting an error of No record found")
	res, err := http.Get("http://localhost:1323/books/653784654378")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("Printing Book..")
	result := string(body)
	log.Println("Result: ", result)
	if result == "record not found" {
		log.Println("Test Passed.")
	} else {
		log.Fatalln("Test failing..")
	}
}

func wrongQuery() {
	log.Println("Testing QueryBook.. Expecting an array of Books")
	res, err := http.Get("http://localhost:1323/books/q/?name=ThisIsJustARidiculouslyLongNameForABook")
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	println("Printing Books..")
	result := string(body)
	log.Println("Result: ", result)
	if result == "[]\n" {
		log.Println("Test Passed.")
	} else {
		log.Fatalln("Test Failing...")
	}
}

func editAtWrongIndex() {
	client := &http.Client{}
	log.Println("Testing editing book...")
	reqBody, err := json.Marshal(Book{"EditedBook", "EditedAuthor2", 45})
	if err != nil {
		log.Fatalln(err)
	}
	req, err := http.NewRequest("PUT", "http://localhost:1323/books/265743865784", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	log.Println("Result: ", result)
	if result == "record not found" {
		log.Println("Test Passed....")
	} else {
		log.Fatalln("Test Failing...")
	}

}

func wrongDelete() {
	client := &http.Client{}
	log.Println("Testing DeleteBook at invalid index.. Expecting an error")
	req, err := http.NewRequest("DELETE", "http://localhost:1323/books/1574957349574835483", nil)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	result := string(body)
	log.Println("Result: ", result)
	if result == "record not found" {
		log.Println("Test Passed.")
	} else {
		log.Fatalln("Test failing...")
	}
}
