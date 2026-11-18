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

func init() {
	models.RegisterAction("http", func() models.Action { return &HttpRequest{} })
}

type HttpRequest struct {
	Endpoint   string                `json:"endpoint" bson:"endpoint"`
	Verb       string                `json:"verb" bson:"verb"`
	Body       map[string]any        `json:"body" bson:"body"`
	Headers    map[string]string     `json:"headers" bson:"headers"`
	NextAction models.ActionEnvelope `json:"nextAction" bson:"nextAction"`
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
	return h.NextAction.Action
}

func (h *HttpRequest) Outputs() []models.Action {
	if h.NextAction.Action == nil {
		return []models.Action{}
	}
	return []models.Action{h.NextAction.Action}
}

func (h *HttpRequest) GetType() string {
	return "http"
}
