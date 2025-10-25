import {Anonymous} from "sdkgen-client";
import {Client} from "./Client";
import {Todo} from "./Todo";
import {ErrorException} from "./ErrorException";

async function main() {
    const credentials = new Anonymous();
    const client = new Client('http://127.0.0.1:1080', credentials);

    await assertGetHello(client);
    await assertGetEntries(client);
    await assertInsert(client);
    await assertThrowException(client);
    await assertBinary(client);
    await assertForm(client);
    await assertJson(client);
    await assertText(client);
    await assertXml(client);
}

async function assertGetHello(client: Client) {
    const message = await client.test().getHello();

    if (message.message != 'Hello World!') {
        throw new Error('Test assertGetHello failed: Message, got: ' + message.message);
    }

    // currently for TS we use interfaces and therefor we cant set default values
    /*
    if (message.category != 'default') {
        throw new Error("Test assertGetHello failed: Category, got: " + message.category);
    }
    */

    if (message.priority != 7) {
        throw new Error("Test assertGetHello failed: Priority, got: " + message.priority);
    }

    if (message.weight != 13.37) {
        throw new Error("Test assertGetHello failed: Weight, got: " + message.weight);
    }

    if (message.done != true) {
        throw new Error("Test assertGetHello failed: Done, got: " + message.done);
    }
}

async function assertGetEntries(client: Client) {
    const todos = await client.test().getEntries(0, 16);

    if (todos.totalResults != 4) {
        throw new Error("Test assertGetEntries failed: TotalResults, got: " + todos.totalResults);
    }

    if (todos.startIndex != 0) {
        throw new Error("Test assertGetEntries failed: StartIndex, got: " + todos.startIndex);
    }

    if (todos.itemsPerPage != 16) {
        throw new Error("Test assertGetEntries failed: ItemsPerPage, got: " + todos.itemsPerPage);
    }

    if (todos.entry?.length != 2) {
        throw new Error("Test assertGetEntries failed: Entry, got: " + todos.entry?.length);
    }

    if (todos.entry[0].title != 'foo') {
        throw new Error("Test assertGetEntries failed: Entry.0.Title, got: " + todos.entry[0].title);
    }

    if (todos.entry[1].title != "bar") {
        throw new Error("Test assertGetEntries failed: Entry.1.Title, got: " + todos.entry[1].title);
    }
}

async function assertInsert(client: Client) {
    const payload: Todo = {
        title: "baz"
    };

    const message = await client.test().insert(payload);

    if (!message.success) {
        throw new Error("Test assertInsert failed: Success, got: " + message.success);
    }

    if (message.message != "Successful") {
        throw new Error("Test assertInsert failed: Message, got: " + message.message);
    }
}

async function assertThrowException(client: Client) {
    try {
        await client.test().throwException();

        throw new Error("Test assertThrowException failed: Expected an error");
    } catch (e) {
        if (e instanceof ErrorException && e.getPayload().message != "Error") {
            throw new Error("Test assertThrowException failed: Error message does not match, got: " + e.getPayload().message);
        }
    }
}

async function assertBinary(client: Client) {
    const decoder = new TextDecoder();
    const payload = newArrayBuffer();

    const response = await client.test().binary(payload);

    if (decoder.decode(response) !== decoder.decode(payload)) {
        throw new Error("Test assertBinary failed");
    }
}

async function assertForm(client: Client) {
    const payload = new URLSearchParams({foo: 'bar'});

    const response = await client.test().form(payload);

    if (payload.toString() !== response.toString()) {
        throw new Error("Test assertForm failed");
    }
}

async function assertJson(client: Client) {
    const payload = {foo: 'bar'};

    const response = await client.test().json(payload);

    if (JSON.stringify(payload) !== JSON.stringify(response)) {
        throw new Error("Test assertJson failed");
    }
}

async function assertText(client: Client) {
    const payload = 'foobar';

    const response = await client.test().text(payload);

    if (payload !== response) {
        throw new Error("Test assertText failed");
    }
}

async function assertXml(client: Client) {
    const payload = '<foo>bar</foo>';

    const response = await client.test().xml(payload);

    if (payload !== response) {
        throw new Error("Test assertText failed");
    }
}

function newArrayBuffer(): ArrayBuffer {
    const data = Buffer.from('Zm9vYmFy', 'base64');
    return data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength);
}

main();
