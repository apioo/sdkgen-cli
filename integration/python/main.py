from sdkgen import Anonymous

from sdk.client import Client
from sdk.error_exception import ErrorException
from sdk.todo import Todo

def assert_get_hello(client: Client):
    message = client.test().get_hello()

    if message.message != 'Hello World!':
        raise Exception('Test assert_get_hello failed: Message, got: ' + message.message)

    if message.category != 'default':
        raise Exception('Test assert_get_hello failed: Category, got: ' + message.category)

    if message.priority != 7:
        raise Exception('Test assert_get_hello failed: Priority, got: ' + message.priority)

    if message.weight != 13.37:
        raise Exception('Test assert_get_hello failed: Weight, got: ' + message.weight)

    if message.done != true:
        raise Exception('Test assert_get_hello failed: Done, got: ' + message.done)


def assert_get_entries(client: Client):
    todos = client.test().get_entries(0, 16)

    if todos.total_results != 4:
        raise Exception("Test assert_get_entries failed: TotalResults, got: " + str(todos.total_results))

    if todos.start_index != 0:
        raise Exception("Test assert_get_entries failed: StartIndex, got: " + str(todos.start_index))

    if todos.items_per_page != 16:
        raise Exception("Test assert_get_entries failed: ItemsPerPage, got: " + str(todos.items_per_page))

    if len(todos.entry) != 2:
        raise Exception("Test assert_get_entries failed: Entry, got: " + str(len(todos.entry)))

    if todos.entry[0].title != 'foo':
        raise Exception("Test assert_get_entries failed: Entry.0.Title, got: " + todos.entry[0].title)

    if todos.entry[1].title != "bar":
        raise Exception("Test assert_get_entries failed: Entry.1.Title, got: " + todos.entry[1].title)


def assert_insert(client: Client):
    payload = Todo()
    payload.title = "baz"

    message = client.test().insert(payload)

    if not message.success:
        raise Exception("Test assert_insert failed: Success, got: " + str(message.success))

    if message.message != "Successful":
        raise Exception("Test assert_insert failed: Message, got: " + message.message)


def assert_throw_exception(client: Client):
    try:
        client.test().throw_exception()

        raise Exception("Test assert_throw_exception failed: Expected an error")
    except ErrorException as e:
        if e.get_payload().message != "Error":
            raise Exception("Test assert_throw_exception failed: Error message does not match, got: " + e.payload.message)


def assert_binary(client: Client):
    payload = bytes([0x66, 0x6F, 0x6F, 0x62, 0x61, 0x72])

    response = client.test().binary(payload)

    if payload != response:
        raise Exception("Test assert_binary failed")


def assert_form(client: Client):
    payload = {
        "foo": "bar"
    }

    response = client.test().form(payload)

    if payload["foo"] != response["foo"][0]:
        raise Exception("Test assert_form failed")


def assert_json(client: Client):
    payload = {
        "foo": "bar"
    }

    response = client.test().json(payload)

    if payload != response:
        raise Exception("Test assert_json failed")


def assert_text(client: Client):
    payload = "foobar"

    response = client.test().text(payload)

    if payload != response:
        raise Exception("Test assert_text failed")


def assert_xml(client: Client):
    payload = "<foo>bar</foo>"

    response = client.test().xml(payload)

    if payload != response:
        raise Exception("Test assert_xml failed")


credentials = Anonymous()
client = Client('http://127.0.0.1:1080', credentials)

assert_get_hello(client)
assert_get_entries(client)
assert_insert(client)
assert_throw_exception(client)
assert_binary(client)
assert_form(client)
assert_json(client)
assert_text(client)
assert_xml(client)
