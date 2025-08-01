package config

type RouterConfig struct {
	EnableProfiling bool
}

func loadRouterConfig() *RouterConfig {
	return &RouterConfig{
		EnableProfiling: getEnvBool("ENABLE_PROFILING", false),
	}
}
