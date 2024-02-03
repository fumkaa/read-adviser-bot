package telegram

import (
	"errors"
	"fmt"

	"github.com/fumkaa/read-adviser-bot/clients/telegram"
	"github.com/fumkaa/read-adviser-bot/events"
	"github.com/fumkaa/read-adviser-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatId   int
	UserName string
}

var (
	ErrUnknowEventType = errors.New("unknow event type")
	ErrUnknowMetaType  = errors.New("unknow meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("[Fetch]can't get updetes: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.ProcessMessage(event)
	default:
		return fmt.Errorf("can't process message: %w", ErrUnknowEventType)
	}
}

func (p *Processor) ProcessMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't process message: %w", err)
	}

	if err := p.doCmd(event.Text, meta.ChatId, meta.UserName); err != nil {
		return fmt.Errorf("can't doCmd: %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, fmt.Errorf("can't get meta: %w", ErrUnknowMetaType)
	}
	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = Meta{
			ChatId:   upd.Message.Chat.ChatId,
			UserName: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.UnknowType
	}
	return events.Message
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
