package uds

import (
	"net"
	"strings"

	"media-box-ui/pkg/logger"
)

type ClientConfig struct {
	SocketPath string
}

type Client struct {
	cfg  *ClientConfig
	log  *logger.Zerolog
	conn net.Conn
}

func NewUDSClient(cfg *ClientConfig, log *logger.Zerolog) *Client {
	s := &Client{
		cfg: cfg,
		log: log,
	}
	return s
}

func (c *Client) connect() error {
	var err error
	c.conn, err = net.Dial("unix", c.cfg.SocketPath)
	if err != nil {
		//c.log.Error().Msgf("connect error: %v", err)
		c.conn = nil
		return err
	}

	go c.reader()

	c.log.Debug().Msg("client connected")

	return nil
}

func (c *Client) Write(msg []byte) error {
	if c.conn == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}

	_, err := c.conn.Write(msg)
	if err != nil {
		c.conn = nil
		c.log.Error().Msgf("write error: %v", err)
		return err
	}

	return nil
}

func (c *Client) reader() {
	buf := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buf[:])
		if err != nil {
			c.log.Error().Msgf("client read error: %v", err)
			return
		}

		msg := strings.TrimSpace(string(buf[:n]))
		mpvUseCase.MPVResponse([]byte(msg))
	}
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}
