package controllers

import (
	"github.com/tjcain/theFieldBiologist/models"
	"github.com/tjcain/theFieldBiologist/views"
)

// Admin represents a new admin view
type Admin struct {
	AdminDash *views.View
	us        models.UserService
}

// NewAdmin is used to create a new Admin controller.
// Any error in rendering templates will cause this function to panic.
func NewAdmin(us models.UserService) *Admin {
	return &Admin{
		AdminDash: views.NewView("pages", "admin/dashboard"),
		us:        us,
	}
}
