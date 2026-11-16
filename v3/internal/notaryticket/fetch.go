package notaryticket

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const lookupURL = "https://api.apple-cloudkit.com/database/1/com.apple.gk.ticket-delivery/production/public/records/lookup"

type Client struct{ HTTPClient *http.Client }

// Fetch uses the ticket-delivery endpoint observed in Apple's stapler output.
// No credentials are sent. Apple does not document this service as a public API.
func (c Client) Fetch(ctx context.Context, expected []Digest) (*Ticket, error) {
	if err := validateDigests(expected); err != nil {
		return nil, err
	}
	type recordQuery struct {
		Name string `json:"recordName"`
	}
	query := struct {
		Records []recordQuery `json:"records"`
	}{}
	requested := map[string]bool{}
	for _, hash := range expected {
		name := fmt.Sprintf("2/%d/%x", hash.Algorithm, hash.Value)
		if !requested[name] {
			query.Records = append(query.Records, recordQuery{name})
			requested[name] = true
		}
	}
	data, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, lookupURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: time.Minute}
	}
	copyClient := *client
	copyClient.Jar = nil
	copyClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := copyClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("fetch notarization ticket: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ticket service returned HTTP %d", response.StatusCode)
	}
	data, err = io.ReadAll(io.LimitReader(response.Body, 8<<20+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 8<<20 {
		return nil, errors.New("ticket response exceeds 8 MiB")
	}
	var reply struct {
		Records []struct {
			Name    string `json:"recordName"`
			Type    string `json:"recordType"`
			Deleted bool   `json:"deleted"`
			Error   string `json:"serverErrorCode"`
			Fields  struct {
				SignedTicket struct {
					Type  string `json:"type"`
					Value []byte `json:"value"`
				} `json:"signedTicket"`
			} `json:"fields"`
		} `json:"records"`
	}
	if err := json.Unmarshal(data, &reply); err != nil {
		return nil, err
	}
	if len(reply.Records) != len(requested) {
		return nil, errors.New("ticket response record count differs from request")
	}
	var result *Ticket
	for _, record := range reply.Records {
		if !requested[record.Name] || record.Deleted || record.Error != "" || record.Type != "DeveloperIDTicket" || record.Fields.SignedTicket.Type != "BYTES" {
			return nil, errors.New("ticket response contains a missing, deleted or unexpected record")
		}
		delete(requested, record.Name)
		ticket, err := Verify(record.Fields.SignedTicket.Value, expected)
		if err != nil {
			return nil, err
		}
		result = ticket
	}
	return result, nil
}
