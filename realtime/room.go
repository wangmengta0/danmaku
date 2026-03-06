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
	if _, ok := r.clients[c]; ok {
		delete(r.clients, c)
	}
	close(c.Send)
}
func (r *Room) Size() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}
func (r *Room) Broadcast(data []byte, hub *Hub) {
	r.mu.RLock()
	clients := make([]*Client, 0, len(r.clients))
	for c := range r.clients {
		clients = append(clients, c)
	}
	r.mu.RUnlock()
	for _, client := range clients {
		select {
		case client.Send <- data:
		default:
			hub.Unregister <- client
		}
	}
}
