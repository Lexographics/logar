package logar

import (
	"fmt"
	"reflect"
)

type ActionManager interface {
	InvokeAction(path string, args ...any) ([]any, error)
	GetActionArgTypes(path string) ([]reflect.Type, error)
	GetActionsMap() Actions
	GetAllActions() []string
	GetActionDetails(path string) (Action, bool)
	AddAction(action Action)
	RemoveAction(path string)
}

func (l *AppImpl) InvokeAction(path string, args ...any) ([]any, error) {
	action, ok := l.GetActionDetails(path)
	if !ok {
		return nil, fmt.Errorf("action '%s' not found", path)
	}

	actionFunc := reflect.ValueOf(action.Func)
	if actionFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("path '%s' does not point to a function", path)
	}

	actionType := actionFunc.Type()
	if actionType.NumIn() != len(args) && !(actionType.IsVariadic() && len(args) >= actionType.NumIn()-1) {
		return nil, fmt.Errorf("path '%s' expects %d arguments, got %d", path, actionType.NumIn(), len(args))
	}

	inArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		inArgs[i] = reflect.ValueOf(arg)
	}

	out := actionFunc.Call(inArgs)

	result := make([]any, len(out))
	for i, val := range out {
		result[i] = val.Interface()
	}

	return result, nil
}

func (l *AppImpl) GetActionArgTypes(path string) ([]reflect.Type, error) {
	action, ok := l.GetActionDetails(path)
	if !ok {
		return nil, fmt.Errorf("action '%s' not found", path)
	}

	actionFunc := reflect.ValueOf(action.Func)
	if actionFunc.Kind() != reflect.Func {
		return nil, fmt.Errorf("path '%s' does not point to a function", path)
	}

	actionType := actionFunc.Type()
	numArgs := actionType.NumIn()
	argTypes := make([]reflect.Type, numArgs)
	for i := 0; i < numArgs; i++ {
		argTypes[i] = actionType.In(i)
	}

	return argTypes, nil
}

func (l *AppImpl) GetActionsMap() Actions {
	return l.actions
}

func (l *AppImpl) GetAllActions() []string {
	actions := []string{}
	for _, action := range l.actions {
		actions = append(actions, action.Path)
	}
	return actions
}

func (l *AppImpl) GetActionDetails(path string) (Action, bool) {
	for _, action := range l.actions {
		if action.Path == path {
			return action, true
		}
	}
	return Action{}, false
}

func (l *AppImpl) AddAction(action Action) {
	for i, existingAction := range l.actions {
		if existingAction.Path == action.Path {
			l.actions[i] = action
			return
		}
	}
	l.actions = append(l.actions, action)
}

func (l *AppImpl) RemoveAction(path string) {
	for i, action := range l.actions {
		if action.Path == path {
			l.actions = append(l.actions[:i], l.actions[i+1:]...)
			return
		}
	}
}
