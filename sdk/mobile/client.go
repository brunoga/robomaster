package mobile

import (
	"fmt"
	"net"

	"github.com/brunoga/robomaster/sdk"
)

type Client struct {
	c *sdk.Client
}

func NewClient(ip string) (*Client, error) {
	parsedIp := net.ParseIP(ip)
	if parsedIp == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ip)
	}

	c, err := sdk.NewClient(parsedIp)
	if err != nil {
		return nil, err
	}

	return &Client{c: c}, nil
}

func NewClientUSB() (*Client, error) {
	return NewClient("192.168.42.2")
}

func NewClientWifiDirect() (*Client, error) {
	return NewClient("192.168.2.1")
}

func (c *Client) Open() error {
	return c.c.Open()
}

func (c *Client) Close() error {
	return c.c.Close()
}

func (c *Client) RobotModule() *Robot {
	return &Robot{r: c.c.RobotModule()}
}

func (c *Client) ChassisModule() *Chassis {
	return &Chassis{c: c.c.ChassisModule()}
}
