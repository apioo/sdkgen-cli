import java.lang.annotation.Annotation;
import java.lang.reflect.Field;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.util.*;

// Reads Spring classes via Reflection without external dependencies
public class Extract {

    public static void main(String[] args) {
        try {
            // Note: In a complete implementation, scan package classes from classpath.
            // Here we show the core inspection logic for annotated controller classes.
            Map<String, Object> spec = new LinkedHashMap<>();
            Map<String, Object> operations = new LinkedHashMap<>();
            Map<String, Object> definitions = new LinkedHashMap<>();

            // Example controller inspection loop (pass target controller classes)
            // List<Class<?>> controllers = findControllerClasses(args[0]);

            spec.put("operations", operations);
            spec.put("definitions", definitions);

            System.out.println(toJson(spec, 0));
        } catch (Exception e) {
            System.err.println("Spring Extraction Error: " + e.getMessage());
            e.printStackTrace(System.err);
            System.exit(1);
        }
    }

    public static void processControllerClass(Class<?> controllerClass, Map<String, Object> operations, Map<String, Object> definitions) {
        String basePath = "";

        // Check for @RequestMapping on class level
        for (Annotation ann : controllerClass.getAnnotations()) {
            if (ann.annotationType().getSimpleName().equals("RequestMapping")) {
                basePath = getAnnotationValue(ann, "value");
            }
        }

        for (Method method : controllerClass.getDeclaredMethods()) {
            String httpMethod = null;
            String path = "";

            for (Annotation ann : method.getAnnotations()) {
                String name = ann.annotationType().getSimpleName();
                if (name.equals("GetMapping")) {
                    httpMethod = "GET";
                    path = getAnnotationValue(ann, "value");
                } else if (name.equals("PostMapping")) {
                    httpMethod = "POST";
                    path = getAnnotationValue(ann, "value");
                } else if (name.equals("PutMapping")) {
                    httpMethod = "PUT";
                    path = getAnnotationValue(ann, "value");
                } else if (name.equals("DeleteMapping")) {
                    httpMethod = "DELETE";
                    path = getAnnotationValue(ann, "value");
                } else if (name.equals("RequestMapping")) {
                    httpMethod = getAnnotationValue(ann, "method");
                    path = getAnnotationValue(ann, "value");
                }
            }

            if (httpMethod == null) continue;

            String fullPath = sanitizePath(basePath + "/" + path);
            String opId = httpMethod.toLowerCase() + "_" + fullPath.replaceAll("[/{}]", "_").replaceAll("_+", "_");

            Map<String, Object> operation = new LinkedHashMap<>();
            operation.put("method", httpMethod);
            operation.put("path", fullPath);
            operation.put("description", "");

            Map<String, Object> arguments = new LinkedHashMap<>();
            operation.put("arguments", arguments);

            Map<String, Object> returnObj = new LinkedHashMap<>();
            returnObj.put("code", 200);
            operation.put("return", returnObj);

            // Process parameters (@RequestBody, @PathVariable, @RequestParam)
            for (Parameter param : method.getParameters()) {
                boolean isBody = false;
                boolean isPath = false;

                for (Annotation ann : param.getAnnotations()) {
                    String annName = ann.annotationType().getSimpleName();
                    if (annName.equals("RequestBody")) isBody = true;
                    if (annName.equals("PathVariable")) isPath = true;
                }

                if (isBody) {
                    Class<?> dtoClass = param.getType();
                    String typeName = dtoClass.getSimpleName();

                    extractTypeSchemaDefinition(dtoClass, definitions);

                    Map<String, Object> payloadArg = new LinkedHashMap<>();
                    payloadArg.put("in", "body");
                    Map<String, Object> schemaRef = new LinkedHashMap<>();
                    schemaRef.put("type", "reference");
                    schemaRef.put("target", typeName);
                    payloadArg.put("schema", schemaRef);

                    arguments.put("payload", payloadArg);
                } else {
                    Map<String, Object> paramArg = new LinkedHashMap<>();
                    paramArg.put("in", isPath ? "path" : "query");
                    paramArg.put("schema", javaTypeToTypeSchema(param.getType()));
                    arguments.put(param.getName(), paramArg);
                }
            }

            // Return type DTO
            Class<?> returnType = method.getReturnType();
            if (!returnType.equals(Void.TYPE) && !returnType.equals(String.class)) {
                String typeName = returnType.getSimpleName();
                extractTypeSchemaDefinition(returnType, definitions);

                Map<String, Object> schemaRef = new LinkedHashMap<>();
                schemaRef.put("type", "reference");
                schemaRef.put("target", typeName);
                returnObj.put("schema", schemaRef);
            }

            operations.put(opId, operation);
        }
    }

