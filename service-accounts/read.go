// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"errors"
	"fmt"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
	"gorm.io/gorm"
)

func init() {
	service.Action("read", read)
}

func read(a *kusanagi.Action) (*kusanagi.Action, error) {
	var account Account

	// Get the account
	tx := database.WithContext(a.GetContext())
	id := a.GetParam("id").GetValue().(string)
	if r := tx.Where("id = ?", id).First(&account); r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			a.Error("Account not found", 404, "404 Not Found")
		} else {
			a.Log(fmt.Sprintf("Failed to read account: %v", r.Error), log.ERROR)
			a.Error("Operation failed", 500, "500 Internal Server Error")
		}
		return a, nil
	}

	// Use the account as body
	if _, err := a.SetEntity(account); err != nil {
		return nil, err
	}
	return a, nil
}
