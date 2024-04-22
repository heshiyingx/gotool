package socketio

import (
	"slices"
	"strconv"
	"sync"
)

type clientMaps struct {
	lockMap sync.Map
	m       map[string][]*Client
}

func NewClientMaps() *clientMaps {
	return &clientMaps{
		m:       make(map[string][]*Client),
		lockMap: sync.Map{},
	}

}

func (c *Client) GetKey() string {
	return getClientKey(c.AppID, c.UserID)

}
func getClientKey(appID, userID int64) string {
	return strconv.FormatInt(appID, 10) + "_" + strconv.FormatInt(userID, 10)

}

// AddClient 添加客户端，返回需要踢下线的客户端
func (m *clientMaps) AddClient(c *Client) *Client {
	locker := m.getClientLock(c)
	locker.Lock()
	defer locker.Unlock()
	var kickOffLineCLient *Client
	clients, ok := m.m[c.GetKey()]
	if !ok {
		clients = make([]*Client, 0, 1)
		clients = append(clients, c)
	} else {
		slices.DeleteFunc(clients, func(c *Client) bool {
			if c.Platform == kickOffLineCLient.Platform {
				kickOffLineCLient = c
				return true
			}
			return false
		})

	}
	return kickOffLineCLient
}
func (m *clientMaps) getClientLock(c *Client) *sync.RWMutex {
	return m.getClientLockByKey(c.GetKey())
}
func (m *clientMaps) getClientLockByKey(key string) *sync.RWMutex {
	rwLocker, ok := m.lockMap.Load(key)
	if !ok {
		locker := &sync.RWMutex{}
		rwLocker, _ = m.lockMap.LoadOrStore(key, locker)
	}
	return rwLocker.(*sync.RWMutex)
}
func (m *clientMaps) RemoveClient(c *Client) {
	locker := m.getClientLock(c)
	locker.Lock()
	defer locker.Unlock()
	clients, ok := m.m[c.GetKey()]
	if !ok {
		return
	}
	slices.DeleteFunc(clients, func(c *Client) bool {
		return c.Platform == c.Platform
	})
}
func (m *clientMaps) GetUserClients(appid int64, userID int64) []*Client {
	locker := m.getClientLockByKey(getClientKey(appid, userID))
	locker.RLock()
	defer locker.Unlock()
	clients := m.m[getClientKey(appid, userID)]
	return clients
}
