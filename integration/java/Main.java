package org.typeapi.generator;

import app.sdkgen.client.Credentials.Anonymous;
import app.sdkgen.client.Exception.ClientException;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.apache.hc.core5.http.NameValuePair;
import org.apache.hc.core5.http.message.BasicNameValuePair;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;

public class Main {
    public static void main(String[] args) throws ClientException, JsonProcessingException {
        Anonymous credentials = new Anonymous();
        Client client = new Client("http://127.0.0.1:1080", credentials);

        assertGetHello(client);
        assertGetEntries(client);
        assertInsert(client);
        assertThrowException(client);
        assertGetFormConfig(client);
        assertBinary(client);
        assertForm(client);
        assertJson(client);
        assertText(client);
        assertXml(client);
    }

    private static void assertGetHello(Client client) throws ClientException {
        HelloWorld message = client.test().getHello();

        if (!message.getMessage().equals("Hello World!")) {
            throw new RuntimeException("Test assertGetHello failed: Message, got: " + message.getMessage());
        }

        if (!message.getCategory().equals("default")) {
            throw new RuntimeException("Test assertGetHello failed: Category, got: " + message.getCategory());
        }

        if (message.getPriority() != 7) {
            throw new RuntimeException("Test assertGetHello failed: Priority, got: " + message.getPriority());
        }

        if (message.getWeight() != 13.37) {
            throw new RuntimeException("Test assertGetHello failed: Weight, got: " + message.getWeight());
        }

        if (message.getDone() != true) {
            throw new RuntimeException("Test assertGetHello failed: Done, got: " + message.getDone());
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

        if (todos.getEntry().size() != 2) {
            throw new RuntimeException("Test assertGetEntries failed: Entry, got: " + todos.getEntry().size());
        }

        if (!todos.getEntry().get(0).getTitle().equals("foo")) {
            throw new RuntimeException("Test assertGetEntries failed: Entry.0.Title, got: " + todos.getEntry().get(0).getTitle());
        }

        if (!todos.getEntry().get(1).getTitle().equals("bar")) {
            throw new RuntimeException("Test assertGetEntries failed: Entry.1.Title, got: " + todos.getEntry().get(1).getTitle());
        }
    }

    private static void assertInsert(Client client) throws ClientException {
        Todo payload = new Todo();
        payload.setTitle("baz");

        Response message = client.test().insert(payload);

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

    private static void assertGetFormConfig(Client client) throws ClientException
    {
        CommonFormContainer form = await client.Test().GetFormConfig();

        if (form.getElements().size() != 3) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements, got: " + form.getElements().size();
        }

        if (!form.getElements().get(0) instanceof CommonFormElementInput) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements.0, got: " + form.getElements().get(0).getClass().getName());
        }

        CommonFormElementInput input = (CommonFormElementInput) form.getElements().get(0);
        if (!input.getType().equals("text")) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements.0.Type, got: " + input.getType());
        }

        if (!form.getElements().get(1) instanceof CommonFormElementSelect) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements.1, got: " + form.getElements().get(1).getClass().getName());
        }

        CommonFormElementSelect select = (CommonFormElementSelect) form.getElements().get(1);
        if (select.getOptions().size() != 2) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements.1.Options, got: " + select.getOptions().size());
        }

        if (!form.getElements().get(2) instanceof CommonFormElementTextArea) {
            throw new RuntimeException("Test assertGetFormConfig failed: Elements.2, got: " + form.getElements().get(2).getClass().getName());
        }
    }

    private static void assertBinary(Client client) throws ClientException
    {
        byte[] payload = {0x66, 0x6F, 0x6F, 0x62, 0x61, 0x72};

        var response = client.test().binary(payload);

        if (!Arrays.equals(payload, response)) {
            throw new RuntimeException("Test assertBinary failed");
        }
    }

    private static void assertForm(Client client) throws ClientException
    {
        List<NameValuePair> payload = List.of(new BasicNameValuePair("foo", "bar"));

        var response = client.test().form(payload);

        if (!response.get(0).getName().equals("foo") || !response.get(0).getValue().equals("bar")) {
            throw new RuntimeException("Test assertForm failed");
        }
    }

    private static void assertJson(Client client) throws ClientException, JsonProcessingException {
        HashMap<String, String> payload = new HashMap<>();
        payload.put("foo", "bar");

        var response = client.test().json(payload);

        var objectMapper = new ObjectMapper();
        if (!objectMapper.writeValueAsString(payload).equals(objectMapper.writeValueAsString(response))) {
            throw new RuntimeException("Test assertJson failed");
        }
    }

    private static void assertText(Client client) throws ClientException
    {
        String payload = "foobar";

        var response = client.test().text(payload);

        if (!payload.equals(response)) {
            throw new RuntimeException("Test assertText failed");
        }
    }

    private static void assertXml(Client client) throws ClientException
    {
        String payload = "<foo>bar</foo>";

        var response = client.test().xml(payload);

        if (!payload.equals(response)) {
            throw new RuntimeException("Test assertXml failed");
        }
    }
}
