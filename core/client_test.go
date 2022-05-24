package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types"
)

func TestRemoteClient_Status(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	_, client := StartTestClient(ctx, t)
	status, err := client.Status(ctx)
	require.NoError(t, err)
	require.NotNil(t, status)
}

func TestRemoteClient_StartBlockSubscription_And_GetBlock(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	_, client := StartTestClient(ctx, t)
	eventChan, err := client.Subscribe(ctx, newBlockSubscriber, newBlockEventQuery)
	require.NoError(t, err)

	for i := 1; i <= 3; i++ {
		select {
		case evt := <-eventChan:
			// check that `Block` works as intended (passing nil to get block at latest height)
			h := evt.Data.(types.EventDataNewBlock).Block.Height
			block, err := client.Block(ctx, &h)
			require.NoError(t, err)
			require.Equal(t, int64(i), block.Block.Height)
		case <-ctx.Done():
			require.NoError(t, ctx.Err())
		}
	}
	// unsubscribe to event channel
	require.NoError(t, client.Unsubscribe(ctx, newBlockSubscriber, newBlockEventQuery))
}
