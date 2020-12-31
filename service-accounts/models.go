// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Account defines a group for one or more users.
type Account struct {
	ID        string    `json:"id" gorm:"primaryKey; size:36"`
	Name      string    `json:"name" gorm:"not null"`
	Active    bool      `json:"active" gorm:"not null; default:true"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	a.ID = fmt.Sprintf("%s", uuid.NewV4())
	return nil
}

func (Account) TableName() string {
	return "account.accounts"
}
