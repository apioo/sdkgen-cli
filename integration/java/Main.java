package org.typeapi.generator;

import app.sdkgen.client.Credentials.HttpBearer;
import app.sdkgen.client.Exception.ClientException;

public class Main {
    public static void main(String[] args) throws ClientException {
        HttpBearer credentials = new HttpBearer("foo");
        Client client = new Client("http://localhost:1080", credentials);

        assertGetHello(client);
        assertGetEntries(client);
        assertInsert(client);
        assertThrowException(client);
    }

    private static void assertGetHello(Client client) throws ClientException {
        HelloWorld message = client.test().getHello();

        if (!message.getMessage().equals("Hello World!")) {
            throw new RuntimeException("Test assertGetHello failed: Message, got: " + message.getMessage());
        }
    }

    private static void assertGetEntries(Client client) throws ClientException {
        Todos todos = client.test().getEntries(0, 16);

        if (todos.getTotalResults() != 4) {
            throw new RuntimeException("Test assertGetEntries failed: TotalResults, got: " + todos.getTotalResults());
        }

        if (todos.getStartIndex() != 0) {
            throw new RuntimeException("Test assertGetEntries failed: StartIndex, got: " + todos.getStartIndex());
        }

        if (todos.getItemsPerPage() != 16) {
            throw new RuntimeException("Test assertGetEntries failed: ItemsPerPage, got: " + todos.getItemsPerPage());
        }

        if (todos.getEntry().length != 2) {
            throw new RuntimeException("Test assertGetEntries failed: Entry, got: " + todos.getEntry().length);
        }

        if (!todos.getEntry()[0].getTitle().equals("foo")) {
            throw new RuntimeException("Test assertGetEntries failed: Entry.0.Title, got: " + todos.getEntry()[0].getTitle());
        }

        if (!todos.getEntry()[1].getTitle().equals("bar")) {
            throw new RuntimeException("Test assertGetEntries failed: Entry.1.Title, got: " + todos.getEntry()[1].getTitle());
        }
    }

    private static void assertInsert(Client client) throws ClientException {
        Todo payload = new Todo();
        payload.setTitle("baz");

        Message message = client.test().insert(payload);

        if (!message.getSuccess()) {
            throw new RuntimeException("Test assertInsert failed: Success, got: " + message.getSuccess());
        }

        if (!message.getMessage().equals("Successful")) {
            throw new RuntimeException("Test assertInsert failed: Message, got: " + message.getMessage());
        }
    }

    private static void assertThrowException(Client client) throws ClientException
    {
        try {
            client.test().throwException();

            throw new RuntimeException("Test assertThrowException failed: Expected an error");
        } catch (ErrorException e) {
            if (!e.getPayload().getMessage().equals("Error")) {
                throw new RuntimeException("Test assertThrowException failed: Error message does not match, got: " + e.getPayload().getMessage());
            }
        }
    }
}
