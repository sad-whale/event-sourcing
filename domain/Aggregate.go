package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

var factories = make(map[AggregateType]AggregateFactory)

func RegisterAggregate(aggregateType AggregateType, factory AggregateFactory) error {
	if factory == nil {
		return errors.New("nil factory")
	}

	if _, ok := factories[aggregateType]; ok {
		return fmt.Errorf("aggregate type %s already registered", aggregateType)
	}

	factories[aggregateType] = factory

	return nil
}

func CreateAggregate(aggregateType AggregateType) (AggregateRoot, error) {
	return CreateAggregateFromId(aggregateType, uuid.New())
}

func CreateAggregateFromId(aggregateType AggregateType, id uuid.UUID) (AggregateRoot, error) {
	return CreateAggregateFromIdAndVersion(aggregateType, id, 0)
}

func CreateAggregateFromIdAndVersion(aggregateType AggregateType, id uuid.UUID, version int32) (AggregateRoot, error) {
	factory, registered := factories[aggregateType]

	if !registered {
		return nil, fmt.Errorf("aggregate type %s is not registered", aggregateType)
	}

	base := AggregateRootBase{
		id:                id,
		version:           version,
		uncommittedEvents: make([]interface{}, 0),
		eventApplier:      nil,
	}

	var aggregate, isAggregateRoot = factory(base).(AggregateRoot)

	if !isAggregateRoot {
		return nil, fmt.Errorf("%s is not an AggregateRoot", aggregateType)
	}

	eventApplier, isEventApplier := aggregate.(EventApplier)

	if !isEventApplier {
		eventApplier = NewReflectEventApplier(aggregate)
	}

	base.setEventApplier(eventApplier)

	return aggregate, nil
}

type AggregateType string

type AggregateFactory func(base AggregateRootBase) interface{}

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
