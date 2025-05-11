package logar

import (
	"fmt"
	"reflect"
)

type ActionManager interface {
	Common

	InvokeAction(path string, args ...any) ([]any, error)
	GetActionArgTypes(path string) ([]reflect.Type, error)
	GetActionsMap() Actions
	GetAllActions() []string
	GetActionDetails(path string) (Action, bool)
	AddAction(action Action)
	RemoveAction(path string)
}

type ActionManagerImpl struct {
	core *AppImpl
}

func (a *ActionManagerImpl) GetApp() App {
	return a.core
}

func (a *ActionManagerImpl) InvokeAction(path string, args ...any) ([]any, error) {
	action, ok := a.GetActionDetails(path)
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

func (a *ActionManagerImpl) GetActionArgTypes(path string) ([]reflect.Type, error) {
	action, ok := a.GetActionDetails(path)
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

func (a *ActionManagerImpl) GetActionsMap() Actions {
	return a.core.actions
}

func (a *ActionManagerImpl) GetAllActions() []string {
	actions := []string{}
	for _, action := range a.core.actions {
		actions = append(actions, action.Path)
	}
	return actions
}

func (a *ActionManagerImpl) GetActionDetails(path string) (Action, bool) {
	for _, action := range a.core.actions {
		if action.Path == path {
			return action, true
		}
	}
	return Action{}, false
}

func (a *ActionManagerImpl) AddAction(action Action) {
	for i, existingAction := range a.core.actions {
		if existingAction.Path == action.Path {
			a.core.actions[i] = action
			return
		}
	}
	a.core.actions = append(a.core.actions, action)
}

func (a *ActionManagerImpl) RemoveAction(path string) {
	for i, action := range a.core.actions {
		if action.Path == path {
			a.core.actions = append(a.core.actions[:i], a.core.actions[i+1:]...)
			return
		}
	}
}
