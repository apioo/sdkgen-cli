
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
        await AssertGetFormConfig(client);
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

        if (message.Category != "default") {
            throw new Exception("Test assertGetHello failed: Category, got: " + message.Category);
        }

        if (message.Priority != 7) {
            throw new Exception("Test assertGetHello failed: Priority, got: " + message.Priority);
        }

        if (message.Weight != 13.37) {
            throw new Exception("Test assertGetHello failed: Weight, got: " + message.Weight);
        }

        if (message.Done != true) {
            throw new Exception("Test assertGetHello failed: Done, got: " + message.Done);
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

    private static async Task AssertGetFormConfig(Client client)
    {
        CommonFormContainer form = await client.Test().GetFormConfig();

        if (form.Elements.Count != 3) {
            throw new Exception("Test assertForm failed: Elements, got: " + form.Elements.Count);
        }

        if (form.Elements[0].GetType().Name != "CommonFormElementInput") {
            throw new Exception("Test assertForm failed: Elements.0, got: " + form.Elements[0].GetType().Name);
        }

        CommonFormElementInput input = (CommonFormElementInput) form.Elements[0];
        if (input.Element != "input") {
            throw new Exception("Test assertForm failed: Elements.0.Element, got: " + input.Element);
        }

        if (input.Type != "text") {
            throw new Exception("Test assertForm failed: Elements.0.Type, got: " + input.Type);
        }

        if (form.Elements[1].GetType().Name != "CommonFormElementSelect") {
            throw new Exception("Test assertForm failed: Elements.1, got: " + form.Elements[1].GetType().Name);
        }

        CommonFormElementSelect select = (CommonFormElementSelect) form.Elements[1];
        if (select.Element != "select") {
            throw new Exception("Test assertForm failed: Elements.1.Element, got: " + select.Element);
        }

        if (select.Options.Count != 2) {
            throw new Exception("Test assertForm failed: Elements.1.Options, got: " + select.Options.Count);
        }

        if (form.Elements[2].GetType().Name != "CommonFormElementTextArea") {
            throw new Exception("Test assertForm failed: Elements.2, got: " + form.Elements[2].GetType().Name);
        }

        CommonFormElementTextArea textArea = (CommonFormElementTextArea) form.Elements[2];
        if (textArea.Element != "textarea") {
            throw new Exception("Test assertForm failed: Elements.2.Element, got: " + textArea.Element);
        }
    }

    private static async Task AssertBinary(Client client)
    {
        var payload = new byte[] {0x66, 0x6F, 0x6F, 0x62, 0x61, 0x72};

        var response = await client.Test().Binary(payload);

        if (!payload.SequenceEqual(response)) {
            throw new Exception("Test AssertBinary failed");
        }
    }

    private static async Task AssertForm(Client client)
    {
        var payload = new System.Collections.Specialized.NameValueCollection();
        payload.Add("foo", "bar");

        var response = await client.Test().Form(payload);

        if (payload.Get("foo") != response.Get("foo")) {
            throw new Exception("Test AssertForm failed");
        }
    }

    private static async Task AssertJson(Client client)
    {
        var payload = new Dictionary<string, string>();
        payload.Add("foo", "bar");

        var response = await client.Test().Json(payload);

        if (System.Text.Json.JsonSerializer.Serialize(payload) != System.Text.Json.JsonSerializer.Serialize(response)) {
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
