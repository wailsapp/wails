package mailbox

import (
	"sync"
	"testing"
)

func TestMailboxProcessesMessagesInOrder(t *testing.T) {
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	consumed := make(chan int, 3)

	messages := New(func(message int) {
		if message == 1 {
			close(firstStarted)
			<-releaseFirst
		}
		consumed <- message
	})

	messages.Send(1)
	<-firstStarted
	messages.Send(2)
	messages.Send(3)
	close(releaseFirst)

	for want := 1; want <= 3; want++ {
		if got := <-consumed; got != want {
			t.Fatalf("message %d consumed as %d", want, got)
		}
	}
}

func TestMailboxAcceptsConcurrentSenders(t *testing.T) {
	const count = 1000

	consumed := make(chan int, count)
	messages := New(func(message int) {
		consumed <- message
	})

	var senders sync.WaitGroup
	senders.Add(count)
	for message := range count {
		go func() {
			defer senders.Done()
			messages.Send(message)
		}()
	}
	senders.Wait()

	seen := make([]bool, count)
	for range count {
		message := <-consumed
		if seen[message] {
			t.Fatalf("message %d consumed twice", message)
		}
		seen[message] = true
	}
}
