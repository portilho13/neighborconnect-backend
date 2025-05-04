package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	controllers_models "github.com/portilho13/neighborconnect-backend/models"
	repositoryControllers "github.com/portilho13/neighborconnect-backend/repository/controlers/marketplace"
)

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	listingID string
}

type HubS struct {
	Rooms      map[string]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
}

type BroadcastMessage struct {
	ListingID string
	Message   []byte
}

var Hub = HubS{
	Rooms:      make(map[string]map[*Client]bool),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan BroadcastMessage),
}

func (h *HubS) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.listingID]; !ok {
				h.Rooms[client.listingID] = make(map[*Client]bool)
			}
			h.Rooms[client.listingID][client] = true

		case client := <-h.Unregister:
			if clients, ok := h.Rooms[client.listingID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
			}

		case b := <-h.Broadcast:
			if clients, ok := h.Rooms[b.ListingID]; ok {
				for client := range clients {
					select {
					case client.send <- b.Message:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(w http.ResponseWriter, r *http.Request, dbPool *pgxpool.Pool) {
	vars := r.URL.Query()
	listingID := vars.Get("listing_id")
	if listingID == "" {
		http.Error(w, "missing listing_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{
		conn:      conn,
		send:      make(chan []byte, 256),
		listingID: listingID,
	}
	Hub.Register <- client

	listingIdInt, err := strconv.Atoi(listingID)
	if err != nil {
		log.Fatal(err)
	}

	bids, err := repositoryControllers.GetBidByListningId(listingIdInt, dbPool)
	if err != nil {
		log.Fatal(err)
	}

	listing, err := repositoryControllers.GetListingById(listingIdInt, dbPool)
	if err != nil {
		log.Fatal(err)
	}

	var bidJson controllers_models.BidInfo
	var bidJsonMarshed []byte

	if len(bids) == 0 {
		bidJson.Id = nil
		bidJson.Bid_Ammount = listing.Start_Price
		bidJson.Bid_Time = nil
		bidJson.User_Id = nil
		bidJson.Listing_Id = listingIdInt

		bidJsonMarshed, err = json.Marshal(bidJson)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		highestBid := bids[0]

		bidJson.Id = highestBid.Id
		bidJson.Bid_Ammount = highestBid.Bid_Ammount
		bidJson.Bid_Time = &highestBid.Bid_Time
		bidJson.User_Id = highestBid.User_Id
		bidJson.Listing_Id = *highestBid.Listing_Id

		bidJsonMarshed, err = json.Marshal(bidJson)
		if err != nil {
			log.Fatal(err)
		}
	}

	Hub.Broadcast <- BroadcastMessage{
		ListingID: listingID,
		Message:   bidJsonMarshed,
	}

	go client.writePump()
}

func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
