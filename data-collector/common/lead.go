// /data-collector/common/lead.go
package common

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Lead struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`

	BusinessName   string       `gorm:"size:255" json:"BusinessName"`
	RegisteredName string       `gorm:"size:255" json:"registered_name"`
	FoundationDate sql.NullTime `gorm:"type:date" json:"foundation_date"`
	Address        string       `gorm:"type:text" json:"address"`
	City           string       `gorm:"type:text" json:"city"`
	State          string       `gorm:"type:text" json:"state"`
	Country        string       `gorm:"type:text" json:"country"`
	ZIPCode        string       `gorm:"type:text" json:"zip_code"`
	Owner          string       `gorm:"type:text" json:"owner"`
	Source         string       `gorm:"type:text" json:"source"`
	Phone          string       `gorm:"size:50" json:"phone"`
	Whatsapp       string       `gorm:"size:50" json:"whatsapp"`
	Website        string       `gorm:"type:text" json:"website"`
	Email          string       `gorm:"type:text" json:"email"`

	Instagram             string  `gorm:"type:text" json:"instagram"`
	Facebook              string  `gorm:"type:text" json:"facebook"`
	TikTok                string  `gorm:"type:text" json:"tiktok"`
	CompanyRegistrationID string  `gorm:"type:text" json:"company_registration_id"`
	Categories            string  `gorm:"type:text" json:"categories"`
	Rating                float64 `gorm:"type:numeric" json:"rating"`
	PriceLevel            int     `gorm:"default:0" json:"price_level"`
	UserRatingsTotal      int     `gorm:"default:0" json:"user_ratings_total"`
	Vicinity              string  `gorm:"type:text" json:"vicinity"`
	PermanentlyClosed     bool    `gorm:"default:false" json:"permanently_closed"`

	CompanySize    string  `gorm:"size:50" json:"company_size"`
	Revenue        float64 `gorm:"type:numeric" json:"revenue"`
	EmployeesCount int     `gorm:"default:0" json:"employees_count"`
	Description    string  `gorm:"type:text" json:"description"`

	PrimaryActivity     string  `gorm:"type:text" json:"primary_activity"`
	SecondaryActivities string  `gorm:"type:text" json:"secondary_activities"`
	Types               string  `gorm:"type:text" json:"types"`
	EquityCapital       float64 `gorm:"type:numeric" json:"equity_capital"`

	BusinessStatus string `gorm:"type:text" json:"business_status"`

	Quality      string `gorm:"size:50" json:"quality"`
	SearchTerm   string `gorm:"size:50" json:"search_term"`
	FieldsFilled int    `gorm:"default:0" json:"fields_filled"`
	GoogleId     string `gorm:"type:text" json:"google_id"`

	Category string `gorm:"type:text" json:"category"`
	Radius   int    `gorm:"default:0" json:"radius"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
