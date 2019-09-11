package domain

import (
	"fmt"
	"reflect"
	"strings"
)

type EventApplier interface {
	Apply(event interface{}) error
}

func NewReflectEventApplier(target interface{}) EventApplier {
	targetType := reflect.TypeOf(target)
	appliers := make(map[reflect.Type]reflect.Method)
	for i := 0; i < targetType.NumMethod(); i++ {
		method := targetType.Method(i)
		if strings.HasPrefix(method.Name, "apply") && method.Type.NumIn() == 2 {
			appliers[method.Type.In(1)] = method
		}
	}

	return reflectEventApplier{applyTarget: target, eventsAppliers: appliers}
}

type reflectEventApplier struct {
	applyTarget    interface{}
	eventsAppliers map[reflect.Type]reflect.Method
}

func (d reflectEventApplier) Apply(event interface{}) error {
	eventType := reflect.TypeOf(event)

	if applier, exists := d.eventsAppliers[eventType]; exists {
		input := []reflect.Value{reflect.ValueOf(d.applyTarget), reflect.ValueOf(event)}
		applier.Func.Call(input)
		return nil
	}

	return fmt.Errorf("applier for event type %s not found", eventType)
}
