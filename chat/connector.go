package chat

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
)

var Connector = &connector{}

type connector struct {
	socket     net.PacketConn
	coord_addr net.Addr
}

type Message struct {
	From string `json:"from"`
	Text string `json:"text"`
}

func (c *connector) writeToCoordinator(r *Request) error {
	data, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "create abstract request")
	}

	if _, err := c.socket.WriteTo(data, c.coord_addr); err != nil {
		return errors.Wrap(err, "write to coordinator")
	}

	return nil
}

func (c *connector) CreateGroup(group_name string, addrs []string) error {
	create_group_request := CreateGroupRequest{
		GroupName: group_name,
		Addrs:     addrs,
	}

	request_data, err := json.Marshal(create_group_request)
	if err != nil {
		return errors.Wrap(err, "create specific request")
	}

	request := Request{
		ReqType: CREATE_GROUP_TYPE,
		Data:    request_data,
	}

	return c.writeToCoordinator(&request)
}

func (c *connector) MakeBCast(group_name string, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "serialize message")
	}

	bcast_request := BCastRequest{
		GroupName: group_name,
		Data:      data,
	}

	request_data, err := json.Marshal(bcast_request)
	if err != nil {
		return errors.Wrap(err, "create specific request")
	}

	request := Request{
		ReqType: BCAST_TYPE,
		Data:    request_data,
	}

	return c.writeToCoordinator(&request)
}

func (c *connector) RecvMessage() (*Message, error) {
	buf := make([]byte, MAX_PACKET_SIZE)

	n, _, err := c.socket.ReadFrom(buf)
	if err != nil {
		return nil, errors.Wrap(err, "read from chat connector")
	}

	msg := Message{}
	if err := json.Unmarshal(buf[:n], &msg); err != nil {
		return nil, errors.Wrap(err, "deserialize coordinator response")
	}

	return &msg, nil
}

func (c *connector) Init(port_to_listen uint32, coord_hostname string, coord_port uint32) error {
	coord_conn := fmt.Sprintf("%s:%d", coord_hostname, coord_port)

	addr, err := net.ResolveUDPAddr("udp4", coord_conn)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("resolve addr: %s", coord_conn))
	}

	socket, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port_to_listen))
	if err != nil {
		return errors.Wrap(err, "create udp socket")
	}

	Connector.coord_addr = addr
	Connector.socket = socket

	return nil
}
