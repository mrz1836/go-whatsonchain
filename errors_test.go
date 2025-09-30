package whatsonchain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorConstants tests all error constant definitions
func TestErrorConstants(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		err         error
		expectedMsg string
	}{
		{"ErrAddressNotFound", ErrAddressNotFound, "address not found"},
		{"ErrBlockNotFound", ErrBlockNotFound, "block not found"},
		{"ErrChainInfoNotFound", ErrChainInfoNotFound, "chain info not found"},
		{"ErrChainTipsNotFound", ErrChainTipsNotFound, "chain tips not found"},
		{"ErrPeerInfoNotFound", ErrPeerInfoNotFound, "peer info not found"},
		{"ErrExchangeRateNotFound", ErrExchangeRateNotFound, "exchange rate not found"},
		{"ErrMempoolInfoNotFound", ErrMempoolInfoNotFound, "mempool info not found"},
		{"ErrHeadersNotFound", ErrHeadersNotFound, "headers not found"},
		{"ErrScriptNotFound", ErrScriptNotFound, "script not found"},
		{"ErrTransactionNotFound", ErrTransactionNotFound, "transaction not found"},
		{"ErrMaxAddressesExceeded", ErrMaxAddressesExceeded, "max limit of addresses exceeded"},
		{"ErrMaxScriptsExceeded", ErrMaxScriptsExceeded, "max limit of scripts exceeded"},
		{"ErrBroadcastFailed", ErrBroadcastFailed, "error broadcasting transaction"},
		{"ErrMaxTransactionsExceeded", ErrMaxTransactionsExceeded, "max transactions limit exceeded"},
		{"ErrMaxPayloadSizeExceeded", ErrMaxPayloadSizeExceeded, "max overall payload size exceeded"},
		{"ErrMaxTransactionSizeExceeded", ErrMaxTransactionSizeExceeded, "max transaction size exceeded"},
		{"ErrMaxUTXOsExceeded", ErrMaxUTXOsExceeded, "max limit of UTXOs exceeded"},
		{"ErrMaxRawTransactionsExceeded", ErrMaxRawTransactionsExceeded, "max limit of raw transactions exceeded"},
		{"ErrMissingRequest", ErrMissingRequest, "missing request"},
		{"ErrBadRequest", ErrBadRequest, "bad request"},
		{"ErrBSVChainRequired", ErrBSVChainRequired, "operation is only available for BSV chain"},
		{"ErrBTCChainRequired", ErrBTCChainRequired, "operation is only available for BTC chain"},
		{"ErrTokenNotFound", ErrTokenNotFound, "token not found"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test that error implements error interface
			assert.Implements(t, (*error)(nil), tc.err)

			// Test error message
			assert.Equal(t, tc.expectedMsg, tc.err.Error())

			// Test that error is not nil
			assert.Error(t, tc.err)
		})
	}
}

// TestErrorComparison tests error equality and comparison
func TestErrorComparison(t *testing.T) {
	t.Parallel()

	// Test that different errors are not equal
	assert.NotEqual(t, ErrAddressNotFound, ErrBlockNotFound)
	assert.NotEqual(t, ErrBSVChainRequired, ErrBTCChainRequired)

	// Test error identity with errors.Is
	assert.NotErrorIs(t, ErrAddressNotFound, ErrBlockNotFound)
}

// TestChainSpecificErrors tests chain-specific error behavior
func TestChainSpecificErrors(t *testing.T) {
	t.Parallel()

	// Test BSV chain error
	assert.Contains(t, ErrBSVChainRequired.Error(), "BSV chain")
	assert.Contains(t, ErrBSVChainRequired.Error(), "operation is only available")

	// Test BTC chain error
	assert.Contains(t, ErrBTCChainRequired.Error(), "BTC chain")
	assert.Contains(t, ErrBTCChainRequired.Error(), "operation is only available")

	// Test chain errors are different
	assert.NotEqual(t, ErrBSVChainRequired, ErrBTCChainRequired)
}

// TestNotFoundErrors tests all "not found" type errors
func TestNotFoundErrors(t *testing.T) {
	t.Parallel()

	notFoundErrors := []error{
		ErrAddressNotFound,
		ErrBlockNotFound,
		ErrChainInfoNotFound,
		ErrChainTipsNotFound,
		ErrPeerInfoNotFound,
		ErrExchangeRateNotFound,
		ErrMempoolInfoNotFound,
		ErrHeadersNotFound,
		ErrScriptNotFound,
		ErrTransactionNotFound,
		ErrTokenNotFound,
	}

	for _, err := range notFoundErrors {
		assert.Contains(t, err.Error(), "not found")
	}
}

// TestLimitExceededErrors tests all limit exceeded errors
func TestLimitExceededErrors(t *testing.T) {
	t.Parallel()

	limitErrors := []error{
		ErrMaxAddressesExceeded,
		ErrMaxScriptsExceeded,
		ErrMaxTransactionsExceeded,
		ErrMaxPayloadSizeExceeded,
		ErrMaxTransactionSizeExceeded,
		ErrMaxUTXOsExceeded,
		ErrMaxRawTransactionsExceeded,
	}

	for _, err := range limitErrors {
		errMsg := err.Error()
		assert.True(t,
			contains(errMsg, "max") || contains(errMsg, "limit") || contains(errMsg, "exceeded"),
			"Expected error message to contain 'max', 'limit', or 'exceeded': %s", errMsg)
	}
}

// TestErrorWrapping tests error wrapping behavior
func TestErrorWrapping(t *testing.T) {
	t.Parallel()

	// Test custom wrapping
	customErr := fmt.Errorf("custom: %w", ErrBlockNotFound)
	if !errors.Is(customErr, ErrBlockNotFound) {
		t.Errorf("expected error to wrap ErrBlockNotFound")
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
