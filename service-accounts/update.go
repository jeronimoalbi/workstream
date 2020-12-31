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
	service.Action("update", update)
}

func update(a *kusanagi.Action) (*kusanagi.Action, error) {
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

	// List of fields to update
	var fields []string

	// Update the account values
	if a.HasParam("name") {
		count := 0
		name := a.GetParam("name").GetValue().(string)
		query := tx.Model(&Account{}).Select("COUNT(*)").Where("name = ?", name)
		if r := query.Scan(&count); r.Error != nil {
			a.Log(fmt.Sprintf("Failed to check the account: %v", r.Error), log.ERROR)
			a.Error("Operation failed", 500, "500 Internal Server Error")
			return a, nil
		} else if count != 0 {
			a.Error("An account with the same name exists", 400, "400 Bad Request")
			return a, nil
		}

		// Update the name
		account.Name = name
		fields = append(fields, "name")
	}

	if a.HasParam("active") {
		// Update the active value
		account.Active = a.GetParam("active").GetValue().(bool)
		fields = append(fields, "active")
	}

	// Update the account
	if len(fields) > 0 {
		if r := tx.Model(&account).Select(fields).Updates(&account); r.Error != nil {
			a.Log(fmt.Sprintf("Failed to update the account: %v", r.Error), log.ERROR)
			a.Error("Operation failed", 500, "500 Internal Server Error")
			return a, nil
		}
	}

	// Use the account as body
	if _, err := a.SetEntity(account); err != nil {
		return nil, err
	}
	return a, nil
}
