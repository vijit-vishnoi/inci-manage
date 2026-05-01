package workflow

import (
	"context"
	"testing"
)

func TestValidateTransition(t *testing.T) {
	ctx := context.Background()

	// Mock validators
	passValidator := func(ctx context.Context, id int) error { return nil }
	failValidator := func(ctx context.Context, id int) error { return ErrIncompleteRCA }

	tests := []struct {
		name          string
		current       string
		next          string
		validator     RCAValidator
		expectedError error
	}{
		{"Valid: OPEN -> INVESTIGATING", "OPEN", "INVESTIGATING", passValidator, nil},
		{"Invalid: OPEN -> RESOLVED", "OPEN", "RESOLVED", passValidator, ErrInvalidTransition},
		{"Invalid: OPEN -> CLOSED", "OPEN", "CLOSED", passValidator, ErrInvalidTransition},
		
		{"Valid: INVESTIGATING -> RESOLVED", "INVESTIGATING", "RESOLVED", passValidator, nil},
		{"Invalid: INVESTIGATING -> CLOSED", "INVESTIGATING", "CLOSED", passValidator, ErrInvalidTransition},
		
		{"Valid: RESOLVED -> CLOSED (With RCA)", "RESOLVED", "CLOSED", passValidator, nil},
		{"Invalid: RESOLVED -> CLOSED (No RCA)", "RESOLVED", "CLOSED", failValidator, ErrIncompleteRCA},
		{"Invalid: RESOLVED -> INVESTIGATING", "RESOLVED", "INVESTIGATING", passValidator, ErrInvalidTransition},

		{"Invalid: CLOSED -> OPEN", "CLOSED", "OPEN", passValidator, ErrInvalidTransition},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransition(ctx, tt.current, tt.next, 1, tt.validator)
			if err != tt.expectedError {
				t.Errorf("expected error: %v, got: %v", tt.expectedError, err)
			}
		})
	}
}
