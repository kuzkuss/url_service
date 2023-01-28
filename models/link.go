package models

type Link struct {
	OriginalLink string `json:"original_link,omitempty" validate:"required" gorm:"column:original_link"`
	ShortLink    string `json:"short_link,omitempty" readonly:"true" gorm:"column:short_link"`
}
