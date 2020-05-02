package env

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
)

func Test_FatalOnError(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		err error
	}{
		"nil": {},
		"err": {fmt.Errorf("error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var logged string
			var canceled bool
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			if tc.err != nil {
				logger.EXPECT().Error(tc.err).Do(func(err error) {
					logged = err.Error()
				}).Times(1)
			}
			e := &env{
				logger:        logger,
				cancelContext: func() { canceled = true },
			}
			e.FatalOnError(tc.err)
			if tc.err != nil {
				assert.Equal(t, logged, tc.err.Error())
				assert.True(t, canceled)
			} else {
				assert.Empty(t, logged)
				assert.False(t, canceled)
			}
		})
	}
}

func Test_PrintVersion(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		program        string
		commandVersion string
		commandErr     error
	}{
		"no data": {},
		"data":    {"binu", "2.3-5", nil},
		"error":   {"binu", "", fmt.Errorf("error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var logged string
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			if tc.commandErr != nil {
				logger.EXPECT().Error(tc.commandErr).Do(func(err error) {
					logged = err.Error()
				}).Times(1)
			} else {
				logger.EXPECT().Info("%s version: %s", tc.program, tc.commandVersion).
					Do(func(format, program, version string) {
						logged = fmt.Sprintf(format, program, version)
					}).Times(1)
			}
			e := &env{logger: logger}
			commandFn := func(ctx context.Context) (string, error) { return tc.commandVersion, tc.commandErr }
			e.PrintVersion(context.Background(), tc.program, commandFn)
			if tc.commandErr != nil {
				assert.Equal(t, logged, tc.commandErr.Error())
			} else {
				assert.Equal(t, logged, fmt.Sprintf("%s version: %s", tc.program, tc.commandVersion))
			}
		})
	}
}
