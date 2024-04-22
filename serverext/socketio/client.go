package socketio

import socketio "github.com/googollee/go-socket.io"

type Client struct {
	Conn     socketio.Conn
	UserID   int64
	AppID    int64
	Platform int8
}

func (c *Client) Reset() {
	c.Conn = nil
	c.UserID = 0
	c.AppID = 0
	c.Platform = 0
}

func (c *Client) SetValue(conn socketio.Conn, appID, userID int64, platform int8) {
	c.Conn = conn
	c.UserID = userID
	c.AppID = appID
	c.Platform = platform
}
