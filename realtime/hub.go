package realtime

type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			room := h.getOrCreateRoom(client.RoomId)
			room.AddClient(client)
		case client := <-h.Unregister:
			if room, ok := h.Rooms[client.RoomId]; ok {
				room.RemoveClient(client)
			}
		}
	}
}
func (h *Hub) getOrCreateRoom(roomId string) *Room {
	if room, ok := h.Rooms[roomId]; ok {
		return room
	}
	room := NewRoom(roomId)
	h.Rooms[roomId] = room
	return room
}
