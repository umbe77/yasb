package handlers

import (
	"context"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO: Implement Sample and example pipeline

type Message map[string]interface{}

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
				break
			case m, ok := <-msg:
				if ok {
					val, isok := m["name"]
					if !isok {
						cancel()
						return
					}
					m["name"] = strings.Title(val.(string))

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
