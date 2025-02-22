<?php

require_once __DIR__ . '/vendor/autoload.php';

require_once __DIR__ . '/TestTag.php';
require_once __DIR__ . '/Error.php';
require_once __DIR__ . '/ErrorException.php';
require_once __DIR__ . '/HelloWorld.php';
require_once __DIR__ . '/Response.php';
require_once __DIR__ . '/Todo.php';
require_once __DIR__ . '/Todos.php';
require_once __DIR__ . '/Client.php';

$credentials = new \Sdkgen\Client\Credentials\Anonymous();
$client = new \SDK\Client('http://127.0.0.1:1080', $credentials);

assertGetHello($client);
assertGetEntries($client);
assertInsert($client);
assertThrowException($client);
assertBinary($client);
assertForm($client);
assertJson($client);
assertText($client);
assertXml($client);

function assertGetHello(\SDK\Client $client): void
{
    $message = $client->test()->getHello();

    if ($message->getMessage() != 'Hello World!') {
        throw new RuntimeException('Test assertGetHello failed: Message, got: ' . $message->getMessage());
    }
}

function assertGetEntries(\SDK\Client $client): void
{
    $todos = $client->test()->getEntries(0, 16);

    if ($todos->getTotalResults() != 4) {
        throw new RuntimeException("Test assertGetEntries failed: TotalResults, got: " . $todos->getTotalResults());
    }

    if ($todos->getStartIndex() != 0) {
        throw new RuntimeException("Test assertGetEntries failed: StartIndex, got: " . $todos->getStartIndex());
    }

    if ($todos->getItemsPerPage() != 16) {
        throw new RuntimeException("Test assertGetEntries failed: ItemsPerPage, got: " . $todos->getItemsPerPage());
    }

    if (count($todos->getEntry()) != 2) {
        throw new RuntimeException("Test assertGetEntries failed: Entry, got: " . count($todos->getEntry()));
    }

    if ($todos->getEntry()[0]->getTitle() != 'foo') {
        throw new RuntimeException("Test assertGetEntries failed: Entry.0.Title, got: " . $todos->getEntry()[0]->getTitle());
    }

    if ($todos->getEntry()[1]->getTitle() != "bar") {
        throw new RuntimeException("Test assertGetEntries failed: Entry.1.Title, got: " . $todos->getEntry()[1]->getTitle());
    }
}

function assertInsert(\SDK\Client $client): void
{
    $payload = new \SDK\Todo();
    $payload->setTitle("baz");

    $message = $client->test()->insert($payload);

    if (!$message->getSuccess()) {
        throw new RuntimeException("Test assertInsert failed: Success, got: " . $message->getSuccess());
    }

    if ($message->getMessage() != "Successful") {
        throw new RuntimeException("Test assertInsert failed: Message, got: " . $message->getMessage());
    }
}

function assertThrowException(\SDK\Client $client): void
{
    try {
        $client->test()->throwException();

        throw new RuntimeException("Test assertThrowException failed: Expected an error");
    } catch (\SDK\ErrorException $e) {
        if ($e->getPayload()->getMessage() != "Error") {
            throw new RuntimeException("Test assertThrowException failed: Error message does not match, got: " . $e->getPayload()->getMessage());
        }
    }
}

function assertBinary(\SDK\Client $client): void
{
    $handle = fopen('php://memory', 'w+');
    fwrite($handle, 'foobar');
    fseek($handle, 0);

    $payload = new \GuzzleHttp\Psr7\Stream($handle);

    $response = $client->test()->binary($payload);

    if ($response->getContents() !== 'foobar') {
        throw new RuntimeException("Test assertBinary failed");
    }
}

function assertForm(\SDK\Client $client): void
{
    $payload = ['foo' => 'bar'];

    $response = $client->test()->form($payload);

    if ($payload !== $response) {
        throw new RuntimeException("Test assertForm failed");
    }
}

function assertJson(\SDK\Client $client): void
{
    $payload = ['foo' => 'bar'];

    $response = $client->test()->json($payload);

    if (\json_encode($payload) !== \json_encode($response)) {
        throw new RuntimeException("Test assertJson failed");
    }
}

function assertText(\SDK\Client $client): void
{
    $payload = 'foobar';

    $response = $client->test()->text($payload);

    if ($payload !== $response) {
        throw new RuntimeException("Test assertText failed");
    }
}

function assertXml(\SDK\Client $client): void
{
    $payload = '<foo>bar</foo>';

    $response = $client->test()->xml($payload);

    if ($payload !== $response) {
        throw new RuntimeException("Test assertText failed");
    }
}
