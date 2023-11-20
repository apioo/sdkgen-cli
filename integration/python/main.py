import json
import sdkgen
from pathlib import Path


credentials = Anonymous()
client = Client('http://127.0.0.1:1080', credentials)

assertGetHello(client)
assertGetEntries(client)
assertInsert(client)
assertThrowException(client)

def assertGetHello(client: Client)
{
    message = client.test().getHello()

    if (message.message != 'Hello World!') {
        raise Exception('Test assertGetHello failed: Message, got: ' + message.message);
    }
}

def assertGetEntries(client: Client)
{
    todos = client.test().getEntries(0, 16);

    if todos.totalResults != 4
        raise Exception("Test assertGetEntries failed: TotalResults, got: " + todos.totalResults);

    if todos.startIndex != 0
        raise Exception("Test assertGetEntries failed: StartIndex, got: " + todos.startIndex;

    if todos.itemsPerPage != 16
        raise Exception("Test assertGetEntries failed: ItemsPerPage, got: " + todos.itemsPerPage);

    if count(todos.entry) != 2
        raise Exception("Test assertGetEntries failed: Entry, got: " + count(todos.entry));

    if todos.entry[0].title != 'foo'
        raise Exception("Test assertGetEntries failed: Entry.0.Title, got: " + todos.entry[0].title);

    if todos.entry[1].title() != "bar"
        raise Exception("Test assertGetEntries failed: Entry.1.Title, got: " + todos.entry[1].title);
}

def assertInsert(client: Client)
{
    payload = Todo()
    payload.title = "baz"

    message = client.test().insert(payload);

    if message.success != true
        raise Exception("Test assertInsert failed: Success, got: " + message.success);

    if message.message != "Successful"
        raise Exception("Test assertInsert failed: Message, got: " + message.message);
}

def assertThrowException(client: Client)
{
    try:
        client.test().throwException();

        raise Exception("Test assertThrowException failed: Expected an error");
    except Exception as e:
        if e.payload.message != "Error"
            raise Exception("Test assertThrowException failed: Error message does not match, got: " + e.payload.message);
}

