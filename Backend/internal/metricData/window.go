package metricData 
import (
	"sync"

	"backend/internal/prediction"
)

const windowCapacity = 20
type RingBuffer struct {
	mu       sync.RWMutex
	data     []prediction.MetricPoint
	capacity int
	head     int
	size     int
}


func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data:     make([]prediction.MetricPoint, capacity),
		capacity: capacity,
	}
}

func (r *RingBuffer) Push(point prediction.MetricPoint) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.size < r.capacity {
		index := (r.head + r.size) % r.capacity
		r.data[index] = point
		r.size++
		return
	}

	r.data[r.head] = point
	r.head = (r.head + 1) % r.capacity
}


func (r *RingBuffer) GetWindow() []prediction.MetricPoint {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.size == 0 {
		return nil
	}

	window := make([]prediction.MetricPoint, r.size)
	for i := 0; i < r.size; i++ {
		window[i] = r.data[(r.head+i)%r.capacity]
	}
	return window
}
func (r *RingBuffer) IsFull() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.size == r.capacity
}

type WindowStore struct {
	mu      sync.RWMutex
	buffers map[string]*RingBuffer
}
func NewWindowStore() *WindowStore {
	return &WindowStore{
		buffers: make(map[string]*RingBuffer),
	}
}


func (w *WindowStore) Push(userID string, point prediction.MetricPoint) {
	w.getOrCreateBuffer(userID).Push(point)
}


func (w *WindowStore) GetWindow(userID string) ([]prediction.MetricPoint, bool) {
	w.mu.RLock()
	buffer, exists := w.buffers[userID]
	w.mu.RUnlock()

	if !exists {
		return nil, false
	}
	return buffer.GetWindow(), true
}

func (w *WindowStore) IsReady(userID string) bool {
	w.mu.RLock()
	buffer, exists := w.buffers[userID]
	w.mu.RUnlock()

	return exists && buffer.IsFull()
}

func (w *WindowStore) getOrCreateBuffer(userID string) *RingBuffer {
	w.mu.RLock()
	buffer, exists := w.buffers[userID]
	w.mu.RUnlock()

	if exists {
		return buffer
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if buffer, exists = w.buffers[userID]; exists {
		return buffer
	}

	buffer = NewRingBuffer(windowCapacity)
	w.buffers[userID] = buffer
	return buffer
}