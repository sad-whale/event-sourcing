package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type AggregateFactoryMethod func(base *AggregateRootBase) AggregateRoot

type AggregateFactory interface {
	RegisterAggregate(aggregateType AggregateType, factory AggregateFactoryMethod) error
	CreateAggregate(aggregateType AggregateType) (AggregateRoot, error)
	CreateAggregateFromId(aggregateType AggregateType, id uuid.UUID) (AggregateRoot, error)
	CreateAggregateFromIdAndVersion(aggregateType AggregateType, id uuid.UUID, version int32) (AggregateRoot, error)
}

type aggregateFactory map[AggregateType]AggregateFactoryMethod

func (factories *aggregateFactory) RegisterAggregate(aggregateType AggregateType, factory AggregateFactoryMethod) error {
	if factory == nil {
		return errors.New("nil factory")
	}

	if _, ok := (*factories)[aggregateType]; ok {
		return fmt.Errorf("aggregate type %s already registered", aggregateType)
	}

	(*factories)[aggregateType] = factory

	return nil
}

func (factories *aggregateFactory) CreateAggregate(aggregateType AggregateType) (AggregateRoot, error) {
	return factories.CreateAggregateFromId(aggregateType, uuid.New())
}

func (factories *aggregateFactory) CreateAggregateFromId(aggregateType AggregateType, id uuid.UUID) (AggregateRoot, error) {
	return factories.CreateAggregateFromIdAndVersion(aggregateType, id, 0)
}

func (factories *aggregateFactory) CreateAggregateFromIdAndVersion(aggregateType AggregateType, id uuid.UUID, version int32) (AggregateRoot, error) {
	factory, registered := (*factories)[aggregateType]

	if !registered {
		return nil, fmt.Errorf("aggregate type %s is not registered", aggregateType)
	}

	base := &AggregateRootBase{
		id:                id,
		version:           version,
		uncommittedEvents: make([]interface{}, 0),
		eventApplier:      nil,
	}

	aggregate := factory(base)

	eventApplier, isEventApplier := aggregate.(EventApplier)

	if !isEventApplier {
		eventApplier = NewReflectEventApplier(aggregate)
	}

	base.setEventApplier(eventApplier)

	return aggregate, nil
}

var factory = aggregateFactory(make(map[AggregateType]AggregateFactoryMethod))

func GetAggregateFactory() AggregateFactory {
	return &factory
}
