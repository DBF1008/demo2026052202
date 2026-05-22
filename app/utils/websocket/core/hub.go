package core

type Hub struct {

	Register chan *Client

	UnRegister chan *Client

	Clients map[*Client]bool
}

func CreateHubFactory() *Hub {
	return &Hub{
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
		case client := <-h.UnRegister:
			if _, ok := h.Clients[client]; ok {
				_ = client.Conn.Close()
				delete(h.Clients, client)
			}
		}
	}
}
