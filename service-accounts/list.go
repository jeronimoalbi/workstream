package main

import (
	"fmt"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
)

func init() {
	service.Action("list", list)
}

func list(a *kusanagi.Action) (*kusanagi.Action, error) {
	var accounts []Account

	// Get all the accounts
	tx := database.WithContext(a.GetContext())
	if r := tx.Find(&accounts); r.Error != nil {
		a.Log(fmt.Sprintf("Failed to list accounts: %v", r.Error), log.ERROR)
		a.Error("Operation failed", 500, "500 Internal Server Error")
		return a, nil
	}

	// Use the accounts as body
	if _, err := a.SetCollection(accounts); err != nil {
		return nil, err
	}
	return a, nil
}
