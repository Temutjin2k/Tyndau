package config

type RouteConfig struct {
	Method      string   `yaml:"method"`
	Path        string   `yaml:"path"`
	GRPCService string   `yaml:"grpc_service"`
	GRPCMethod  string   `yaml:"grpc_method"`
	RequestType string   `yaml:"request_type"`
	PathParams  []string `yaml:"path_params"`
	QueryParams []string `yaml:"query_params"`
}
