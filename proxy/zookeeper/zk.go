package zookeeper

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// ZkManager ZkManager
type ZkManager struct {
	hosts      []string
	conn       *zk.Conn
	pathPrefix string
}

// NewZkManager 创建NewZkManager
func NewZkManager(hosts []string) *ZkManager {
	return &ZkManager{hosts: hosts, pathPrefix: "/gateway_servers_"}
}

// GetConnect 连接 zk 服务器
func (z *ZkManager) GetConnect() error {
	conn, _, err := zk.Connect(z.hosts, 5*time.Second)
	if err != nil {
		return err
	}
	z.conn = conn
	return nil
}

// Close 关闭服务
func (z *ZkManager) Close() {
	z.conn.Close()
	return
}

// GetPathData 获取配置
func (z *ZkManager) GetPathData(nodePath string) ([]byte, *zk.Stat, error) {
	return z.conn.Get(nodePath)
}

// SetPathData 更新配置
func (z *ZkManager) SetPathData(nodePath string, config []byte, version int32) (err error) {
	ex, _, _ := z.conn.Exists(nodePath)
	if !ex {
		z.conn.Create(nodePath, config, 0, zk.WorldACL(zk.PermAll))
		return nil
	}
	_, dStat, err := z.GetPathData(nodePath)
	if err != nil {
		return
	}
	_, err = z.conn.Set(nodePath, config, dStat.Version)
	if err != nil {
		fmt.Println("Update node error", err)
		return err
	}
	fmt.Println("SetData ok")
	return
}

// RegistServerPath 创建临时节点
// nodePath 主节点 host 要注册的子节点
func (z *ZkManager) RegistServerPath(nodePath, host string) (err error) {
	// 父节点不存在就先添加父节点，并且设置成持久化节点
	ex, _, err := z.conn.Exists(nodePath)
	if err != nil {
		fmt.Println("Exists error", nodePath)
		return err
	}
	if !ex {
		//持久化节点，思考题：如果不是持久化节点会怎么样？
		// 因为主节点是临时的，会导致会话结束时，所有子节点都丢失
		// 持久化保证能够获取到下游的临时节点
		_, err = z.conn.Create(nodePath, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("Create error", nodePath)
			return err
		}
	}
	//临时节点
	// 生成子节点的 path
	subNodePath := nodePath + "/" + host
	// 查寻是否已存在
	ex, _, err = z.conn.Exists(subNodePath)
	if err != nil {
		fmt.Println("Exists error", subNodePath)
		return err
	}
	// 如果不存在就创建子节点（这里设置成临时节点）
	if !ex {
		_, err = z.conn.Create(subNodePath, nil, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("Create error", subNodePath)
			return err
		}
	}
	return
}

// GetServerListByPath 获取服务列表
func (z *ZkManager) GetServerListByPath(path string) (list []string, err error) {
	list, _, err = z.conn.Children(path)
	return
}

// WatchServerListByPath watch机制，服务器有断开或者重连，收到消息
func (z *ZkManager) WatchServerListByPath(path string) (chan []string, chan error) {
	conn := z.conn
	// 快照 channel ，只要节点有变化就会往这个 channel 中输送内容
	snapshots := make(chan []string)
	// error channle ，出现错误是就往该 channel 中写数据
	errors := make(chan error)
	// 开启协程，循环监听 path 地址的更新
	go func() {
		for {
			// 这里监听子节点的变化
			// 快照接收的是子节点列表内容
			// events 是监听的子节点发生变化的主节点的地址（例如监听到某个子节点移除了，会返回它的主节点地址，
			// 该子节点触发的消息类型－移除，也会返回到 events 中）
			// err 接收监听不成功返回的错误
			snapshot, _, events, err := conn.ChildrenW(path)
			if err != nil {
				errors <- err
			}
			snapshots <- snapshot
			select {
			// 当 events 中有数据时，写入到 evt 里面
			case evt := <-events:
				// 如果这个数据是错误的话，写到错误 channel 里
				if evt.Err != nil {
					errors <- evt.Err
				}
				// 打印发生变更的子节点的地址和它所发生的变更事件（移除，修改...）
				fmt.Printf("ChildrenW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()

	return snapshots, errors
}

// WatchPathData watch机制，监听节点值变化
func (z *ZkManager) WatchPathData(nodePath string) (chan []byte, chan error) {
	conn := z.conn
	snapshots := make(chan []byte)
	errors := make(chan error)

	go func() {
		for {
			dataBuf, _, events, err := conn.GetW(nodePath)
			if err != nil {
				errors <- err
				return
			}
			snapshots <- dataBuf
			select {
			case evt := <-events:
				if evt.Err != nil {
					errors <- evt.Err
					return
				}
				fmt.Printf("GetW Event Path:%v, Type:%v\n", evt.Path, evt.Type)
			}
		}
	}()
	return snapshots, errors
}
