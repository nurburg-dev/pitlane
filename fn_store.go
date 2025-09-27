package pitlane

import (
	"fmt"

	"github.com/nurburg-dev/pitlane/internal/utils"
)

var (
	workflowStore map[string]any = map[string]any{}
	activityStore map[string]any = map[string]any{}
)

func RegisterWorkflow(workflowFunc interface{}) error {
	funcName, err := utils.GetFunctionName(workflowFunc)
	if err != nil {
		return fmt.Errorf("failed to get workflow function name: %w", err)
	}

	_, exists := workflowStore[funcName]
	if exists {
		return fmt.Errorf("workflow %s already registered", funcName)
	}

	workflowStore[funcName] = workflowFunc
	return nil
}

func RegisterActivity(activityFunc interface{}) error {
	funcName, err := utils.GetFunctionName(activityFunc)
	if err != nil {
		return fmt.Errorf("failed to get activity function name: %w", err)
	}

	_, exists := activityStore[funcName]
	if exists {
		return fmt.Errorf("activity %s already registered", funcName)
	}

	activityStore[funcName] = activityFunc
	return nil
}

func GetWorkflowStore() map[string]any {
	return workflowStore
}

func GetActivityStore() map[string]any {
	return activityStore
}
