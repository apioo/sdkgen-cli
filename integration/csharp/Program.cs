
using Sdkgen.Client.Credentials;

class Program
{
    static void Main()
    {
        Anonymous credentials = new Anonymous();
        Client client = new Client("http://127.0.0.1:1080", credentials);

        AssertGetHello(client);
        AssertGetEntries(client);
        AssertInsert(client);
        AssertThrowException(client);
    }

    private static void AssertGetHello(Client client)
    {
        HelloWorld message = client.Test().GetHello();

        if (message.Message != "Hello World!") {
            throw new Exception("Test assertGetHello failed: Message, got: " + message.Message);
        }
    }

    private static void AssertGetEntries(Client client)
    {
        Todos todos = client.Test().GetEntries(0, 16);

        if (todos.TotalResults != 4) {
            throw new Exception("Test assertGetEntries failed: TotalResults, got: " + todos.TotalResults);
        }

        if (todos.StartIndex != 0) {
            throw new Exception("Test assertGetEntries failed: StartIndex, got: " + todos.StartIndex);
        }

        if (todos.ItemsPerPage != 16) {
            throw new Exception("Test assertGetEntries failed: ItemsPerPage, got: " + todos.ItemsPerPage);
        }

        if (todos.Entry.Length != 2) {
            throw new Exception("Test assertGetEntries failed: Entry, got: " + todos.Entry.Length);
        }

        if (todos.Entry[0].Title != "foo") {
            throw new Exception("Test assertGetEntries failed: Entry.0.Title, got: " + todos.Entry[0].Title);
        }

        if (todos.Entry[1].Title != "bar") {
            throw new Exception("Test assertGetEntries failed: Entry.1.Title, got: " + todos.Entry[1].Title);
        }
    }

    private static void AssertInsert(Client client)
    {
        Todo payload = new Todo();
        payload.Title = "baz";

        Message message = client.Test().Insert(payload);

        if (!message.Success) {
            throw new Exception("Test assertInsert failed: Success, got: " + message.Success);
        }

        if (message.Message != "Successful") {
            throw new Exception("Test assertInsert failed: Message, got: " + message.Message);
        }
    }

    private static void AssertThrowException(Client client)
    {
        try {
            client.Test().ThrowException();

            throw new Exception("Test assertThrowException failed: Expected an error");
        } catch (ErrorException e) {
            if (e.Payload.Message != "Error") {
                throw new Exception("Test assertThrowException failed: Error message does not match, got: " + e.Payload.Message);
            }
        }
    }
}
