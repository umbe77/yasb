package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/umbe77/yasb/models"
)

type HttpRequest struct {
	Endpoint   string            `json:"endpoint"`
	Verb       string            `json:"verb"`
	Body       map[string]any    `json:"body"`
	Headers    map[string]string `json:"headers"`
	NextAction models.Action     `json:"nextAction"`
	resp       map[string]any
}

func (h *HttpRequest) Exec(ctx context.Context, m *models.Message) error {
	c := http.Client{
		Timeout: time.Duration(10) * time.Second,
	}

	var b []byte
	if h.Body != nil {
		var err error
		if b, err = json.Marshal(h.Body); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(h.Verb, h.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	for k, v := range h.Headers {
		req.Header.Add(k, v)
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resp_message map[string]any
	if err := json.Unmarshal(body, &resp_message); err != nil {
		return err
	}

	h.resp = resp_message

	return nil
}

func (h *HttpRequest) Next() models.Action {
	return h.NextAction
}

func (h *HttpRequest) Outputs() []models.Action {
	return []models.Action{
		h.NextAction,
	}
}

func (h *HttpRequest) GetType() string {
	return "http"
}
