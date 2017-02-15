package main

import (
	"github.com/daiLlew/golang-exercises/trace"
	"log"
	"net/http"
)

type room struct {
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
	tracer  trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			// joining.
			r.clients[client] = true
			r.tracer.Trace("New client joined.")
		case client := <-r.leave:
			// leaving
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left.")
		case msg := <-r.forward:
			// forward
			r.tracer.Trace("Message recieved: " + string(msg))
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("---- Message sent to client.")
			}

		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}
