package pitlane

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
)

type Workflow struct {
	workflowFunc     interface{}
	workflowFuncName string
	activityFuncMaps map[string]interface{}
}

func NewWorkflow(wfFunc interface{}) (*Workflow, error) {
	if err := validateFunc(wfFunc, "workflow"); err != nil {
		return nil, err
	}

	return &Workflow{
		workflowFunc:     wfFunc,
		activityFuncMaps: make(map[string]interface{}),
		workflowFuncName: runtime.FuncForPC(reflect.ValueOf(wfFunc).Pointer()).Name(),
	}, nil
}

func (w *Workflow) AddActivities(activityFuncs ...interface{}) error {
	for _, activityFunc := range activityFuncs {
		if err := validateFunc(activityFunc, "activity"); err != nil {
			return err
		}

		funcName := runtime.FuncForPC(reflect.ValueOf(activityFunc).Pointer()).Name()
		_, exists := w.activityFuncMaps[funcName]
		if exists {
			return fmt.Errorf("duplicate activity found in workflow %s", funcName)
		}
		w.activityFuncMaps[funcName] = activityFunc
	}
	return nil
}

func validateFunc(fn interface{}, funcType string) error {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("%s must be a function, got %T", funcType, fn)
	}

	if fnType.NumIn() < 1 {
		return fmt.Errorf("%s function must have at least one parameter (context.Context)", funcType)
	}

	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !fnType.In(0).Implements(contextType) {
		return fmt.Errorf("first parameter must be context.Context")
	}

	if fnType.NumOut() != 2 {
		return fmt.Errorf("%s function must return exactly 2 values (result, error)", funcType)
	}

	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !fnType.Out(1).Implements(errorType) {
		return fmt.Errorf("second return value must be error")
	}

	return nil
}
