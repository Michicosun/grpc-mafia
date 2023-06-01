package registry

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenRandomName() string {
	return uuid.New().String()
}

func SaveAvatar(header *multipart.FileHeader, filename string) error {
	file, err := header.Open()
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return Server.ava_storage.Write(filename, bytes)
}

func EndWithError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}
