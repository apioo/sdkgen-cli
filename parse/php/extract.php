<?php

use Symfony\Component\Routing\Annotation\Route as LegacyRoute;
use Symfony\Component\Routing\Attribute\Route;
use Symfony\HttpKernel\Attribute\MapRequestPayload;

$projectPath = $argv[1] ?? getcwd();

// 1. Locate and require Composer Autoloader
$autoloadFile = $projectPath . '/vendor/autoload.php';
if (!file_exists($autoloadFile)) {
    fwrite(STDERR, "Error: vendor/autoload.php not found in {$projectPath}\n");
    exit(1);
}
require_once $autoloadFile;

// Helper to convert PHP native types to TypeSchema formats
function phpTypeToTypeSchema(?ReflectionType $type): array {
    if (!$type) {
        return ["type" => "string"];
    }

    $typeName = $type instanceof ReflectionNamedType ? $type->getName() : (string) $type;

    return match ($typeName) {
        'int' => ["type" => "integer"],
        'float' => ["type" => "number"],
        'bool' => ["type" => "boolean"],
        'string' => ["type" => "string"],
        'array' => ["type" => "array", "schema" => ["type" => "string"]],
        default => class_exists($typeName)
            ? ["type" => "reference", "target" => (new ReflectionClass($typeName))->getShortName()]
            : ["type" => "string"],
    };
}

// Helper to inspect DTO classes mapped via #[MapRequestPayload]
function extractDtoDefinition(string $class, array &$definitions): void {
    $refClass = new ReflectionClass($class);
    $shortName = $refClass->getShortName();

    if (isset($definitions[$shortName])) {
        return;
    }

    $properties = [];
    $required = [];

    foreach ($refClass->getProperties(ReflectionProperty::IS_PUBLIC) as $property) {
        $name = $property->getName();
        $properties[$name] = phpTypeToTypeSchema($property->getType());

        if ($property->getType() && !$property->getType()->allowsNull() && !$property->hasDefaultValue()) {
            $required[] = $name;
        }
    }

    $definition = [
        "type" => "object",
        "properties" => $properties
    ];

    if (!empty($required)) {
        $definition["required"] = $required;
    }

    $definitions[$shortName] = $definition;
}

// 2. Boot Kernel to get the Router
try {
    $kernelFile = $projectPath . '/src/Kernel.php';
    if (!file_exists($kernelFile)) {
        fwrite(STDERR, "Error: src/Kernel.php not found\n");
        exit(1);
    }
    require_once $kernelFile;

    $env = $_SERVER['APP_ENV'] ?? 'dev';
    $debug = (bool) ($_SERVER['APP_DEBUG'] ?? true);

    $kernel = new App\Kernel($env, $debug);
    $kernel->boot();

    $container = $kernel->getContainer();
    $router = $container->get('router');
    $routeCollection = $router->getRouteCollection();
} catch (\Throwable $e) {
    fwrite(STDERR, "Failed to boot Symfony application: " . $e->getMessage() . "\n");
    exit(1);
}

$operations = [];
$definitions = [];

// 3. Process Symfony Routes
foreach ($routeCollection->all() as $name => $route) {
    // Skip internal Symfony/_profiler routes
    if (str_starts_with($name, '_')) {
        continue;
    }

    $controller = $route->getDefault('_controller');
    if (!$controller || !is_string($controller) || !str_contains($controller, '::')) {
        continue;
    }

    [$controllerClass, $methodName] = explode('::', $controller);

    if (!class_exists($controllerClass) || !method_exists($controllerClass, $methodName)) {
        continue;
    }

    $refMethod = new ReflectionMethod($controllerClass, $methodName);
    $methods = $route->getMethods() ?: ['GET'];

    foreach ($methods as $httpMethod) {
        $operationId = strtolower($httpMethod) . '_' . str_replace(['/', '{', '}'], ['_', '', ''], $route->getPath());

        $operation = [
            "method" => strtoupper($httpMethod),
            "path" => $route.getPath(),
            "description" => "",
            "arguments" => [],
            "return" => [
                "code" => 200
            ]
        ];

        // Process Method Parameters (Query, Path, and Body Payloads)
        foreach ($refMethod->getParameters() as $param) {
            $paramName = $param->getName();
            $paramType = $param->getType();

            // Check for #[MapRequestPayload] attribute (Symfony 6.3+)
            $payloadAttrs = $param->getAttributes(MapRequestPayload::class);
            if (!empty($payloadAttrs) && $paramType && class_exists($paramType->getName())) {
                $dtoClass = $paramType->getName();
                $shortName = (new ReflectionClass($dtoClass))->getShortName();

                extractDtoDefinition($dtoClass, $definitions);

                $operation["arguments"]["payload"] = [
                    "in" => "body",
                    "schema" => [
                        "type" => "reference",
                        "target" => $shortName
                    ]
                ];
                continue;
            }

            // Determine if parameter is path or query
            $in = str_contains($route->getPath(), '{' . $paramName . '}') ? 'path' : 'query';
            $operation["arguments"][$paramName] = [
                "in" => $in,
                "schema" => phpTypeToTypeSchema($paramType)
            ];
        }

        // Process Return Type (if typed controller method or DTO return)
        $returnType = $refMethod->getReturnType();
        if ($returnType && $returnType instanceof ReflectionNamedType && class_exists($returnType->getName())) {
            $dtoClass = $returnType->getName();
            $shortName = (new ReflectionClass($dtoClass))->getShortName();

            extractDtoDefinition($dtoClass, $definitions);

            $operation["return"]["schema"] = [
                "type" => "reference",
                "target" => $shortName
            ];
        }

        $operations[$operationId] = $operation;
    }
}

// Output final TypeAPI JSON
echo json_encode([
    "operations" => $operations,
    "definitions" => $definitions
], JSON_PRETTY_PRINT | JSON_UNESCAPED_SLASHES);
