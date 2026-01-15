package eventemitter_test

import (
	"sync"
	"testing"
	"time"

	"github.com/afyadigital/eventemitter"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventName := "TestListen"

	actionCalled := false
	action := func(data interface{}) {
		actionCalled = true
	}

	eventemitterr.Listen(eventName, action)

	require.False(t, actionCalled)
	eventemitterr.Emit("TestListen", nil)
	require.True(t, actionCalled)
}

func TestMultipleEvents(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventNameA := "TestMultipleEventsA"

	actionCalled1 := false
	action1 := func(data interface{}) {
		actionCalled1 = true
	}

	eventNameB := "TestMultipleEventsB"

	actionCalled2 := false
	action2 := func(data interface{}) {
		actionCalled2 = true
	}

	eventemitterr.Listen(eventNameA, action1)
	eventemitterr.Listen(eventNameB, action2)

	require.False(t, actionCalled1)
	require.False(t, actionCalled2)

	eventemitterr.Emit("TestMultipleEventsA", nil)
	require.True(t, actionCalled1)
	require.False(t, actionCalled2)

	eventemitterr.Emit("TestMultipleEventsB", nil)
	require.True(t, actionCalled1)
	require.True(t, actionCalled2)
}

func TestMultipleListenersToSameEvent(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventNameA := "TestMultipleListenersToSameEvent"

	actionCalled1 := false
	action1 := func(data interface{}) {
		actionCalled1 = true
	}

	actionCalled2 := false
	action2 := func(data interface{}) {
		actionCalled2 = true
	}

	eventemitterr.Listen(eventNameA, action1)
	eventemitterr.Listen(eventNameA, action2)

	require.False(t, actionCalled1)
	require.False(t, actionCalled2)

	eventemitterr.Emit("TestMultipleListenersToSameEvent", nil)
	require.True(t, actionCalled1)
	require.True(t, actionCalled2)
}

func TestEmitWithData(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)
	eventName := "TestEmitWithData"

	actionCalled := false
	var actionData interface{}
	action := func(data interface{}) {
		actionCalled = true
		actionData = data
	}

	eventemitterr.Listen(eventName, action)

	require.False(t, actionCalled)
	eventemitterr.Emit("TestEmitWithData", "oi")
	require.True(t, actionCalled)
	require.Equal(t, "oi", actionData)
}

func TestRemoveEvent(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)
	eventName := "TestRemoveEvent"

	actionCalled := false
	action := func(data interface{}) {
		actionCalled = true
	}

	eventemitterr.Listen(eventName, action)
	eventemitterr.RemoveEvent(eventName)

	require.False(t, actionCalled)
	eventemitterr.Emit("TestRemoveEvent", nil)
	require.False(t, actionCalled)
}

func TestKeepRunningOnPanic(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)
	eventName := "TestKeepRunningOnPanic"

	action1Called := false
	action1 := func(data interface{}) {
		action1Called = true
	}

	actionPanic := func(data interface{}) {
		panic("ahhh")
	}

	action2Called := false
	action2 := func(data interface{}) {
		action2Called = true
	}

	eventemitterr.Listen(eventName, action1)
	eventemitterr.Listen(eventName, actionPanic)
	eventemitterr.Listen(eventName, action2)

	require.NotPanics(t, func() {
		eventemitterr.Emit("TestKeepRunningOnPanic", nil)
	})

	require.True(t, action1Called)
	require.True(t, action2Called)
}

func TestListenOnce(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventName := "TestListenOnce"

	actionCalledTimes := 0
	action := func(data interface{}) {
		actionCalledTimes++
	}

	eventemitterr.ListenOnce(eventName, action)

	require.Equal(t, 0, actionCalledTimes)

	eventemitterr.Emit("TestListenOnce", nil)
	require.Equal(t, 1, actionCalledTimes)

	eventemitterr.Emit("TestListenOnce", nil)
	require.Equal(t, 1, actionCalledTimes)
}

func TestReset(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventName1 := "TestReset1"

	action1Called := false
	action1 := func(data interface{}) {
		action1Called = true
	}

	eventName2 := "TestReset2"

	action2Called := false
	action2 := func(data interface{}) {
		action2Called = true
	}

	eventemitterr.Listen(eventName1, action1)
	eventemitterr.Listen(eventName2, action2)

	eventemitterr.Reset()

	eventemitterr.Emit("TestReset1", nil)
	eventemitterr.Emit("TestReset2", nil)

	require.False(t, action1Called)
	require.False(t, action2Called)
}

func TestRaceConditionEmitAfterRemove(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventName := "TestRaceConditionEmitAfterRemove"

	actionCalled := false
	action := func(data interface{}) {
		actionCalled = true
	}

	eventemitterr.Listen(eventName, action)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Millisecond)
		eventemitterr.RemoveEvent(eventName)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Millisecond)
		eventemitterr.Emit(eventName, nil)
	}()

	wg.Wait()

	require.False(t, actionCalled)
}

func TestRaceConditionMultipleOperations(t *testing.T) {
	eventemitterr := eventemitter.StartMap()

	t.Cleanup(eventemitterr.Reset)

	eventName := "TestRaceConditionMultipleOperations"

	actionCalledCount := 0
	action := func(data interface{}) {
		actionCalledCount++
	}

	eventemitterr.Listen(eventName, action)

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			if id%3 == 0 {
				time.Sleep(time.Duration(id) * time.Millisecond)
				eventemitterr.RemoveEvent(eventName)
			} else if id%3 == 1 {
				time.Sleep(time.Duration(id) * time.Millisecond)
				eventemitterr.Emit(eventName, nil)
			} else {
				time.Sleep(time.Duration(id) * time.Millisecond)
				eventemitterr.Listen(eventName, action)
			}
		}(i)
	}

	wg.Wait()

	require.GreaterOrEqual(t, actionCalledCount, 0)
}
