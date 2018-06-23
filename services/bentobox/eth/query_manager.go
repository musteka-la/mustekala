package eth

import "sync"

// queryManager keeps track of the queries
// being made
type queryManager struct {
	queue queryQueue
}

// queryQueue is the queue element of
// the queryManager
type queryQueue struct {
	sync.RWMutex
	items map[string]struct{}
	count int
}

// newQueryManager initializes the queryManager
func newQueryManager() *queryManager {
	q := &queryManager{
		queue: queryQueue{
			items: make(map[string]struct{}),
			count: 0,
		},
	}

	return q
}

// addQuery keeps track in a map of the queries
// being made
func (q *queryManager) addQuery(id string) {
	q.queue.Lock()
	defer q.queue.Unlock()

	if _, ok := q.queue.items[id]; !ok {
		q.queue.items[id] = struct{}{}
		q.queue.count += 1
	}
}

// removeQuery takes the registered query out
// of the map
func (q *queryManager) removeQuery(id string) {
	q.queue.Lock()
	defer q.queue.Unlock()

	delete(q.queue.items, id)
	q.queue.count -= 1
}

// getQueueCount returns the number of queries
// being made
func (q *queryManager) getQueueCount() int {
	q.queue.Lock()
	defer q.queue.Unlock()

	return q.queue.count
}
