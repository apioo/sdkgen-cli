
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
        await AssertBinary(client);
        await AssertForm(client);
        await AssertJson(client);
        await AssertText(client);
        await AssertXml(client);
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

    private static async Task AssertBinary(Client client)
    {
        var payload = new byte[] {0x66, 0x6F, 0x6F, 0x62, 0x61, 0x72};

        var response = await client.Test().Binary(payload);

        if (payload != response) {
            throw new Exception("Test AssertBinary failed");
        }
    }

    private static async Task AssertForm(Client client)
    {
        var payload = new System.Collections.Specialized.NameValueCollection
        {
            { "foo", "bar" }
        };

        var response = await client.Test().Form(payload);

        if (payload != response) {
            throw new Exception("Test AssertForm failed");
        }
    }

    private static async Task AssertJson(Client client)
    {
        var payload = new Dictionary<string, string>
        {
            { "string", "bar" }
        };

        var response = await client.Test().Json(payload);

        if (payload != response) {
            throw new Exception("Test AssertJson failed");
        }
    }

    private static async Task AssertText(Client client)
    {
        var payload = "foobar";

        var response = await client.Test().Text(payload);

        if (payload != response) {
            throw new Exception("Test AssertText failed");
        }
    }

    private static async Task AssertXml(Client client)
    {
        var payload = "<foo>bar</foo>";

        var response = await client.Test().Xml(payload);

        if (payload != response) {
            throw new Exception("Test AssertXml failed");
        }
    }
}
