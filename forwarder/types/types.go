// forwarder/types/types.go
package types

type CombinedLeadData struct {
	Lead        interface{}            `json:"lead"`
	TavilyData  interface{}            `json:"tavily_data,omitempty"`
	TavilyExtra interface{}            `json:"tavily_extra,omitempty"`
	SerperData  map[string]interface{} `json:"serper_data,omitempty"`
	CNPJData    map[string]interface{} `json:"cnpj_data,omitempty"`
	Prompt      string                 `json:"prompt,omitempty"`
}

type OlhamaPayload struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}
