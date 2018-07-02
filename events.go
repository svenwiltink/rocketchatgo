package rocketchatgo

import (
	"errors"
)

const (
	MessageCreate = "messageCreate"
	ChannelJoin   = "channelJoin"
	ChannelLeave  = "channelLeave"
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
	return MessageCreate
}

type MessageCreateEventHandler func(session *Session, event *MessageCreateEvent)

func (eh MessageCreateEventHandler) Type() string {
	return MessageCreate
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
	return ChannelJoin
}

type ChannelJoinEventHandler func(session *Session, event *ChannelJoinEvent)

func (eh ChannelJoinEventHandler) Type() string {
	return ChannelJoin
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
	return ChannelLeave
}

type ChannelLeaveEventHandler func(session *Session, event *ChannelLeaveEvent)

func (eh ChannelLeaveEventHandler) Type() string {
	return ChannelLeave
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
