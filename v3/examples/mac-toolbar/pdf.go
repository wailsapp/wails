package main

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

const (
	pdfPageWidth  = 612
	pdfPageHeight = 792
)

type daymarkPDFPage struct {
	continuation bool
	bodyLines    []string
}

type daymarkPDFObject struct {
	id   int
	body []byte
}

// renderDaymarkPDF produces a complete PDF document using PDF's built-in
// fonts. It has no AppKit dependency, so the lazy share callback can render on
// its background goroutine without blocking the native application thread.
func renderDaymarkPDF(note sharePayload) ([]byte, error) {
	title := strings.TrimSpace(note.Title)
	if title == "" {
		title = "Untitled note"
	}
	subtitle := strings.TrimSpace(note.Subtitle)
	bodyLines := wrapDaymarkPDFText(strings.TrimSpace(note.Body), 76)
	pages := paginateDaymarkPDF(bodyLines)

	const (
		catalogID    = 1
		pagesID      = 2
		helveticaID  = 3
		boldID       = 4
		obliqueID    = 5
		timesRomanID = 6
		infoID       = 7
		firstPageID  = 8
	)

	objects := []daymarkPDFObject{
		{id: catalogID, body: []byte("<< /Type /Catalog /Pages 2 0 R >>")},
		{id: helveticaID, body: []byte("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>")},
		{id: boldID, body: []byte("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>")},
		{id: obliqueID, body: []byte("<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Oblique >>")},
		{id: timesRomanID, body: []byte("<< /Type /Font /Subtype /Type1 /BaseFont /Times-Roman >>")},
		{id: infoID, body: []byte(fmt.Sprintf("<< /Title (%s) /Author (Daymark) /Creator (Wails) >>", pdfLiteral(title)))},
	}

	pageIDs := make([]int, 0, len(pages))
	for index, page := range pages {
		pageID := firstPageID + index*2
		contentID := pageID + 1
		pageIDs = append(pageIDs, pageID)
		content := buildDaymarkPDFPage(title, subtitle, page, index+1, len(pages))
		objects = append(objects,
			daymarkPDFObject{id: pageID, body: []byte(fmt.Sprintf(
				"<< /Type /Page /Parent %d 0 R /MediaBox [0 0 %d %d] "+
					"/Resources << /Font << /F1 %d 0 R /F2 %d 0 R /F3 %d 0 R /F4 %d 0 R >> >> "+
					"/Contents %d 0 R >>",
				pagesID, pdfPageWidth, pdfPageHeight, helveticaID, boldID, obliqueID, timesRomanID, contentID))},
			daymarkPDFObject{id: contentID, body: []byte(fmt.Sprintf("<< /Length %d >>\nstream\n%s\nendstream", len(content), content))},
		)
	}

	kids := make([]string, len(pageIDs))
	for index, id := range pageIDs {
		kids[index] = fmt.Sprintf("%d 0 R", id)
	}
	objects = append(objects, daymarkPDFObject{id: pagesID, body: []byte(fmt.Sprintf(
		"<< /Type /Pages /Kids [%s] /Count %d >>", strings.Join(kids, " "), len(pageIDs)))})

	return assembleDaymarkPDF(objects, catalogID, infoID)
}

func paginateDaymarkPDF(lines []string) []daymarkPDFPage {
	if len(lines) == 0 {
		return []daymarkPDFPage{{}}
	}
	const firstPageLines = 21
	const continuationLines = 27
	result := make([]daymarkPDFPage, 0, 1+len(lines)/continuationLines)
	for len(lines) > 0 {
		limit := firstPageLines
		continuation := len(result) > 0
		if continuation {
			limit = continuationLines
		}
		if limit > len(lines) {
			limit = len(lines)
		}
		result = append(result, daymarkPDFPage{
			continuation: continuation,
			bodyLines:    append([]string(nil), lines[:limit]...),
		})
		lines = lines[limit:]
	}
	return result
}

func wrapDaymarkPDFText(value string, maxRunes int) []string {
	if value == "" {
		return nil
	}
	paragraphs := strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
	var result []string
	for paragraphIndex, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		if len(words) == 0 {
			if len(result) > 0 && result[len(result)-1] != "" {
				result = append(result, "")
			}
			continue
		}
		line := words[0]
		for _, word := range words[1:] {
			if len([]rune(line))+1+len([]rune(word)) <= maxRunes {
				line += " " + word
				continue
			}
			result = append(result, line)
			line = word
		}
		result = append(result, line)
		if paragraphIndex < len(paragraphs)-1 {
			result = append(result, "")
		}
	}
	for len(result) > 0 && result[len(result)-1] == "" {
		result = result[:len(result)-1]
	}
	return result
}

