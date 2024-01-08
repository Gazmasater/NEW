package config

func Must(cfg *AgentConfig, err error) *AgentConfig {
	if err != nil {
		panic(err)
	}
	return cfg
}
