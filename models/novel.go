package models

import (
    "gorm.io/gorm"
    "time"
)

type Author struct {
    gorm.Model
    Name    string   `gorm:"type:varchar(255);not null"`
    Bio     string   `gorm:"type:text"`
    Novels  []Novel  `gorm:"many2many:author_novels;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Genre struct {
    gorm.Model
    Name   string   `gorm:"type:varchar(100);unique;not null"`
    Novels []Novel  `gorm:"many2many:novel_genres;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Tags struct {
    gorm.Model
    Name   string   `gorm:"type:varchar(50);unique;not null"`
    Novels []Novel  `gorm:"many2many:novel_tags;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Novel struct {
    gorm.Model
    Title          string      `gorm:"type:varchar(255);not null"`
    Slug           string      `gorm:"type:varchar(255);uniqueIndex"` // for SEO URLs
    Summary        string      `gorm:"type:text"`
    PublishedDate  time.Time
    PageCount      int
    Language       string      `gorm:"type:varchar(50)"`
    Status         string      `gorm:"type:varchar(50);default:'Published'"` // Draft/Published/Upcoming
    CoverImageURL  string      `gorm:"type:text"` // For cover images

    Authors        []Author    `gorm:"many2many:author_novels;"`
    Genres         []Genre     `gorm:"many2many:novel_genres;"`
    Tags           []Tags       `gorm:"many2many:novel_tags;"`
    Reviews        []Review    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
    Comments       []Comment   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}


type Review struct {
    gorm.Model
    Rating  int     `gorm:"not null"`
    Comment string  `gorm:"type:text"`
    NovelID uint    `gorm:"not null;index"`
    UserID  uint    `gorm:"not null;index"`
}

type Comment struct {
    gorm.Model
    Content string `gorm:"type:text;not null"`
    NovelID uint   `gorm:"not null;index"`
    UserID  uint   `gorm:"not null;index"`
}