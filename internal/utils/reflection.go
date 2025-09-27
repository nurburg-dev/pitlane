package utils

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
)

func GetFunctionName(f any) (string, error) {
	if err := validateFunc(f); err != nil {
		return "", err
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), nil
}

func validateFunc(fn any) error {
	fnType := reflect.TypeOf(fn)

	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("must be a function, got %T", fn)
	}
	if fnType.NumIn() < 1 {
		return fmt.Errorf("function must have at least one parameter (context.Context)")
	}
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if !fnType.In(0).Implements(contextType) {
		return fmt.Errorf("first parameter must be context.Context")
	}
	if fnType.NumOut() != 2 {
		return fmt.Errorf("function must return exactly 2 values (result, error)")
	}
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !fnType.Out(1).Implements(errorType) {
		return fmt.Errorf("second return value must be error")
	}
	return nil
}

func ValidateArgs(fn any, args ...any) error {
	if err := validateFunc(fn); err != nil {
		return err
	}

	fnType := reflect.TypeOf(fn)

	// Check argument count (excluding context.Context which is the first parameter)
	expectedArgCount := fnType.NumIn() - 1
	actualArgCount := len(args)

	if actualArgCount != expectedArgCount {
		return fmt.Errorf("function expects %d arguments (excluding context), got %d", expectedArgCount, actualArgCount)
	}

	// Validate each argument type
	for i, arg := range args {
		paramIndex := i + 1 // Skip context.Context parameter
		expectedType := fnType.In(paramIndex)
		actualType := reflect.TypeOf(arg)

		if actualType == nil {
			if expectedType.Kind() == reflect.Ptr || expectedType.Kind() == reflect.Interface {
				continue // nil is acceptable for pointers and interfaces
			}
			return fmt.Errorf("argument %d: expected %s, got nil", i, expectedType)
		}

		if !actualType.AssignableTo(expectedType) {
			return fmt.Errorf("argument %d: expected %s, got %s", i, expectedType, actualType)
		}
	}

	return nil
}

func InvokeFunction(fn any, args ...any) ([]reflect.Value, error) {
	if err := validateFunc(fn); err != nil {
		return nil, err
	}

	fnValue := reflect.ValueOf(fn)
	fnType := fnValue.Type()

	// Check argument count
	if len(args) != fnType.NumIn() {
		return nil, fmt.Errorf("function expects %d arguments, got %d", fnType.NumIn(), len(args))
	}

	// Convert arguments to reflect.Values and validate types
	reflectArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		expectedType := fnType.In(i)

		if arg == nil {
			if expectedType.Kind() == reflect.Ptr || expectedType.Kind() == reflect.Interface {
				reflectArgs[i] = reflect.Zero(expectedType)
			} else {
				return nil, fmt.Errorf("argument %d: cannot pass nil to non-pointer/interface type %s", i, expectedType)
			}
		} else {
			argValue := reflect.ValueOf(arg)
			if !argValue.Type().AssignableTo(expectedType) {
				return nil, fmt.Errorf("argument %d: expected %s, got %s", i, expectedType, argValue.Type())
			}
			reflectArgs[i] = argValue
		}
	}

	// Call the function
	results := fnValue.Call(reflectArgs)
	return results, nil
}
