package realtime

import "sync"

type Room struct {
	Id      string
	clients map[*Client]struct{}
	mu      sync.RWMutex
}

func NewRoom(id string) *Room {
	return &Room{
		Id:      id,
		clients: make(map[*Client]struct{}),
	}
}
func (r *Room) AddClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}
func (r *Room) RemoveClient(c *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
	close(c.Send)
}
func (r *Room) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.clients)
}
