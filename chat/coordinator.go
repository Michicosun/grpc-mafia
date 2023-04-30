package chat

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
)

const (
	MAX_PACKET_SIZE = 1024
	QUEUE_SIZE      = 128
)

type BCastRequest struct {
	GroupName string `json:"group_name"`
	Data      []byte `json:"data"`
}

type CreateGroupRequest struct {
	GroupName string   `json:"group_name"`
	Addrs     []string `json:"addrs"`
}

type RequestType uint32

const (
	BCAST_TYPE        = 0
	CREATE_GROUP_TYPE = 1
)

type Request struct {
	ReqType RequestType `json:"req_type"`
	Data    []byte      `json:"data"`
}

type Coordinator struct {
	socket      net.PacketConn
	groups      map[string][]net.Addr
	bcast_queue chan BCastRequest
}

func (c *Coordinator) Listen() {
	buf := make([]byte, 1024)
	for {
		n, addr, err := c.socket.ReadFrom(buf)
		if err != nil {
			zlog.Error().Err(err).Msg("read from socket")
			continue
		}

		zlog.Info().Str("buf", string(buf[:n])).Msgf("read from socket from %s", addr)

		if err := c.processRequest(buf[:n]); err != nil {
			zlog.Error().Err(err).Msg("process request data")
		}
	}
}

func (c *Coordinator) processRequest(buf []byte) error {
	r := Request{}
	if err := json.Unmarshal(buf, &r); err != nil {
		return errors.Wrap(err, "parse request")
	}

	switch r.ReqType {
	case BCAST_TYPE:
		bcast_req := BCastRequest{}
		if err := json.Unmarshal(r.Data, &bcast_req); err != nil {
			return errors.Wrap(err, "parse bcast request")
		}

		zlog.Info().Str("group", bcast_req.GroupName).Str("text", string(bcast_req.Data)).Msg("get bcast request")

		c.bcast_queue <- bcast_req

	case CREATE_GROUP_TYPE:
		create_req := CreateGroupRequest{}
		if err := json.Unmarshal(r.Data, &create_req); err != nil {
			return errors.Wrap(err, "parse create group request")
		}

		zlog.Info().Str("group", create_req.GroupName).Msg("get create group request")

		c.createGroup(&create_req)
	}

	return nil
}

func (c *Coordinator) createGroup(r *CreateGroupRequest) {
	addrs := make([]net.Addr, 0)

	for _, addr := range r.Addrs {
		net_addr, err := net.ResolveUDPAddr("udp4", addr)
		if err != nil {
			zlog.Error().Err(err).Str("address", addr).Msg("cannot resolve address")
			continue
		}

		zlog.Info().Str("address", addr).Str("to", net_addr.String()).Msg("successfully resolved")

		addrs = append(addrs, net_addr)
	}

	c.groups[r.GroupName] = addrs
}

func (c *Coordinator) WorkerRoutine() {
	for r := range c.bcast_queue {
		c.sendBCast(&r)
	}
}

func (c *Coordinator) sendBCast(r *BCastRequest) {
	group, ok := c.groups[r.GroupName]
	if !ok {
		zlog.Error().Str("name", r.GroupName).Msg("undefined group")
		return
	}

	for i, addr := range group {
		zlog.Info().Str("progress", fmt.Sprintf("%d/%d", i, len(group))).Msgf("sending message to %s", addr)
		c.socket.WriteTo(r.Data, addr)
	}
}

func (c *Coordinator) CreateGroup(group_name string, addrs []net.Addr) {
	c.groups[group_name] = addrs
}

func MakeCoordinator() (*Coordinator, error) {
	port := os.Getenv("PORT")

	zlog.Info().Str("port", port).Msg("configure udp socket")

	socket, err := net.ListenPacket("udp", fmt.Sprintf(":%s", port))
	if err != nil {
		return nil, err
	}

	coordinator := &Coordinator{
		socket:      socket,
		groups:      make(map[string][]net.Addr),
		bcast_queue: make(chan BCastRequest, QUEUE_SIZE),
	}

	for i := 0; i < QUEUE_SIZE; i += 1 {
		go coordinator.WorkerRoutine()
	}

	return coordinator, nil
}