func buildDaymarkPDFPage(title, subtitle string, page daymarkPDFPage, pageNumber, pageCount int) string {
	var content strings.Builder
	// Warm paper background and Daymark's orange accent.
	content.WriteString("0.972 0.965 0.945 rg 0 0 612 792 re f\n")
	content.WriteString("0.87 0.29 0.15 rg 54 716 7 26 re f\n")
	pdfText(&content, "F2", 9, 70, 726, "0.20 0.19 0.17", "DAYMARK  /  FIELD NOTES")

	displayTitle := title
	if page.continuation {
		displayTitle += " (continued)"
	}
	titleLines := wrapDaymarkPDFText(displayTitle, 36)
	if len(titleLines) > 2 {
		titleLines = titleLines[:2]
	}
	titleY := 674.0
	for _, line := range titleLines {
		pdfText(&content, "F2", 28, 54, titleY, "0.10 0.10 0.09", line)
		titleY -= 34
	}

	bodyTop := titleY - 14
	if !page.continuation && subtitle != "" {
		subtitleLines := wrapDaymarkPDFText(subtitle, 70)
		if len(subtitleLines) > 2 {
			subtitleLines = subtitleLines[:2]
		}
		for _, line := range subtitleLines {
			pdfText(&content, "F3", 12, 56, bodyTop, "0.39 0.37 0.33", line)
			bodyTop -= 17
		}
		bodyTop -= 8
	}

	content.WriteString(fmt.Sprintf("0.78 0.74 0.67 RG 0.7 w 54 %.1f m 558 %.1f l S\n", bodyTop, bodyTop))
	bodyTop -= 34
	cardBottom := 76.0
	cardTop := bodyTop + 20
	content.WriteString(fmt.Sprintf("0.995 0.992 0.982 rg 46 %.1f 520 %.1f re f\n",
		cardBottom, cardTop-cardBottom))
	content.WriteString(fmt.Sprintf("0.86 0.83 0.77 RG 0.6 w 46 %.1f 520 %.1f re S\n",
		cardBottom, cardTop-cardBottom))

	lineY := bodyTop
	for _, line := range page.bodyLines {
		if line == "" {
			lineY -= 10
			continue
		}
		pdfText(&content, "F4", 13, 66, lineY, "0.16 0.15 0.14", line)
		lineY -= 20
	}

	footer := "DAYMARK"
	if pageCount > 1 {
		footer = fmt.Sprintf("DAYMARK  /  %d OF %d", pageNumber, pageCount)
	}
	pdfText(&content, "F1", 8, 54, 40, "0.48 0.45 0.40", footer)
	pdfText(&content, "F1", 8, 418, 40, "0.48 0.45 0.40", "SHARED FROM A WAILS APP")
	return content.String()
}

func pdfText(content *strings.Builder, font string, size, x, y float64, color, value string) {
	content.WriteString(fmt.Sprintf("BT /%s %.1f Tf %s rg %.1f %.1f Td (%s) Tj ET\n",
		font, size, color, x, y, pdfLiteral(value)))
}

func pdfLiteral(value string) string {
	var result strings.Builder
	for _, current := range value {
		switch current {
		case '\\', '(', ')':
			result.WriteByte('\\')
			result.WriteRune(current)
		case '\n', '\r', '\t':
			result.WriteByte(' ')
		case '‘', '’':
			result.WriteByte('\'')
		case '“', '”':
			result.WriteByte('"')
		case '–', '—':
			result.WriteByte('-')
		case '…':
			result.WriteString("...")
		default:
			if current >= 32 && current <= 126 {
				result.WriteRune(current)
			} else if unicode.IsSpace(current) {
				result.WriteByte(' ')
			} else {
				result.WriteByte('?')
			}
		}
	}
	return result.String()
}

func assembleDaymarkPDF(objects []daymarkPDFObject, catalogID, infoID int) ([]byte, error) {
	sort.Slice(objects, func(left, right int) bool { return objects[left].id < objects[right].id })
	if len(objects) == 0 || objects[len(objects)-1].id != len(objects) {
		return nil, fmt.Errorf("PDF object identifiers are not contiguous")
	}

	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n%Daymark\n")
	offsets := make([]int, len(objects)+1)
	for _, object := range objects {
		offsets[object.id] = document.Len()
		fmt.Fprintf(&document, "%d 0 obj\n", object.id)
		document.Write(object.body)
		document.WriteString("\nendobj\n")
	}

	xrefOffset := document.Len()
	fmt.Fprintf(&document, "xref\n0 %d\n", len(offsets))
	document.WriteString("0000000000 65535 f \n")
	for id := 1; id < len(offsets); id++ {
		document.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[id]))
	}
	fmt.Fprintf(&document,
		"trailer\n<< /Size %d /Root %d 0 R /Info %d 0 R >>\nstartxref\n%s\n%%%%EOF\n",
		len(offsets), catalogID, infoID, strconv.Itoa(xrefOffset))
	return document.Bytes(), nil
}
