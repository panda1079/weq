package config

type Mysql struct {
}

func (r *Mysql) Run() map[string]map[string]string {
	return map[string]map[string]string{
		//数据库的主库配置
		"MYSQL_WRITE": map[string]string{
			"host":     "192.168.2.3",
			"port":     "3306",
			"user":     "test",
			"password": "123456",
			"database": "test",
			"timeout":  "5",
			"charset":  "utf8mb4",
			"MAX_NUM":  "20",
			"TIME_OUT": "3",
		},
		//数据库的从库配置
		"MYSQL_READ": map[string]string{
			"host":     "192.168.2.3",
			"port":     "3306",
			"user":     "test",
			"password": "123456",
			"database": "test",
			"timeout":  "5",
			"charset":  "utf8mb4",
			"MAX_NUM":  "20",
			"TIME_OUT": "3",
		},
	}
}
