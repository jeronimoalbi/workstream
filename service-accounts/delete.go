package main

import (
	"errors"
	"fmt"

	kusanagi "github.com/kusanagi/kusanagi-sdk-go/v2"
	"github.com/kusanagi/kusanagi-sdk-go/v2/lib/log"
	"gorm.io/gorm"
)

func init() {
	service.Action("delete", delete_)
}

func delete_(a *kusanagi.Action) (*kusanagi.Action, error) {
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

	// Delete the account
	if r := tx.Delete(&account); r.Error != nil {
		a.Log(fmt.Sprintf("Failed to delete the account: %v", r.Error), log.ERROR)
		a.Error("Operation failed", 500, "500 Internal Server Error")
	}

	return a, nil
}
