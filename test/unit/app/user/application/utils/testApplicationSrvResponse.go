package utils

import (
	returnLog "github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/message"
	"github.com/rcrespodev/user_manager/pkg/kernel/cqrs/returnLog/domain/valueObjects"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type Want struct {
	Status         valueObjects.Status
	HttpCodeReturn valueObjects.HttpCodeReturn
	Error          error
	ErrorMessage   *message.MessageData
	SuccessMessage *message.MessageData
}

func (w *Want) TestResponse(t *testing.T, log *returnLog.ReturnLog) {
	// ReturnLog check
	require.EqualValues(t, w.Status, log.Status())
	require.EqualValues(t, w.HttpCodeReturn, log.HttpCode())

	// Check Internal error
	switch w.Error {
	case nil:
		if log.Error() != nil {
			require.Nil(t, log.Error().InternalError())
		}
	default:
		require.EqualValues(t, w.Error, log.Error().InternalError().Error())
	}

	// Check Client error messages
	switch w.ErrorMessage {
	case nil:
		if log.Error() != nil {
			require.Nil(t, log.Error().Message())
		}
	default:
		gotMessage := log.Error().Message()
		gotMessage.Time = time.Time{}
		require.EqualValues(t, w.ErrorMessage, gotMessage)
	}

	// Check Success message
	switch w.SuccessMessage {
	case nil:
		require.Nil(t, log.Success())
	default:
		gotMessage := log.Success().MessageData()
		gotMessage.Time = w.SuccessMessage.Time
		require.EqualValues(t, w.SuccessMessage, gotMessage)
	}
}
