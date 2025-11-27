package pdf

import (
	"bytes"
	"fmt"
	"time"

	"codeberg.org/go-pdf/fpdf"
	"github.com/et0/links-checker-25-11-2025/internal/model"
)

func Create(lists []*model.LinkList) ([]byte, error) {
	// Создаем новый PDF документ
	pdf := fpdf.New("P", "mm", "A4", "")

	// Добавляем страницу
	pdf.AddPage()

	// Устанавливаем шрифт для заголовка
	pdf.SetFont("Arial", "B", 16)

	// Заголовок
	pdf.CellFormat(0, 10, "Link Status Report", "", 0, "C", false, 0, "")
	pdf.Ln(12)

	// Дата генерации
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")), "", 0, "L", false, 0, "")
	pdf.Ln(15)

	// Обрабатываем каждый batch
	for i, list := range lists {
		if i > 0 {
			pdf.Ln(10) // Добавляем отступ между batches
		}

		// Заголовок batch
		pdf.SetFont("Arial", "B", 14)
		pdf.CellFormat(0, 10, fmt.Sprintf("Batch ID: %d", list.ID), "", 0, "L", false, 0, "")
		pdf.Ln(8)

		// Информация о batch
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(0, 10, fmt.Sprintf("Created: %s", list.CreatedAt.Format("2006-01-02 15:04:05")), "", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(0, 10, fmt.Sprintf("Status: %s", list.Status), "", 0, "L", false, 0, "")
		pdf.Ln(5)
		pdf.CellFormat(0, 10, fmt.Sprintf("Links count: %d", len(*list.Links)), "", 0, "L", false, 0, "")
		pdf.Ln(8)

		// Таблица со ссылками, если они есть
		if len(*list.Links) > 0 {
			// Заголовок таблицы
			pdf.SetFont("Arial", "B", 10)
			pdf.SetFillColor(240, 240, 240)
			pdf.CellFormat(150, 8, "URL", "1", 0, "L", true, 0, "")
			pdf.CellFormat(30, 8, "Status", "1", 1, "C", true, 0, "")

			// Данные таблицы
			pdf.SetFont("Arial", "", 9)
			for _, link := range *list.Links {

				fmt.Println(link.URL)
				status := "Unavailable"
				if link.Status == model.Available {
					status = "Available"
				}

				pdf.CellFormat(150, 8, link.URL, "1", 0, "L", false, 0, "")
				pdf.CellFormat(30, 8, status, "1", 1, "C", false, 0, "")
			}
		} else {
			pdf.SetFont("Arial", "I", 10)
			pdf.CellFormat(0, 10, "No links to display", "", 1, "L", false, 0, "")
		}
	}

	// Генерируем PDF в буфер
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}
