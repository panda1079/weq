package config

type Redis struct {
}

func (r *Redis) Run() map[string]map[string]string {
	return map[string]map[string]string{
		//数据库的主库配置
		"write": map[string]string{
			"host": "192.168.2.3", // Redis 地址
			"port": "6379",        // Redis 端口
			"auth": "123456",      // Redis 密码
			"db":   "1",           // redis数据库
		},
	}
}
