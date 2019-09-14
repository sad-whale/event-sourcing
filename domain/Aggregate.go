package domain

import (
	"errors"
	"github.com/google/uuid"
)

type AggregateType string

type AggregateRoot interface {
	Id() uuid.UUID
	Version() int32
	UncommittedEvents() []interface{}
	Commit(version int32) error
}

type AggregateRootBase struct {
	id                uuid.UUID
	version           int32
	uncommittedEvents []interface{}
	eventApplier      EventApplier
}

func (arb *AggregateRootBase) setEventApplier(applier EventApplier) {
	arb.eventApplier = applier
}

func (arb AggregateRootBase) Id() uuid.UUID {
	return arb.id
}

func (arb AggregateRootBase) Version() int32 {
	return arb.version
}

func (arb AggregateRootBase) UncommittedEvents() []interface{} {
	return arb.uncommittedEvents
}

func (arb *AggregateRootBase) Commit(version int32) error {
	if version <= arb.version {
		return errors.New("invalid version")
	}

	arb.version = version
	arb.uncommittedEvents = make([]interface{}, 0)

	return nil
}
