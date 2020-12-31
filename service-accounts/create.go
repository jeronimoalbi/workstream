// Copyright 2020 Jerónimo José Albi. All rights reserved.
//
// Distributed under the MIT license.
//
// For the full copyright and license information, please view the LICENSE
// file that was distributed with this source code.

package main

import (
	"fmt"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
)

func init() {
	service.Action("create", create)
}

func create(a *kusanagi.Action) (*kusanagi.Action, error) {
	ctx := a.GetContext()
	tx := database.WithContext(ctx)

	// Check that the an account with the same name doesn't exist
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

	// Create the account
	account := Account{Name: name}
	if a.HasParam("active") {
		account.Active = a.GetParam("active").GetValue().(bool)
	}
	result := tx.Create(&account)
	if result.Error != nil {
		a.Log(fmt.Sprintf("Failed to create the account: %v", result.Error), log.ERROR)
		a.Error("Account creation failed", 500, "500 Internal Server Error")
		return a, nil
	}

	// Use the account as body
	if _, err := a.SetEntity(account); err != nil {
		return nil, err
	}
	return a, nil
}
