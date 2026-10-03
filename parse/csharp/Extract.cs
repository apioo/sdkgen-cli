using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text.Json;
using System.Text.Json.Nodes;

namespace Sdkgen.Extractor
{
    public class Extract
    {
        public static void Main(string[] args)
        {
            try
            {
                string projectPath = args.Length > 0 ? args[0] : Directory.GetCurrentDirectory();

                var operations = new Dictionary<string, object>();
                var definitions = new Dictionary<string, object>();

                // Locate compiled assemblies in bin directory
                var binDir = Path.Combine(projectPath, "bin");
                if (!Directory.Exists(binDir))
                {
                    Console.Error.WriteLine($"Error: bin directory not found in {projectPath}. Please build the project first (dotnet build).");
                    Environment.Exit(1);
                }

                var dllFiles = Directory.GetFiles(binDir, "*.dll", SearchOption.AllDirectories)
                    .Where(f => !f.Contains("Microsoft.") && !f.Contains("System."))
                    .ToList();

                foreach (var dllPath in dllFiles)
                {
                    try
                    {
                        var assembly = Assembly.LoadFrom(dllPath);
                        ProcessAssembly(assembly, operations, definitions);
                    }
                    catch
                    {
                        // Skip non-executable or third-party assemblies
                    }
                }

                var result = new Dictionary<string, object>
                {
                    ["operations"] = operations,
                    ["definitions"] = definitions
                };

                var options = new JsonSerializerOptions { WriteIndented = true };
                Console.WriteLine(JsonSerializer.Serialize(result, options));
            }
            catch (Exception ex)
            {
                Console.Error.WriteLine($"ASP.NET Extraction Error: {ex.Message}");
                Environment.Exit(1);
            }
        }

        private static void ProcessAssembly(Assembly assembly, Dictionary<string, object> operations, Dictionary<string, object> definitions) {
            var controllerTypes = assembly.GetTypes()
                .Where(t => t.IsClass && !t.IsAbstract &&
                            (t.Name.EndsWith("Controller") || t.GetCustomAttributes().Any(a => a.GetType().Name == "ApiControllerAttribute")))
                .ToList();

            foreach (var controller in controllerTypes)
            {
                string basePath = "";
                var routeAttr = controller.GetCustomAttributes().FirstOrDefault(a => a.GetType().Name == "RouteAttribute");
                if (routeAttr != null)
                {
                    var templateProp = routeAttr.GetType().GetProperty("Template");
                    basePath = templateProp?.GetValue(routeAttr)?.ToString() ?? "";
                    basePath = basePath.Replace("[controller]", controller.Name.Replace("Controller", ""));
                }

                foreach (var method in controller.GetMethods(BindingFlags.Public | BindingFlags.Instance | BindingFlags.DeclaredOnly))
                {
                    string httpMethod = null;
                    string routePath = "";

                    foreach (var attr in method.GetCustomAttributes())
                    {
                        var attrName = attr.GetType().Name;
                        if (attrName == "HttpGetAttribute") httpMethod = "GET";
                        else if (attrName == "HttpPostAttribute") httpMethod = "POST";
                        else if (attrName == "HttpPutAttribute") httpMethod = "PUT";
                        else if (attrName == "HttpDeleteAttribute") httpMethod = "DELETE";
                        else if (attrName == "HttpPatchAttribute") httpMethod = "PATCH";

                        if (httpMethod != null)
                        {
                            var templateProp = attr.GetType().GetProperty("Template");
                            routePath = templateProp?.GetValue(attr)?.ToString() ?? "";
                            break;
                        }
                    }

                    if (httpMethod == null) continue;

                    string fullPath = SanitizePath("/" + basePath + "/" + routePath);
                    string opId = $"{httpMethod.ToLower()}_{fullPath.Replace("/", "_").Replace("{", "").Replace("}", "").Replace(":", "_")}";

                    var operation = new Dictionary<string, object>
                    {
                        ["method"] = httpMethod,
                        ["path"] = fullPath,
                        ["description"] = "",
                        ["arguments"] = new Dictionary<string, object>(),
                        ["return"] = new Dictionary<string, object> { ["code"] = 200 }
                    };

                    var argsDict = (Dictionary<string, object>)operation["arguments"];

                    // Process Parameters
                    foreach (var param in method.GetParameters())
                    {
                        bool isBody = param.GetCustomAttributes().Any(a => a.GetType().Name == "FromBodyAttribute");
                        bool isPath = fullPath.Contains("{" + param.Name + "}") || param.GetCustomAttributes().Any(a => a.GetType().Name == "FromRouteAttribute");

                        if (isBody || (!param.ParameterType.IsPrimitive && param.ParameterType != typeof(string) && !isPath))
                        {
                            var dtoType = param.ParameterType;
                            string typeName = dtoType.Name;

                            ExtractTypeSchemaDefinition(dtoType, definitions);

                            argsDict["payload"] = new Dictionary<string, object>
                            {
                                ["in"] = "body",
                                ["schema"] = new Dictionary<string, object>
                                {
                                    ["type"] = "reference",
                                    ["target"] = typeName
                                }
                            };
                        }
                        else
                        {
                            argsDict[param.Name] = new Dictionary<string, object>
                            {
                                ["in"] = isPath ? "path" : "query",
                                ["schema"] = CSharpTypeToTypeSchema(param.ParameterType)
                            };
                        }
                    }

                    // Process Return Type
                    var returnType = method.ReturnType;
                    if (returnType != typeof(void) && returnType != typeof(System.Threading.Tasks.Task))
                    {
                        if (returnType.IsGenericType && returnType.GetGenericTypeDefinition().Name.StartsWith("Task"))
                        {
                            returnType = returnType.GetGenericArguments()[0];
                        }

                        if (!returnType.IsPrimitive && returnType != typeof(string))
                        {
                            string typeName = returnType.Name;
                            ExtractTypeSchemaDefinition(returnType, definitions);

                            var returnObj = (Dictionary<string, object>)operation["return"];
                            returnObj["schema"] = new Dictionary<string, object>
                            {
                                ["type"] = "reference",
                                ["target"] = typeName
                            };
                        }
                    }

                    operations[opId] = operation;
                }
            }
        }

