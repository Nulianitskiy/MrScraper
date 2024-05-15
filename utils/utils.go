package utils

import (
	"code.sajari.com/docconv/v2"
	"fmt"
	"io"
	"net/http"
	"os"
)

func ExtractTextFromPDF(url string) (string, error) {
	// Создаем временный файл для сохранения скачанного PDF
	tmpFile, err := os.CreateTemp("", "temp-*.pdf")
	if err != nil {
		return "", fmt.Errorf("не удалось создать временный файл: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Удаляем временный файл после завершения работы
	defer tmpFile.Close()

	// Загружаем PDF по ссылке
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("не удалось скачать файл: %w", err)
	}
	defer resp.Body.Close()

	// Сохраняем содержимое в временный файл
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("не удалось сохранить файл: %w", err)
	}

	// Извлекаем текст из PDF с использованием docconv
	res, err := docconv.ConvertPath(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("не удалось извлечь текст из PDF: %w", err)
	}

	return res.Body, nil
}
