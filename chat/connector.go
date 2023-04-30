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

func (c *connector) MakeBCast(group_name string, text string) error {
	bcast_request := BCastRequest{
		GroupName: group_name,
		Data:      []byte(text),
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
