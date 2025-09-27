package pitlane

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/nurburg-dev/pitlane/internal/db"
	"github.com/nurburg-dev/pitlane/internal/dbrepo"
	"github.com/nurburg-dev/pitlane/internal/entities"
	"github.com/nurburg-dev/pitlane/internal/utils"
)

func ExecuteActivity(ctx context.Context, activityFunc, retVal interface{}, args ...interface{}) error {
	activityName, err := utils.GetFunctionName(activityFunc)
	if err != nil {
		return fmt.Errorf("failed to get activity function name: %w", err)
	}

	activityHistory, ok := ctx.Value(activityHistoryKey).([]*entities.DBActivityRun)
	if !ok {
		return fmt.Errorf("activity history not found in context")
	}

	activityRepo, ok := ctx.Value(activityRepoKey).(*dbrepo.PGActivityRunRepository)
	if !ok {
		return fmt.Errorf("activity repository not found in context")
	}

	workflowRunID, ok := ctx.Value(workflowRunIDKey).(string)
	if !ok {
		return fmt.Errorf("workflow run ID not found in context")
	}

	activityHistoryIndexPtr, ok := ctx.Value(activityHistoryIndexKey).(*[]int)
	if !ok {
		return fmt.Errorf("activity history index not found in context")
	}
	activityHistoryIndex := (*activityHistoryIndexPtr)[0]

	if activityHistoryIndex < len(activityHistory) {
		currentActivity := activityHistory[activityHistoryIndex]

		if currentActivity.ActivityName != activityName {
			return fmt.Errorf("%w: expected activity %s, got %s",
				ErrActivityHistoryMismatch, activityName, currentActivity.ActivityName)
		}

		var expectedInput []interface{}
		if unmarshalErr := json.Unmarshal(currentActivity.Input, &expectedInput); unmarshalErr != nil {
			return fmt.Errorf("failed to unmarshal activity input: %w", unmarshalErr)
		}

		if !reflect.DeepEqual(args, expectedInput) {
			return fmt.Errorf("%w: activity input mismatch for %s",
				ErrActivityHistoryMismatch, activityName)
		}

		if currentActivity.Status == entities.ActivityStatusFinished && currentActivity.Output != nil {
			if unmarshalErr := json.Unmarshal(*currentActivity.Output, retVal); unmarshalErr != nil {
				return fmt.Errorf("failed to unmarshal activity output: %w", unmarshalErr)
			}
			(*activityHistoryIndexPtr)[0]++
			return nil
		}

		if currentActivity.Status == entities.ActivityStatusFailed {
			(*activityHistoryIndexPtr)[0]++
			return fmt.Errorf("activity %s failed: %s", activityName,
				getStringValue(currentActivity.ErrorMessage))
		}
	}

	inputBytes, err := json.Marshal(args)
	if err != nil {
		return fmt.Errorf("failed to marshal activity input: %w", err)
	}

	now := time.Now()
	activityRunID := db.GenerateReadableID()
	newActivityRun := &entities.DBActivityRun{
		ID:            activityRunID,
		ActivityName:  activityName,
		WorkflowRunID: workflowRunID,
		Input:         inputBytes,
		Status:        entities.ActivityStatusPending,
		ScheduledAt:   now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := activityRepo.CreateActivityRun(ctx, newActivityRun); err != nil {
		return fmt.Errorf("failed to create activity run: %w", err)
	}

	return ErrActivitySchedulingNeeded
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