    private static void extractTypeSchemaDefinition(Class<?> clazz, Map<String, Object> definitions) {
        String typeName = clazz.getSimpleName();
        if (definitions.containsKey(typeName) || clazz.isPrimitive() || clazz.getName().startsWith("java.")) {
            return;
        }

        Map<String, Object> definition = new LinkedHashMap<>();
        definition.put("type", "object");

        Map<String, Object> properties = new LinkedHashMap<>();
        List<String> required = new ArrayList<>();

        for (Field field : clazz.getDeclaredFields()) {
            String fieldName = field.getName();
            properties.put(fieldName, javaTypeToTypeSchema(field.getType()));

            // Check basic nullability or validation annotations
            boolean isNotNull = false;
            for (Annotation ann : field.getAnnotations()) {
                String name = ann.annotationType().getSimpleName();
                if (name.equals("NotNull") || name.equals("NotBlank")) {
                    isNotNull = true;
                }
            }
            if (isNotNull) required.add(fieldName);
        }

        definition.put("properties", properties);
        if (!required.isEmpty()) {
            definition.put("required", required);
        }

        definitions.put(typeName, definition);
    }

    private static Map<String, Object> javaTypeToTypeSchema(Class<?> clazz) {
        Map<String, Object> schema = new LinkedHashMap<>();
        if (clazz == Integer.class || clazz == int.class || clazz == Long.class || clazz == long.class) {
            schema.put("type", "integer");
        } else if (clazz == Float.class || clazz == float.class || clazz == Double.class || clazz == double.class) {
            schema.put("type", "number");
        } else if (clazz == Boolean.class || clazz == boolean.class) {
            schema.put("type", "boolean");
        } else if (Collection.class.isAssignableFrom(clazz)) {
            schema.put("type", "array");
            Map<String, Object> itemSchema = new LinkedHashMap<>();
            itemSchema.put("type", "string");
            schema.put("schema", itemSchema);
        } else if (!clazz.getName().startsWith("java.")) {
            schema.put("type", "reference");
            schema.put("target", clazz.getSimpleName());
        } else {
            schema.put("type", "string");
        }
        return schema;
    }

    private static String getAnnotationValue(Annotation ann, String attribute) {
        try {
            Method m = ann.annotationType().getMethod(attribute);
            Object val = m.invoke(ann);
            if (val instanceof String[]) {
                String[] arr = (String[]) val;
                return arr.length > 0 ? arr[0] : "";
            }
            return val.toString();
        } catch (Exception e) {
            return "";
        }
    }

    private static String sanitizePath(String path) {
        return path.replaceAll("//+", "/");
    }

    // Helper to print formatted JSON without external library dependencies (like Jackson/Gson)
    private static String toJson(Object obj, int indent) {
        if (obj instanceof Map) {
            Map<?, ?> map = (Map<?, ?>) obj;
            StringBuilder sb = new StringBuilder("{\n");
            int i = 0;
            for (Map.Entry<?, ?> entry : map.entrySet()) {
                sb.append("  ".repeat(indent + 1))
                  .append("\"").append(entry.getKey()).append("\": ")
                  .append(toJson(entry.getValue(), indent + 1));
                if (++i < map.size()) sb.append(",");
                sb.append("\n");
            }
            sb.append("  ".repeat(indent)).append("}");
            return sb.toString();
        } else if (obj instanceof List) {
            List<?> list = (List<?>) obj;
            StringBuilder sb = new StringBuilder("[\n");
            for (int i = 0; i < list.size(); i++) {
                sb.append("  ".repeat(indent + 1)).append(toJson(list.get(i), indent + 1));
                if (i < list.size() - 1) sb.append(",");
                sb.append("\n");
            }
            sb.append("  ".repeat(indent)).append("]");
            return sb.toString();
        } else if (obj instanceof String) {
            return "\"" + obj + "\"";
        } else {
            return String.valueOf(obj);
        }
    }
}
