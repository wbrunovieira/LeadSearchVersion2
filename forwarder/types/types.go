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

type OlhamaResponse struct {
	Think          string `json:"<think>"`
	RegisteredName string `json:"RegisteredName"`
	CNPJ           string `json:"CNPJ"`
	Contatos       struct {
		Nome     string `json:"Nome"`
		Telefone string `json:"Telefone"`
		Email    string `json:"Email"`
	} `json:"Contatos"`
	DataFundacao string `json:"DataFundacao"`
	Website      string `json:"Website"`
	RedesSociais struct {
		Facebook  string `json:"Facebook"`
		Instagram string `json:"Instagram"`
		TikTok    string `json:"TikTok"`
		WhatsApp  string `json:"WhatsApp"`
	} `json:"RedesSociais"`
	AnaliseEmpresa string      `json:"AnaliseEmpresa"`
	Message        OlhamaReply `json:"Message"`
}

type OlhamaReply struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
