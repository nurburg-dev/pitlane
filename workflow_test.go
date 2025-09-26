package pitlane_test

import (
	"context"
	"testing"

	"github.com/nurburg-dev/pitlane"
	"github.com/stretchr/testify/require"
)

func Activity1(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func Activity2(_ context.Context, a, b string) (string, error) {
	return a + b, nil
}

func SampleWorkflow(_ context.Context, name string) (string, error) {
	return "Hello " + name, nil
}

func TestNewWorkflow(t *testing.T) {
	workflow, err := pitlane.NewWorkflow(SampleWorkflow)
	require.NoError(t, err)
	require.NotNil(t, workflow)
}

func TestWorkflow_AddActivities(t *testing.T) {
	workflow, err := pitlane.NewWorkflow(SampleWorkflow)
	require.NoError(t, err)

	err = workflow.AddActivities(Activity1, Activity2)
	require.NoError(t, err)
}

func TestWorkflow_AddActivities_Duplicate(t *testing.T) {
	workflow, err := pitlane.NewWorkflow(SampleWorkflow)
	require.NoError(t, err)

	activityFunc := func(_ context.Context) (interface{}, error) {
		return "activity result", nil
	}

	err = workflow.AddActivities(activityFunc)
	require.NoError(t, err)

	err = workflow.AddActivities(activityFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "duplicate activity found in workflow")
}

func TestWorkflow_AddActivities_InvalidFunction(t *testing.T) {
	workflow, err := pitlane.NewWorkflow(SampleWorkflow)
	require.NoError(t, err)

	// Test with non-function
	err = workflow.AddActivities("not a function")
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be a function")

	// Test with function missing context parameter
	invalidFunc := func(a int) (int, error) {
		return a, nil
	}
	err = workflow.AddActivities(invalidFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "first parameter must be context.Context")

	// Test with function returning wrong number of values
	invalidReturnFunc := func(_ context.Context, a int) int {
		return a
	}
	err = workflow.AddActivities(invalidReturnFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must return exactly 2 values")
}

func TestNewWorkflow_InvalidFunction(t *testing.T) {
	// Test with non-function
	_, err := pitlane.NewWorkflow("not a function")
	require.Error(t, err)
	require.Contains(t, err.Error(), "must be a function")

	// Test with function missing context parameter
	invalidFunc := func(a int) (int, error) {
		return a, nil
	}
	_, err = pitlane.NewWorkflow(invalidFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "first parameter must be context.Context")

	// Test with function returning wrong number of values
	invalidReturnFunc := func(_ context.Context, a int) int {
		return a
	}
	_, err = pitlane.NewWorkflow(invalidReturnFunc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "must return exactly 2 values")
}
