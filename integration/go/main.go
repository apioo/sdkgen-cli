package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/apioo/sdkgen-go/v2"
	"log"
	"net/url"
	"strconv"
)

func main() {
	credentials := sdkgen.Anonymous{}

	client, err := NewClient("http://127.0.0.1:1080", credentials)
	if err != nil {
		log.Fatal(err)
	}

	assertGetHello(client)
	assertGetEntries(client)
	assertInsert(client)
	assertThrowException(client)
	assertBinary(client)
	assertForm(client)
	assertJson(client)
	assertText(client)
	assertXml(client)
}

func assertGetHello(client *Client) {
	message, err := client.Test().GetHello()
	if err != nil {
		log.Fatal(err)
	}

	if message.Message != "Hello World!" {
		log.Fatal("Test assertGetHello failed: Message, got: " + message.Message)
	}

	// Go has no option to set a default value at a struct
	/*
		if message.Category != "default" {
			log.Fatal("Test assertGetHello failed: Category, got: " + message.Category)
		}
	*/

	if message.Priority != 7 {
		log.Fatal("Test assertGetHello failed: Priority, got: " + strconv.FormatInt(int64(message.Priority), 10))
	}

	if message.Weight != 13.37 {
		log.Fatal("Test assertGetHello failed: Weight, got: " + strconv.FormatFloat(message.Weight, 'f', 10, 2))
	}

	if message.Done != true {
		log.Fatal("Test assertGetHello failed: Done, got: " + strconv.FormatBool(message.Done))
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

	var exception *ErrorException
	switch {
	case errors.As(err, &exception):
		if exception.Payload.Message != "Error" {
			log.Fatal("Test assertThrowException failed: Error message does not match, got: " + err.Error())
		}
	default:
		log.Fatal("Test assertThrowException failed: Error message does not match, got: " + err.Error())
	}
}

func assertBinary(client *Client) {
	var payload = []byte{0x66, 0x6F, 0x6F, 0x62, 0x61, 0x72}

	response, err := client.Test().Binary(payload)
	if err != nil {
		log.Fatal(err)
	}

	if !bytes.Equal(*response, payload) {
		log.Fatal("Test assertBinary failed")
	}
}

func assertForm(client *Client) {
	var payload = url.Values{}
	payload.Set("foo", "bar")

	response, err := client.Test().Form(payload)
	if err != nil {
		log.Fatal(err)
	}

	if response.Get("foo") != "bar" {
		log.Fatal("Test assertForm failed")
	}
}

func assertJson(client *Client) {
	var payload = make(map[string]string)
	payload["foo"] = "bar"

	response, err := client.Test().Json(payload)
	if err != nil {
		log.Fatal(err)
	}

	if response == nil {
		log.Fatal("Test assertJson failed, no response")
	}

	left, _ := json.Marshal(payload)
	right, _ := json.Marshal(&response)

	if !bytes.Equal(left, right) {
		log.Fatal("Test assertJson failed")
	}
}

func assertText(client *Client) {
	var payload = "foobar"

	response, err := client.Test().Text(payload)
	if err != nil {
		log.Fatal(err)
	}

	if payload != response {
		log.Fatal("Test assertText failed")
	}
}

func assertXml(client *Client) {
	var payload = "<foo>bar</foo>"

	response, err := client.Test().Xml(payload)
	if err != nil {
		log.Fatal(err)
	}

	if payload != response {
		log.Fatal("Test assertXml failed")
	}
}
