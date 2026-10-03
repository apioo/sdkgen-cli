import sys
import json
import inspect
import importlib
from typing import get_args, get_origin, Union
from fastapi import FastAPI
from pydantic import BaseModel

def pydantic_type_to_typeschema(annotation):
    """Translates Python / Pydantic types into TypeSchema types."""
    origin = get_origin(annotation)
    args = get_args(annotation)

    # Unwrap Optional[T] or Union[T, None]
    if origin is Union:
        non_none_args = [a for a in args if a is not type(None)]
        if len(non_none_args) == 1:
            return pydantic_type_to_typeschema(non_none_args[0])

    if annotation == str:
        return {"type": "string"}
    elif annotation == int:
        return {"type": "integer"}
    elif annotation == float:
        return {"type": "number"}
    elif annotation == bool:
        return {"type": "boolean"}
    elif origin == list:
        item_schema = pydantic_type_to_typeschema(args[0]) if args else {"type": "string"}
        return {"type": "array", "schema": item_schema}
    elif inspect.isclass(annotation) and issubclass(annotation, BaseModel):
        return {"type": "reference", "target": annotation.__name__}

    return {"type": "string"}

def extract_pydantic_definition(model_cls: type[BaseModel]):
    """Generates a TypeSchema Object definition from a Pydantic Model."""
    properties = {}
    required = []

    for name, field in model_cls.model_fields.items():
        field_type = field.annotation
        properties[name] = pydantic_type_to_typeschema(field_type)
        if field.is_required():
            required.append(name)

    spec = {
        "type": "object",
        "properties": properties
    }
    if required:
        spec["required"] = required
    return spec

def extract_typeapi(app: FastAPI) -> dict:
    operations = {}
    definitions = {}

    for route in app.routes:
        # Skip system routes (e.g., /openapi.json, /docs)
        if not hasattr(route, "methods") or not route.methods:
            continue

        for method in route.methods:
            if method in ("HEAD", "OPTIONS"):
                continue

            op_id = route.unique_id or f"{method.lower()}_{route.path.replace('/', '_')}"

            operation = {
                "method": method.upper(),
                "path": route.path,
                "description": route.summary or route.description or "",
                "arguments": {},
                "return": {
                    "code": 200
                }
            }

            if hasattr(route, "dependant"):
                # Request Body Payload
                if route.dependant.body_params:
                    for body_param in route.dependant.body_params:
                        param_type = body_param.annotation
                        if inspect.isclass(param_type) and issubclass(param_type, BaseModel):
                            type_name = param_type.__name__
                            if type_name not in definitions:
                                definitions[type_name] = extract_pydantic_definition(param_type)

                            operation["arguments"]["payload"] = {
                                "in": "body",
                                "schema": {
                                    "type": "reference",
                                    "target": type_name
                                }
                            }

                # Path & Query Parameters
                for param in route.dependant.path_params + route.dependant.query_params:
                    location = "path" if param in route.dependant.path_params else "query"
                    operation["arguments"][param.name] = {
                        "in": location,
                        "schema": pydantic_type_to_typeschema(param.annotation)
                    }

            # Response Model
            if route.response_model:
                resp_model = route.response_model
                if inspect.isclass(resp_model) and issubclass(resp_model, BaseModel):
                    type_name = resp_model.__name__
                    if type_name not in definitions:
                        definitions[type_name] = extract_pydantic_definition(resp_model)

                    operation["return"]["schema"] = {
                        "type": "reference",
                        "target": type_name
                    }

            operations[op_id] = operation

    return {
        "operations": operations,
        "definitions": definitions
    }

if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.stderr.write("Error: Expected entrypoint argument in format 'module:app'\n")
        sys.exit(1)

    entrypoint = sys.argv[1]
    if ":" not in entrypoint:
        sys.stderr.write("Error: Entrypoint must be formatted as 'module:app' (e.g. main:app)\n")
        sys.exit(1)

    module_name, app_name = entrypoint.split(":", 1)

    # Ensure current working directory is on Python path so imports resolve
    sys.path.insert(0, ".")

    try:
        module = importlib.import_module(module_name)
        app = getattr(module, app_name)

        spec = extract_typeapi(app)
        print(json.dumps(spec, indent=2))
    except Exception as e:
        sys.stderr.write(f"Extraction failed: {str(e)}\n")
        sys.exit(1)