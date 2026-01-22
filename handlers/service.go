package handlers

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// TODO: Implement Sample and example pipeline

type Message map[string]any

func trigger(ctx context.Context, msg Message) (<-chan Message, error) {

	out := make(chan Message)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				log.Print("DONE")
				return
			case out <- msg:
			}
		}
	}()

	return out, nil
}

func capitalize(ctx context.Context, cancel context.CancelFunc, msg <-chan Message) (<-chan Message, error) {
	out := make(chan Message)

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-msg:
				if ok {
					val, isok := m["name"]
					if !isok {
						cancel()
						return
					}

					caser := cases.Title(language.English)

					m["name"] = caser.String(val.(string))

					out <- m
					log.Printf("%v", m)
				} else {
					log.Printf("channel not ready")
					cancel()
					return
				}
			}
		}
	}()
	return out, nil
}

func sink(ctx context.Context, cancel context.CancelFunc, in <-chan Message) Message {
	var msg Message
	for {
		select {
		case <-ctx.Done():
			return nil
		case m, ok := <-in:
			if !ok {
				log.Print("SINK cannot read channel")
				cancel()
				return msg
			}
			msg = m
			return msg
		}
	}
}

func ExecuteService(c *fiber.Ctx) error {
	ctx, cancel := context.WithCancel(c.Context())
	defer cancel()

	msg1, err := trigger(ctx, Message{
		"name": "roberto",
	})
	if err != nil {
		log.Printf("Error: %s", err)
	}

	msg2, err := capitalize(ctx, cancel, msg1)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	c.JSON(sink(ctx, cancel, msg2))
	return nil
}
