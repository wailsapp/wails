package runtime

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	internallogger "github.com/wailsapp/wails/v2/internal/logger"
	"github.com/wailsapp/wails/v2/internal/servicebus"
	"github.com/wailsapp/wails/v2/pkg/logger"

	is2 "github.com/matryer/is"
)

func TestStoreProvider_NewWithScalarDefault(t *testing.T) {
	is := is2.New(t)

	defaultLogger := logger.NewDefaultLogger()
	testLogger := internallogger.New(defaultLogger)
	//testLogger.SetLogLevel(logger.TRACE)
	serviceBus := servicebus.New(testLogger)
	err := serviceBus.Start()
	is.NoErr(err)
	defer serviceBus.Stop()

	testRuntime := New(serviceBus)
	storeProvider := newStore(testRuntime)
	testStore := storeProvider.New("test", 100)
	value := testStore.Get()
	is.Equal(value, 100)
	testStore.resync()
	value = testStore.Get()
	is.Equal(value, 100)
}

func TestStoreProvider_NewWithStructDefault(t *testing.T) {
	is := is2.New(t)

	defaultLogger := logger.NewDefaultLogger()
	testLogger := internallogger.New(defaultLogger)
	//testLogger.SetLogLevel(logger.TRACE)
	serviceBus := servicebus.New(testLogger)
	err := serviceBus.Start()
	is.NoErr(err)
	defer serviceBus.Stop()

	testRuntime := New(serviceBus)
	storeProvider := newStore(testRuntime)

	type TestValue struct {
		Name string
	}
	testValue := &TestValue{
		Name: "hi",
	}

	testStore := storeProvider.New("test", testValue)

	testStore.Update(func(current *TestValue) *TestValue {
		return testValue
	})
	testStore.resync()
	value := testStore.Get()
	is.Equal(value, testValue)
	is.Equal(value.(*TestValue).Name, "hi")

	testValue = &TestValue{
		Name: "there",
	}
	testStore.Update(func(current *TestValue) *TestValue {
		return testValue
	})
	testStore.resync()
	value = testStore.Get()
	is.Equal(value, testValue)
	is.Equal(value.(*TestValue).Name, "there")

}

func TestStoreProvider_RapidReadWrite(t *testing.T) {
	is := is2.New(t)

	defaultLogger := logger.NewDefaultLogger()
	testLogger := internallogger.New(defaultLogger)
	//testLogger.SetLogLevel(logger.TRACE)
	serviceBus := servicebus.New(testLogger)
	err := serviceBus.Start()
	is.NoErr(err)
	defer serviceBus.Stop()

	testRuntime := New(serviceBus)
	storeProvider := newStore(testRuntime)

	testStore := storeProvider.New("test", 1)

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	var wg sync.WaitGroup
	readers := 100
	writers := 100
	wg.Add(readers + writers)
	// Setup readers
	go func(testStore *Store, ctx context.Context) {
		for readerCount := 0; readerCount < readers; readerCount++ {
			go func(store *Store, ctx context.Context, id int) {
				for {
					select {
					case <-ctx.Done():
						wg.Done()
						return
					default:
						store.Get()
					}
				}
			}(testStore, ctx, readerCount)
		}
	}(testStore, ctx)

	// Setup writers
	go func(testStore *Store, ctx context.Context) {
		for writerCount := 0; writerCount < writers; writerCount++ {
			go func(store *Store, ctx context.Context, id int) {
				for {
					select {
					case <-ctx.Done():
						wg.Done()
						return
					default:
						store.Update(func(current int) int {
							return rand.Int()
						})
					}
				}
			}(testStore, ctx, writerCount)
		}
	}(testStore, ctx)

	wg.Wait()
}
