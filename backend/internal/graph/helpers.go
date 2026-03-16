package graph

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func saveBase64Image(imageDir, payload string) (string, error) {
	data := payload

	if strings.Contains(payload, ",") {
		parts := strings.SplitN(payload, ",", 2)
		data = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s.jpg", uuid.NewString())
	fullPath := filepath.Join(imageDir, filename)

	if err := os.WriteFile(fullPath, decoded, 0644); err != nil {
		return "", err
	}

	return "/media/pet-images/" + filename, nil
}