        private static void ExtractTypeSchemaDefinition(Type type, Dictionary<string, object> definitions)
        {
            string typeName = type.Name;
            if (definitions.ContainsKey(typeName) || type.IsPrimitive || type == typeof(string) || type.Assembly.FullName.StartsWith("System"))
            {
                return;
            }

            var definition = new Dictionary<string, object>
            {
                ["type"] = "object"
            };

            var properties = new Dictionary<string, object>();
            var required = new List<string>();

            foreach (var prop in type.GetProperties(BindingFlags.Public | BindingFlags.Instance))
            {
                string propName = prop.Name;
                properties[propName] = CSharpTypeToTypeSchema(prop.PropertyType);

                // Nullability check
                var nullabilityInfo = new NullabilityInfoContext().Create(prop);
                if (nullabilityInfo.WriteState == NullabilityState.NotNull)
                {
                    required.Add(propName);
                }
            }

            definition["properties"] = properties;
            if (required.Count > 0)
            {
                definition["required"] = required;
            }

            definitions[typeName] = definition;
        }

        private static Dictionary<string, object> CSharpTypeToTypeSchema(Type type)
        {
            var schema = new Dictionary<string, object>();

            if (type == typeof(int) || type == typeof(long) || type == typeof(short))
            {
                schema["type"] = "integer";
            }
            else if (type == typeof(float) || type == typeof(double) || type == typeof(decimal))
            {
                schema["type"] = "number";
            }
            else if (type == typeof(bool))
            {
                schema["type"] = "boolean";
            }
            else if (typeof(IEnumerable).IsAssignableFrom(type) && type != typeof(string))
            {
                schema["type"] = "array";
                schema["schema"] = new Dictionary<string, object> { ["type"] = "string" };
            }
            else if (!type.IsPrimitive && type != typeof(string) && !type.Assembly.FullName.StartsWith("System"))
            {
                schema["type"] = "reference";
                schema["target"] = type.Name;
            }
            else
            {
                schema["type"] = "string";
            }

            return schema;
        }

        private static string SanitizePath(string path)
        {
            return System.Text.RegularExpressions.Regex.Replace(path, @"//+", "/");
        }
    }
}
