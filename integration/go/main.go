package main

import (
	"github.com/apioo/sdkgen-go"
	"log"
	"strconv"
)

func main() {
	credentials := sdkgen.HttpBearer{
		Token: "foo",
	}

	client, err := NewClient("http://localhost:1080", credentials)
	if err != nil {
		log.Fatal(err)
	}

	assertGetHello(client)
	assertGetEntries(client)
	assertInsert(client)
	assertThrowException(client)
}

func assertGetHello(client *Client) {
	message, err := client.Test().GetHello()
	if err != nil {
		log.Fatal(err)
	}

	if message.Message != "Hello World!" {
		log.Fatal("Test assertGetHello failed: Message, got: " + message.Message)
	}
}

func assertGetEntries(client *Client) {
	todos, err := client.Test().GetEntries(0, 16)
	if err != nil {
		log.Fatal(err)
	}

	if todos.TotalResults != 4 {
		log.Fatal("Test assertGetEntries failed: TotalResults, got: " + strconv.FormatInt(int64(todos.TotalResults), 10))
	}

	if todos.StartIndex != 0 {
		log.Fatal("Test assertGetEntries failed: StartIndex, got: " + strconv.FormatInt(int64(todos.StartIndex), 10))
	}

	if todos.ItemsPerPage != 16 {
		log.Fatal("Test assertGetEntries failed: ItemsPerPage, got: " + strconv.FormatInt(int64(todos.ItemsPerPage), 10))
	}

	if len(todos.Entry) != 2 {
		log.Fatal("Test assertGetEntries failed: Entry, got: " + strconv.FormatInt(int64(len(todos.Entry)), 10))
	}

	if todos.Entry[0].Title != "foo" {
		log.Fatal("Test assertGetEntries failed: Entry.0.Title, got: " + todos.Entry[0].Title)
	}

	if todos.Entry[1].Title != "bar" {
		log.Fatal("Test assertGetEntries failed: Entry.1.Title, got: " + todos.Entry[1].Title)
	}
}

func assertInsert(client *Client) {
	var payload = Todo{
		Title: "baz",
	}

	message, err := client.Test().Insert(payload)
	if err != nil {
		log.Fatal(err)
	}

	if message.Success != true {
		log.Fatal("Test assertInsert failed: Success, got: " + strconv.FormatBool(message.Success))
	}

	if message.Message != "Successful" {
		log.Fatal("Test assertInsert failed: Message, got: " + message.Message)
	}
}

func assertThrowException(client *Client) {
	_, err := client.Test().ThrowException()
	if err == nil {
		log.Fatal("Test assertThrowException failed: Expected an error")
	}

	if err.Error() != "" {
		log.Fatal("Test assertThrowException failed: Error message does not match, got: " + err.Error())
	}
}
