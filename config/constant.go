package config

const configFileName = "config.json"
const configFileDir = "kpf"

type PortForwardRule struct {
	Prefix string `json:"prefix"`
	Port   string `json:"port"`
}

type ConfigStructure struct {
	Namespace        string            `json:"namespace"`
	PortForwardRules []PortForwardRule `json:"port_forward_rules"`
}

var configTemplate = ConfigStructure{
	Namespace: "<namespace_name>",
	PortForwardRules: []PortForwardRule{
		{Prefix: "<pod_name1>", Port: "<pod_port1>"},
		{Prefix: "<pod_name2>", Port: "<pod_port2>"},
	},
}
