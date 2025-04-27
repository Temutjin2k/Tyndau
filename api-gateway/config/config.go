package config

import (
	"os"
)

type Config struct {
	Environment string
	Version     string
	Port        string
	Services    []Service
	CorsURLs    string
}

type Service struct {
	Name       string
	GrpcAddr   string
	ApiVersion string
	Status     string
	Routes     []RouteConfig
}

func NewConfig() Config {
	return Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Version:     getEnv("VERSION", "v1"),
		Port:        getEnv("PORT", "8080"),
		Services: []Service{
			{
				Name:       "inventory service",
				GrpcAddr:   getEnv("INVENTORY_GRPC_ADDR", "localhost:50051"),
				ApiVersion: "v1",
				Status:     "down",
				Routes: []RouteConfig{
					{
						Method:      "GET",
						Path:        "/inventory",
						GRPCService: "InventoryService",
						GRPCMethod:  "ListProducts",
						RequestType: "ListProductsRequest",
						QueryParams: []string{"Page", "PageSize", "SortBy"},
					},
					{
						Method:      "GET",
						Path:        "/inventory/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "GetProductByID",
						RequestType: "GetProductRequest",
						PathParams:  []string{"Id"},
					},
					{
						Method:      "POST",
						Path:        "/inventory",
						GRPCService: "InventoryService",
						GRPCMethod:  "CreateProduct",
						RequestType: "CreateProductRequest",
					},
					{
						Method:      "PUT",
						Path:        "/inventory/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "UpdateProduct",
						RequestType: "UpdateProductRequest",
						PathParams:  []string{"id"},
					},
					{
						Method:      "DELETE",
						Path:        "/inventory/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "DeleteProduct",
						RequestType: "DeleteProductRequest",
						PathParams:  []string{"Id"},
					},
					{
						Method:      "POST",
						Path:        "/categories",
						GRPCService: "InventoryService",
						GRPCMethod:  "CreateCategory",
						RequestType: "CreateCategoryRequest",
					},
					{
						Method:      "GET",
						Path:        "/categories/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "GetCategoryByID",
						RequestType: "GetCategoryRequest",
						PathParams:  []string{"Id"},
					},
					{
						Method:      "PUT",
						Path:        "/categories/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "UpdateCategory",
						RequestType: "UpdateCategoryRequest",
						PathParams:  []string{"id"},
					},
					{
						Method:      "DELETE",
						Path:        "/categories/{id}",
						GRPCService: "InventoryService",
						GRPCMethod:  "DeleteCategory",
						RequestType: "DeleteCategoryRequest",
						PathParams:  []string{"Id"},
					},
					{
						Method:      "GET",
						Path:        "/categories",
						GRPCService: "InventoryService",
						GRPCMethod:  "ListCategories",
						RequestType: "ListCategoriesRequest",
						QueryParams: []string{"Page", "PageSize", "SortBy"},
					},
				},
			},
			// Orders service configuration
		},
		CorsURLs: getEnv("CORS_URLS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
