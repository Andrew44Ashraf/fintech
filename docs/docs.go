package docs

// Swagger documentation placeholder
// This can be empty for now
type swaggerInfo struct {
    Version     string
    Host        string
    BasePath    string
    Schemes     []string
    Title       string
    Description string
}

var SwaggerInfoInstance = swaggerInfo{
    Version:     "1.0",
    Host:        "localhost:8080",
    BasePath:    "/api",
    Schemes:     []string{"http"},
    Title:       "Fintech Service API",
    Description: "API for fintech services",
}