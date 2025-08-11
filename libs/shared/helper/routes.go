package helper

type Routes struct {
	Name string
	Http struct {
		Method string
		Path   string
		Query  map[string]string
		Body   map[string]string
	}
	Grpc struct {
		Service   string
		Target    string
		Timeoutms int
	}
	Mapping struct {
		PathParams  map[string]string `koanf:"path_params"`
		QueryParams map[string]string `koanf:"query_params"`
		Body        map[string]string `koanf:"body"`
	}
	Auth struct {
		Groups []string
	}
	Cache struct {
		TTLms int
	}
}
