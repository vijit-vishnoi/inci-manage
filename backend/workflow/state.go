package workflow

import (
	"context"
	"errors"
)

var (
	ErrInvalidTransition = errors.New("invalid state transition")
	ErrIncompleteRCA     = errors.New("cannot close ticket: missing or incomplete RCA")
)

type RCAValidator func(ctx context.Context, workItemID int) error

// ValidateTransition returns an error if the transition is invalid based on strict rules.
func ValidateTransition(ctx context.Context, currentStatus, newStatus string, workItemID int, rcaValidator RCAValidator) error {
	switch currentStatus {
	case "OPEN":
		if newStatus != "INVESTIGATING" {
			return ErrInvalidTransition
		}
	case "INVESTIGATING":
		if newStatus != "RESOLVED" {
			return ErrInvalidTransition
		}
	case "RESOLVED":
		if newStatus != "CLOSED" {
			return ErrInvalidTransition
		}
		// Strict validation before closing: must have an RCA
		if rcaValidator != nil {
			if err := rcaValidator(ctx, workItemID); err != nil {
				return err
			}
		} else {
            return ErrIncompleteRCA
        }
	case "CLOSED":
		return ErrInvalidTransition
	default:
		return ErrInvalidTransition
	}
	return nil
}
