package library

import (
	"config"
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

// RedisG 关于连接的公共函数，（由于个人比较懒，就实现了最基本的key增删查）
type RedisG struct {
	Connections    map[string]*redis.Client
	connectionName string //当前连接池名称
}

// InitRedis 加载配置，装载链接
func (r *RedisG) InitRedis() {
	//先把内存分配定下来，不如会出现 panic: assignment to entry in nil map
	r.Connections = make(map[string]*redis.Client)

	//获取启动配置
	deploy := config.Redis{}
	con := deploy.Run()

	//循环加载mysql数据库配置
	for k1, v1 := range con {
		//SetLog(v1, "输出配置详情")
		b := r.Connect(k1, v1)
		if !b {
			SetLog("加载配置"+k1+"失败", "加载配置失败")
			return
		}
	}
}

// Connect 绑定解析，并返回对象
func (r *RedisG) Connect(mName string, con map[string]string) bool {

	db, _ := strconv.Atoi(con["db"]) //转int类型

	rdb := redis.NewClient(&redis.Options{
		Addr:     con["host"] + ":" + con["port"],
		Password: con["auth"], // no password set
		DB:       db,          // use default DB
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		SetLog(err, "连接Redis错误")
		return false
	}

	SetLog(err, "连接Redis-"+con["host"]+":"+con["port"]+"成功~")

	//r.Connections[mName] = new(redis.Client)
	r.Connections[mName] = rdb

	return true
}

// Connection 连接名设置
func (r *RedisG) Connection(name string) *RedisG {
	r.connectionName = name
	return r
}

// Load 如果需要更多功能，就直接用这个继承拉起
func (r *RedisG) Load(name string) *redis.Client {
	return r.Connections[name]
}

// GetConnName 获取connectionName
func (r *RedisG) GetConnName() string {
	if Empty(r.connectionName) {
		return "write"
	}
	return r.connectionName
}

// Set 写入数据
func (r *RedisG) Set(key string, value interface{}, expiration time.Duration) bool {
	err := r.Connections[r.GetConnName()].Set(key, value, time.Second*expiration).Err()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-写入数据错误")
		return false
	}
	return true
}

// Exists 检查key是否存在
func (r *RedisG) Exists(key string) bool {
	val, err := r.Connections[r.GetConnName()].Exists(key).Result()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-检查数据错误")
		return false
	}

	//EXISTS 命令可以检查指定键是否存在，并返回一个整数值，表示键存在时返回1，不存在时返回0
	return val > 0
}

// Get 读取数据
func (r *RedisG) Get(key string) string {
	//先检查一下key是否存在
	if !r.Exists(key) {
		return ""
	}

	val, err := r.Connections[r.GetConnName()].Get(key).Result()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-读取数据错误")
		return ""
	}
	return val
}

// Del 删除数据
func (r *RedisG) Del(key string) bool {
	//先检查一下key是否存在
	if !r.Exists(key) {
		return true
	}

	err := r.Connections[r.GetConnName()].Del(key).Err()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-删除数据错误")
		return false
	}
	return true
}

// LLen 返回列表的长度。 如果列表 key 不存在，则 key 被解释为一个空列表，返回 0。 如果 key 不是列表类型，返回一个错误。
func (r *RedisG) LLen(key string) int64 {
	//先检查一下key是否存在
	if !r.Exists(key) {
		return 0
	}

	val, err := r.Connections[r.GetConnName()].LLen(key).Result()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-读取列表长度错误")
		return 0
	}
	return val
}

// LPush 从列表左边插入数据
func (r *RedisG) LPush(key string, value ...interface{}) bool {
	err := r.Connections[r.GetConnName()].LPush(key, value...).Err()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-从列表左边插入数据错误")
		return false
	}
	return true
}

// RPop 从列表的右边删除第一个数据，并返回删除的数据
func (r *RedisG) RPop(key string) string {
	//先检查一下key是否存在
	if !r.Exists(key) {
		return ""
	}

	val, err := r.Connections[r.GetConnName()].RPop(key).Result()
	if err != nil {
		SetLog(err, "Redis-"+r.GetConnName()+"-读取数据错误")
		return ""
	}
	return val
}
