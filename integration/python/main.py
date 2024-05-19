from sdkgen import Anonymous

from sdk.client import Client
from sdk.error_exception import ErrorException
from sdk.todo import Todo

def assert_get_hello(client: Client):
    message = client.test().get_hello()

    if message.message != 'Hello World!':
        raise Exception('Test assertGetHello failed: Message, got: ' + message.message)


def assert_get_entries(client: Client):
    todos = client.test().get_entries(0, 16)

    if todos.total_results != 4:
        raise Exception("Test assertGetEntries failed: TotalResults, got: " + str(todos.total_results))

    if todos.start_index != 0:
        raise Exception("Test assertGetEntries failed: StartIndex, got: " + str(todos.start_index))

    if todos.items_per_page != 16:
        raise Exception("Test assertGetEntries failed: ItemsPerPage, got: " + str(todos.items_per_page))

    if len(todos.entry) != 2:
        raise Exception("Test assertGetEntries failed: Entry, got: " + str(len(todos.entry)))

    if todos.entry[0].title != 'foo':
        raise Exception("Test assertGetEntries failed: Entry.0.Title, got: " + todos.entry[0].title)

    if todos.entry[1].title != "bar":
        raise Exception("Test assertGetEntries failed: Entry.1.Title, got: " + todos.entry[1].title)


def assert_insert(client: Client):
    payload = Todo("baz")

    message = client.test().insert(payload)

    if not message.success:
        raise Exception("Test assertInsert failed: Success, got: " + str(message.success))

    if message.message != "Successful":
        raise Exception("Test assertInsert failed: Message, got: " + message.message)


def assert_throw_exception(client: Client):
    try:
        client.test().throw_exception()

        raise Exception("Test assertThrowException failed: Expected an error")
    except ErrorException as e:
        if e.get_payload().message != "Error":
            raise Exception("Test assertThrowException failed: Error message does not match, got: " + e.payload.message)


credentials = Anonymous()
client = Client('http://127.0.0.1:1080', credentials)

assert_get_hello(client)
assert_get_entries(client)
assert_insert(client)
assert_throw_exception(client)
