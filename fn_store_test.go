package pitlane_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nurburg-dev/pitlane"
	"github.com/stretchr/testify/require"
)

func Activity2_1(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func Activity2_2(_ context.Context, a, b string) (string, error) {
	return a + b, nil
}

func SampleWorkflow2(_ context.Context, name string, count int) (string, error) {
	return fmt.Sprintf("Hello %s %d", name, count), nil
}

func TestRegisterWorkflow(t *testing.T) {
	err := pitlane.RegisterWorkflow(SampleWorkflow2)
	require.NoError(t, err)

	err = pitlane.RegisterWorkflow(SampleWorkflow2)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already registered")
}

func TestRegisterActivity(t *testing.T) {
	err := pitlane.RegisterActivity(Activity2_1)
	require.NoError(t, err)

	err = pitlane.RegisterActivity(Activity2_2)
	require.NoError(t, err)
}

func TestRegisterActivity_Duplicate(t *testing.T) {
	activityFunc := func(_ context.Context) (interface{}, error) {
		return "activity result", nil
	}

	err := pitlane.RegisterActivity(activityFunc)
	require.NoError(t, err)

	err = pitlane.RegisterActivity(activityFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "already registered")
}

func TestRegisterWorkflow_InvalidFunction(t *testing.T) {
	// Test with non-function
	err := pitlane.RegisterWorkflow("not a function")
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be a function")

	// Test with function missing context parameter
	invalidFunc := func(a int) (int, error) {
		return a, nil
	}
	err = pitlane.RegisterWorkflow(invalidFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "first parameter must be context.Context")

	// Test with function returning wrong number of values
	invalidReturnFunc := func(_ context.Context, a int) int {
		return a
	}
	err = pitlane.RegisterWorkflow(invalidReturnFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must return exactly 2 values")
}
