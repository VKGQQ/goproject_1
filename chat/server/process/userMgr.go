package process

import (
	"fmt"
	"maps"
	"sync"
)

// 因为UserMgr 实例在服务器端有且只有一个
// 因为在很多的地方，都会使用到，因此将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	// onlineUsers 会被多个连接的 goroutine 并发读写，必须用读写锁保护。
	// 读多写少（转发/查询频繁，上下线低频），因此选用 RWMutex。
	mu          sync.RWMutex
	onlineUsers map[int]*UserProcess
}

// 完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers添加（写锁，互斥）
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.onlineUsers[up.UserId] = up
}

// DelOnlineUser 删除离线用户（写锁，互斥）
func (this *UserMgr) DelOnlineUser(userId int) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.onlineUsers, userId)
}

// OnlineUserIds 返回当前所有在线用户 id 的快照切片（读锁）
// 锁内完成拷贝，调用方可在锁外安全遍历
func (this *UserMgr) OnlineUserIds() []int {
	this.mu.RLock()
	defer this.mu.RUnlock()
	ids := make([]int, 0, len(this.onlineUsers))
	for id := range this.onlineUsers {
		ids = append(ids, id)
	}
	return ids
}

// GetAllOnlineUsers 返回在线用户表的快照副本（读锁）
// 必须返回拷贝而不能直接返回原 map：否则调用方在锁外遍历时，
// 其他 goroutine 可能正在并发写 map，依然会 panic。
// 网络发送等耗时操作基于快照在锁外进行，避免持锁做 IO 阻塞上下线。
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	this.mu.RLock()
	defer this.mu.RUnlock()
	snapshot := make(map[int]*UserProcess, len(this.onlineUsers))
	maps.Copy(snapshot, this.onlineUsers)
	return snapshot
}

// GetOnlineUserById 根据id返回对应的在线用户（读锁）
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	this.mu.RLock()
	defer this.mu.RUnlock()
	up, ok := this.onlineUsers[userId]
	if !ok { //说明查找的这个用户，当前不在线。
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
