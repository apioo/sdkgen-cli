const path = require('path');
const fs = require('fs');

const projectPath = process.argv[1] || process.cwd();

// 1. Ensure reflect-metadata is required (standard in NestJS)
try {
    require(path.join(projectPath, 'node_modules', 'reflect-metadata'));
} catch (e) {
    try { require('reflect-metadata'); } catch (err) {}
}

// Helper to convert JS/TS types to TypeSchema formats
function tsTypeToTypeSchema(type) {
    if (!type) return { type: 'string' };

    if (type === String) return { type: 'string' };
    if (type === Number) return { type: 'number' };
    if (type === Boolean) return { type: 'boolean' };
    if (type === Array) return { type: 'array', schema: { type: 'string' } };

    if (typeof type === 'function' && type.name && type.name !== 'Object') {
        return { type: 'reference', target: type.name };
    }

    return { type: 'string' };
}

// Extract properties from DTO class
function extractDtoDefinition(dtoClass, definitions) {
    if (!dtoClass || typeof dtoClass !== 'function') return;
    const className = dtoClass.name;
    if (!className || className === 'Object' || definitions[className]) return;

    const instance = new dtoClass();
    const keys = Object.getOwnPropertyNames(instance);

    const properties = {};
    const required = [];

    // Inspect properties using TypeScript reflect metadata if available
    keys.forEach((key) => {
        const propType = Reflect.getMetadata('design:type', dtoClass.prototype, key);
        properties[key] = tsTypeToTypeSchema(propType);
    });

    const definition = {
        type: 'object',
        properties: properties
    };

    definitions[className] = definition;
}

async function bootstrap() {
    // 2. Locate built dist main file or app module
    const distMain = path.join(projectPath, 'dist', 'main.js');
    const distAppModule = path.join(projectPath, 'dist', 'app.module.js');

    let AppModule;
    let NestFactory;
    let core;

    try {
        core = require(path.join(projectPath, 'node_modules', '@nestjs/core'));
        const common = require(path.join(projectPath, 'node_modules', '@nestjs/common'));

        if (fs.existsSync(distAppModule)) {
            AppModule = require(distAppModule).AppModule;
        } else if (fs.existsSync(distMain)) {
            const mainExports = require(distMain);
            AppModule = mainExports.AppModule;
        } else {
            throw new Error(`Could not find compiled dist/app.module.js or dist/main.js in ${projectPath}. Please build the NestJS project first (npm run build).`);
        }

        // Constants used by NestJS internally for Metadata
        const PATH_METADATA = 'path';
        const METHOD_METADATA = 'method';
        const HTTP_CODE_METADATA = '__httpCode__';

        const HTTP_METHODS = {
            0: 'GET',
            1: 'POST',
            2: 'PUT',
            3: 'DELETE',
            4: 'PATCH',
            5: 'ALL',
            6: 'OPTIONS',
            7: 'HEAD',
        };

        // Create application context to inspect controllers without starting HTTP server
        const app = await core.NestFactory.createApplicationContext(AppModule, { logger: false });
        const modulesContainer = app.get(core.ModulesContainer);

        const operations = {};
        const definitions = {};

        modulesContainer.forEach((moduleRef) => {
            moduleRef.controllers.forEach((wrapper) => {
                const { instance, metatype } = wrapper;
                if (!instance || !metatype) return;

                // Controller base path
                const controllerPath = Reflect.getMetadata(PATH_METADATA, metatype) || '';
                const prototype = Object.getPrototypeOf(instance);

                const methodNames = Object.getOwnPropertyNames(prototype).filter(
                    (item) => item !== 'constructor' && typeof prototype[item] === 'function'
                );

                methodNames.forEach((methodName) => {
                    const handler = prototype[methodName];
                    const routePath = Reflect.getMetadata(PATH_METADATA, handler);
                    const methodCode = Reflect.getMetadata(METHOD_METADATA, handler);

                    if (methodCode === undefined) return;

                    const httpMethod = HTTP_METHODS[methodCode] || 'GET';
                    const fullPath = ('/' + controllerPath + '/' + (routePath || '')).replace(/\/+/g, '/');

                    const opId = `${httpMethod.toLowerCase()}_${fullPath.replace(/[\/\{\}:]/g, '_').replace(/_+/g, '_')}`;

                    const operation = {
                        method: httpMethod,
                        path: fullPath,
                        description: '',
                        arguments: {},
                        return: {
                            code: Reflect.getMetadata(HTTP_CODE_METADATA, handler) || (httpMethod === 'POST' ? 201 : 200)
                        }
                    };

                    // Inspect method parameters (DTOs / Payloads)
                    const paramTypes = Reflect.getMetadata('design:paramtypes', prototype, methodName) || [];
                    paramTypes.forEach((paramType) => {
                        if (paramType && typeof paramType === 'function' && paramType.name !== 'Object') {
                            const dtoName = paramType.name;
                            extractDtoDefinition(paramType, definitions);

                            operation.arguments['payload'] = {
                                in: 'body',
                                schema: {
                                    type: 'reference',
                                    target: dtoName
                                }
                            };
                        }
                    });

                    // Inspect return type metadata
                    const returnType = Reflect.getMetadata('design:returntype', prototype, methodName);
                    if (returnType && typeof returnType === 'function' && returnType.name !== 'Object') {
                        const dtoName = returnType.name;
                        extractDtoDefinition(returnType, definitions);

                        operation.return.schema = {
                            type: 'reference',
                            target: dtoName
                        };
                    }

                    operations[opId] = operation;
                });
            });
        });

        await app.close();

        console.log(JSON.stringify({ operations, definitions }, null, 2));
    } catch (err) {
        process.stderr.write(`NestJS Extraction Error: ${err.message}\n`);
        process.exit(1);
    }
}

bootstrap();
