// Package entity handles concurrent object Entity
package entity 

import (
    "time"
)

// Entity is one of a concurrent base object
type Entity struct {
    ID string `json:"id" gorm:"type:char(42)"`
    Role string `json:"role" gorm:"type:text;default:default"`
    Meta string `json:"meta" gorm:"type:json"`
    CDate time.Time `json:"cdate" gorm:"type:timestamp with time zone;not null;default:clock_timestamp()"`
}

type postRequest struct {
    CCAddr string `json:"ccaddr"`
    Meta string `json:"meta"`
}

// SafeEntity is safe verison of entity
type SafeEntity struct {
    ID string `json:"ccaddr"`
    Role string `json:"role"`
    CDate time.Time `json:"cdate"`
}

