package handler

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	rows, err := h.DB.Query("SELECT id, name FROM users LIMIT 10")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		users = append(users, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	return c.JSON(http.StatusOK, users)
}
