package domain

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type AggregateFactoryMethod func(base *AggregateRootBase) AggregateRoot

type AggregateFactory interface {
	RegisterAggregate(factory AggregateFactoryMethod) error
	CreateAggregate(aggregate *AggregateRoot) error
	CreateAggregateFromId(aggregate *AggregateRoot, id uuid.UUID) error
	CreateAggregateFromIdAndVersion(aggregate *AggregateRoot, id uuid.UUID, version int32) error
}

type mapAggregateFactory map[reflect.Type]AggregateFactoryMethod

func (factories *mapAggregateFactory) RegisterAggregate(factory AggregateFactoryMethod) error {
	if factory == nil {
		return errors.New("nil factory")
	}

	aggregateType := reflect.TypeOf(factory).Out(0)

	if _, ok := (*factories)[aggregateType]; ok {
		return fmt.Errorf("aggregate type %s already registered", aggregateType)
	}

	(*factories)[aggregateType] = factory

	return nil
}

func (factories *mapAggregateFactory) CreateAggregate(aggregate *AggregateRoot) error {
	return factories.CreateAggregateFromId(aggregate, uuid.New())
}

func (factories *mapAggregateFactory) CreateAggregateFromId(aggregate *AggregateRoot, id uuid.UUID) error {
	return factories.CreateAggregateFromIdAndVersion(aggregate, id, 0)
}

func (factories *mapAggregateFactory) CreateAggregateFromIdAndVersion(aggregate *AggregateRoot, id uuid.UUID, version int32) error {
	aggregateType := reflect.TypeOf(aggregate)
	factory, registered := (*factories)[aggregateType]

	if !registered {
		return fmt.Errorf("aggregate type %s is not registered", aggregateType)
	}

	base := &AggregateRootBase{
		id:                id,
		version:           version,
		uncommittedEvents: make([]interface{}, 0),
		eventApplier:      nil,
	}

	*aggregate = factory(base)

	eventApplier := newReflectEventApplier(aggregate)


	base.setEventApplier(eventApplier)

	return nil
}

func newMapAggregateFactory() AggregateFactory  {
	factory := mapAggregateFactory(make(map[reflect.Type]AggregateFactoryMethod))
	return &factory
}

