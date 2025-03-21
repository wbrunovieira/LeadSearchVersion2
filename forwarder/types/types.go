// /forwarder/types/types.go
package types

import (
	"time"

	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/common"
	"github.com/wbrunovieira/LeadSearchVersion2/data-collector/tavily"
)

type CombinedLeadData struct {
	Lead        common.Lead            `json:"lead"`
	TavilyData  *tavily.TavilyResponse `json:"tavily_data,omitempty"`
	TavilyExtra struct {
		CNPJ    string `json:"cnpj,omitempty"`
		Phone   string `json:"phone,omitempty"`
		Owner   string `json:"owner,omitempty"`
		Email   string `json:"email,omitempty"`
		Website string `json:"website,omitempty"`
	} `json:"tavily_extra,omitempty"`
	SerperData map[string]interface{} `json:"serper_data,omitempty"`
	CNPJData   map[string]interface{} `json:"cnpj_data,omitempty"`
	Prompt     string                 `json:"prompt,omitempty"`

	CompanyDetailsCnpjBiz map[string]string      `json:"company_details,omitempty"`
	InverterData          map[string]interface{} `json:"inverter_data,omitempty"`
}

type OlhamaPayload struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
}

type OlhamaOuterResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	DoneReason string `json:"done_reason,omitempty"`
}

type OlhamaResponse struct {
	Think          string `json:"<think>"`
	RegisteredName string `json:"RegisteredName"`
	CNPJ           string `json:"CNPJ"`
	Contatos       string `json:"Contatos"`
	DataDeFundacao string `json:"DataDeFundacao"`
	Website        string `json:"Website"`
	RedesSociais   struct {
		Facebook  string `json:"Facebook"`
		Instagram string `json:"Instagram"`
		TikTok    string `json:"TikTok"`
		WhatsApp  string `json:"WhatsApp"`
	} `json:"RedesSociais"`
	AnaliseEmpresa interface{} `json:"AnaliseEmpresa"`
	Message        OlhamaReply `json:"Message"`
}

type MinimalOlhamaResponse struct {
	RegisteredName string `json:"RegisteredName"`
	CNPJ           string `json:"CNPJ"`

	DataDeFundacao string `json:"DataDeFundacao"`
	Website        string `json:"Website"`
	RedesSociais   struct {
		Facebook  string `json:"Facebook"`
		Instagram string `json:"Instagram"`
		TikTok    string `json:"TikTok"`
		WhatsApp  string `json:"WhatsApp"`
	} `json:"RedesSociais"`
}

type OlhamaReply struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
