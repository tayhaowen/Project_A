package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"project-a/internal/dto"

	"github.com/PuerkitoBio/goquery"
)

const tideTimingsURL = "https://www.nea.gov.sg/corporate-functions/weather/tide-timings"

// TideService encapsulates the logic to scrape tide timings from NEA.
type TideService struct {
	client *http.Client
	url    string
}

// NewTideService constructs a TideService with the provided HTTP client.
func NewTideService(client *http.Client) *TideService {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}

	return &TideService{client: client, url: tideTimingsURL}
}

// GetTideTimings fetches and parses tide timings for the available months.
func (s *TideService) GetTideTimings(ctx context.Context) ([]dto.TideMonth, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return nil, fmt.Errorf("create tide request: %w", err)
	}
	req.Header.Set("User-Agent", "project-a-tide-scraper/1.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch tide timings: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tide timings responded with status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse tide timings page: %w", err)
	}

	labelMap := extractMonthLabels(doc)
	months := make([]dto.TideMonth, 0)

	doc.Find(".forecast-widget__content").Each(func(_ int, section *goquery.Selection) {
		key := strings.TrimSpace(attrOr(section, "data-box"))
		if key == "" {
			key = strings.TrimSpace(attrOr(section, "id"))
		}

		monthLabel := labelMap[key]
		if monthLabel == "" {
			monthLabel = key
		}

		days := parseTideTable(section)
		if len(monthLabel) == 0 || len(days) == 0 {
			return
		}

		months = append(months, dto.TideMonth{Month: monthLabel, Days: days})
	})

	if len(months) == 0 {
		return nil, fmt.Errorf("no tide data found on page")
	}

	return months, nil
}

func extractMonthLabels(doc *goquery.Document) map[string]string {
	labels := make(map[string]string)
	doc.Find(".tab__nav-item").Each(func(_ int, node *goquery.Selection) {
		key := strings.TrimSpace(attrOr(node, "data-box"))
		if key == "" {
			return
		}

		text := strings.TrimSpace(node.Text())
		if text != "" {
			labels[key] = text
		}
	})
	return labels
}

func parseTideTable(section *goquery.Selection) []dto.TideDay {
	dayMap := make(map[int]*dto.TideDay)
	dayOrder := make([]int, 0)
	currentDay := 0

	section.Find("table tbody tr").Each(func(_ int, row *goquery.Selection) {
		cells := extractCells(row)
		if len(cells) == 0 {
			return
		}

		switch len(cells) {
		case 4:
			dayVal, err := strconv.Atoi(cells[0])
			if err != nil {
				return
			}
			currentDay = dayVal
			if _, exists := dayMap[dayVal]; !exists {
				dayMap[dayVal] = &dto.TideDay{Day: dayVal}
				dayOrder = append(dayOrder, dayVal)
			}

			if obs := buildObservation(cells[1:]); obs != nil {
				dayMap[dayVal].Observations = append(dayMap[dayVal].Observations, *obs)
			}
		case 3:
			if currentDay == 0 {
				return
			}
			if obs := buildObservation(cells); obs != nil {
				dayMap[currentDay].Observations = append(dayMap[currentDay].Observations, *obs)
			}
		}
	})

	days := make([]dto.TideDay, 0, len(dayOrder))
	for _, day := range dayOrder {
		if data := dayMap[day]; data != nil && len(data.Observations) > 0 {
			days = append(days, *data)
		}
	}

	return days
}

func extractCells(row *goquery.Selection) []string {
	cells := make([]string, 0)
	row.Find("td").Each(func(_ int, cell *goquery.Selection) {
		value := strings.TrimSpace(cell.Text())
		if value != "" {
			cells = append(cells, value)
		}
	})
	return cells
}

func buildObservation(values []string) *dto.TideObservation {
	if len(values) != 3 {
		return nil
	}

	height, err := strconv.ParseFloat(strings.ReplaceAll(values[1], ",", ""), 64)
	if err != nil {
		return nil
	}

	return &dto.TideObservation{
		Time:   values[0],
		Height: height,
		Level:  values[2],
	}
}

func attrOr(sel *goquery.Selection, attr string) string {
	if val, ok := sel.Attr(attr); ok {
		return val
	}
	return ""
}
