package rocketchatgo

import (
	"errors"
)

type Event interface {
	GetType() string
}

// EventHandler is an interface for Discord events.
type EventHandler interface {
	// Type returns the type of event this handler belongs to.
	Type() string

	// Handle is called whenever an event of Type() happens.
	// It is the receivers responsibility to type assert that the interface
	// is the expected struct.
	Handle(*Session, interface{})
}

type MessageCreateEvent struct {
	Message *Message
}

func (*MessageCreateEvent) GetType() string {
	return EventMessageCreate
}

type MessageCreateEventHandler func(session *Session, event *MessageCreateEvent)

func (eh MessageCreateEventHandler) Type() string {
	return EventMessageCreate
}

func (eh MessageCreateEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*MessageCreateEvent); ok {
		eh(s, t)
	}
}

type ChannelJoinEvent struct {
	Channel *Room
}

func (*ChannelJoinEvent) GetType() string {
	return EventChannelJoin
}

type ChannelJoinEventHandler func(session *Session, event *ChannelJoinEvent)

func (eh ChannelJoinEventHandler) Type() string {
	return EventChannelJoin
}

func (eh ChannelJoinEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelJoinEvent); ok {
		eh(s, t)
	}
}

type ChannelLeaveEvent struct {
	Channel *Room
}

func (*ChannelLeaveEvent) GetType() string {
	return EventChannelLeave
}

type ChannelLeaveEventHandler func(session *Session, event *ChannelLeaveEvent)

func (eh ChannelLeaveEventHandler) Type() string {
	return EventChannelLeave
}

func (eh ChannelLeaveEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*ChannelLeaveEvent); ok {
		eh(s, t)
	}
}

func getHandlerFromInterface(handler interface{}) (EventHandler, error) {
	switch v := handler.(type) {
	case func(session *Session, event *MessageCreateEvent):
		return MessageCreateEventHandler(v), nil
	case func(session *Session, event *ChannelJoinEvent):
		return ChannelJoinEventHandler(v), nil
	case func(session *Session, event *ChannelLeaveEvent):
		return ChannelLeaveEventHandler(v), nil
	default:
		return nil, errors.New("invalid handler")
	}
}
