// /api/db/lead.go

package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lead struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	BusinessName   string       `gorm:"size:255"`
	RegisteredName string       `gorm:"size:255"`
	FoundationDate sql.NullTime `gorm:"type:date"`
	Address        string       `gorm:"type:text"`
	City           string       `gorm:"type:text"`
	State          string       `gorm:"type:text"`
	Country        string       `gorm:"type:text"`
	ZIPCode        string       `gorm:"type:text"`
	Owner          string       `gorm:"type:text"`
	Source         string       `gorm:"type:text"`
	Phone          string       `gorm:"size:50"`
	Whatsapp       string       `gorm:"size:50"`
	Website        string       `gorm:"type:text"`
	Email          string       `gorm:"type:text"`

	Instagram             string  `gorm:"type:text"`
	Facebook              string  `gorm:"type:text"`
	TikTok                string  `gorm:"type:text"`
	CompanyRegistrationID string  `gorm:"type:text"`
	Categories            string  `gorm:"type:text"`
	Rating                float64 `gorm:"type:numeric"`
	PriceLevel            int     `gorm:"default:0"`
	UserRatingsTotal      int     `gorm:"default:0"`
	Vicinity              string  `gorm:"type:text"`
	PermanentlyClosed     bool    `gorm:"default:false"`

	CompanySize    string  `gorm:"size:50"`
	Revenue        float64 `gorm:"type:numeric"`
	EmployeesCount int     `gorm:"default:0"`
	Description    string  `gorm:"type:text"`

	PrimaryActivity     string  `gorm:"type:text"`
	SecondaryActivities string  `gorm:"type:text"`
	Types               string  `gorm:"type:text"`
	EquityCapital       float64 `gorm:"type:numeric"`

	BusinessStatus string `gorm:"type:text"`

	Quality      string `gorm:"size:50"`
	SearchTerm   string `gorm:"size:50"`
	FieldsFilled int    `gorm:"default:0"`
	GoogleId     string `gorm:"type:text"`

	Category string `gorm:"type:text"`
	Radius   int    `gorm:"default:0"`

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func CreateLead(lead *Lead) error {
	var existingLead Lead

	result := DB.Where("google_id = ?", lead.GoogleId).First(&existingLead)
	if result.Error == nil {
		log.Printf("Lead com GoogleId %s já existe. Ignorando inserção.", lead.GoogleId)
		return nil
	}

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return fmt.Errorf("erro ao verificar se o lead já existe: %v", result.Error)
	}

	result = DB.Create(lead)
	if result.Error != nil {
		return fmt.Errorf("falha ao salvar lead no banco de dados: %v", result.Error)
	}
	log.Printf("Lead salvo com sucesso: Nome=%s, WhatsApp=%s", lead.BusinessName, lead.Whatsapp)
	return nil
}

func GetLeads() ([]Lead, error) {
	var leads []Lead
	result := DB.Find(&leads)
	if result.Error != nil {
		return nil, result.Error
	}
	return leads, nil
}

func GetLeadByGoogleId(googleId string) (*Lead, error) {
	var lead Lead
	result := DB.Where("google_id = ?", googleId).First(&lead)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("Lead não encontrado")
		}
		return nil, result.Error
	}
	return &lead, nil
}

func GetLeadIdByGoogleId(googleId string) (uuid.UUID, error) {
	var lead Lead
	result := DB.Select("id").Where("google_id = ?", googleId).First(&lead)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return uuid.Nil, fmt.Errorf("Lead não encontrado para o Google ID: %s", googleId)
		}
		return uuid.Nil, result.Error
	}
	return lead.ID, nil
}

func GetLeadByID(leadID uuid.UUID) (*Lead, error) {
	var lead Lead
	result := DB.First(&lead, "id = ?", leadID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &lead, nil
}

func UpdateLead(lead *Lead) error {
	existingLead, err := GetLeadByID(lead.ID)
	if err != nil {
		return fmt.Errorf("erro ao buscar o lead: %v", err)
	}

	if existingLead == nil {
		return fmt.Errorf("Lead não encontrado para ID: %s", lead.ID)
	}

	existingLead.RegisteredName = lead.RegisteredName
	existingLead.FoundationDate = lead.FoundationDate
	existingLead.Address = lead.Address
	existingLead.City = lead.City
	existingLead.State = lead.State

	existingLead.ZIPCode = lead.ZIPCode
	existingLead.Owner = lead.Owner
	existingLead.Source = lead.Source
	existingLead.Phone = lead.Phone
	existingLead.Whatsapp = lead.Whatsapp
	existingLead.Website = lead.Website
	existingLead.Email = lead.Email
	existingLead.Instagram = lead.Instagram
	existingLead.Facebook = lead.Facebook
	existingLead.TikTok = lead.TikTok
	existingLead.CompanyRegistrationID = lead.CompanyRegistrationID
	existingLead.Categories = lead.Categories
	existingLead.Rating = lead.Rating
	existingLead.PriceLevel = lead.PriceLevel
	existingLead.UserRatingsTotal = lead.UserRatingsTotal
	existingLead.Vicinity = lead.Vicinity
	existingLead.PermanentlyClosed = lead.PermanentlyClosed
	existingLead.CompanySize = lead.CompanySize
	existingLead.Revenue = lead.Revenue
	existingLead.EmployeesCount = lead.EmployeesCount

	if lead.Description != "" {
		if existingLead.Description != "" {
			existingLead.Description = fmt.Sprintf("%s\n%s", existingLead.Description, lead.Description)
		} else {
			existingLead.Description = lead.Description
		}
	}
	existingLead.PrimaryActivity = lead.PrimaryActivity
	existingLead.SecondaryActivities = lead.SecondaryActivities
	existingLead.Types = lead.Types
	existingLead.EquityCapital = lead.EquityCapital
	existingLead.BusinessStatus = lead.BusinessStatus
	existingLead.Quality = lead.Quality
	existingLead.SearchTerm = lead.SearchTerm
	existingLead.FieldsFilled = lead.FieldsFilled

	existingLead.Category = lead.Category
	existingLead.Radius = lead.Radius

	result := DB.Save(existingLead)
	if result.Error != nil {
		return fmt.Errorf("erro ao atualizar o lead: %v", result.Error)
	}
	return nil
}
