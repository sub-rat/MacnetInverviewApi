package utils

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Pagination(c *gin.Context) (int, int, error) {
	pageStr := c.Query("page")
	page := 0
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			return 0, 0, err
		}
		page = p - 1
		if page <= 0 {
			page = 0
		}
	}

	sizeString := c.Query("size")
	limit := 20
	if sizeString != "" {
		l, err := strconv.Atoi(sizeString)
		if err != nil {
			return 0, 0, err
		}
		limit = l
		if limit <= 0 {
			limit = 20
		}
	}
	return page, limit, nil
}

func HashPassword(password string) (string, error) {
	byte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(byte), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Converts a struct to a map while maintaining the json alias as keys
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}
