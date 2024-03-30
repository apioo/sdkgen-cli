
using Sdkgen.Client.Credentials;

class Program
{
    static async Task Main()
    {
        Anonymous credentials = new Anonymous();
        Client client = new Client("http://127.0.0.1:1080", credentials);

        await AssertGetHello(client);
        await AssertGetEntries(client);
        await AssertInsert(client);
        await AssertThrowException(client);
    }

    private static async Task AssertGetHello(Client client)
    {
        HelloWorld message = await client.Test().GetHello();

        if (message.Message != "Hello World!") {
            throw new Exception("Test assertGetHello failed: Message, got: " + message.Message);
        }
    }

    private static async Task AssertGetEntries(Client client)
    {
        Todos todos = await client.Test().GetEntries(0, 16);

        if (todos.TotalResults != 4) {
            throw new Exception("Test assertGetEntries failed: TotalResults, got: " + todos.TotalResults);
        }

        if (todos.StartIndex != 0) {
            throw new Exception("Test assertGetEntries failed: StartIndex, got: " + todos.StartIndex);
        }

        if (todos.ItemsPerPage != 16) {
            throw new Exception("Test assertGetEntries failed: ItemsPerPage, got: " + todos.ItemsPerPage);
        }

        if (todos.Entry.Count != 2) {
            throw new Exception("Test assertGetEntries failed: Entry, got: " + todos.Entry.Count);
        }

        if (todos.Entry[0].Title != "foo") {
            throw new Exception("Test assertGetEntries failed: Entry.0.Title, got: " + todos.Entry[0].Title);
        }

        if (todos.Entry[1].Title != "bar") {
            throw new Exception("Test assertGetEntries failed: Entry.1.Title, got: " + todos.Entry[1].Title);
        }
    }

    private static async Task AssertInsert(Client client)
    {
        Todo payload = new Todo();
        payload.Title = "baz";

        Response message = await client.Test().Insert(payload);

        if (!(message.Success ?? false)) {
            throw new Exception("Test assertInsert failed: Success, got: " + message.Success);
        }

        if (message.Message != "Successful") {
            throw new Exception("Test assertInsert failed: Message, got: " + message.Message);
        }
    }

    private static async Task AssertThrowException(Client client)
    {
        try {
            await client.Test().ThrowException();

            throw new Exception("Test assertThrowException failed: Expected an error");
        } catch (ErrorException e) {
            if (e.Payload.Message != "Error") {
                throw new Exception("Test assertThrowException failed: Error message does not match, got: " + e.Payload.Message);
            }
        }
    }
}
