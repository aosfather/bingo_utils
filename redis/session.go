package redis

import (
	"github.com/go-redis/redis"
	"strings"
)

//分布式session存储实现
type RedisSessionStore struct {
	client redis.Cmdable
	expire int64 //超期时间

}

func (this *RedisSessionStore) InitByCluster(addr []string, pwd string, expire int64) {
	//防止多次初始化
	if this.client != nil {
		return
	}

	if expire > 0 {
		this.expire = expire
	}

	this.client = redis.NewClusterClient(&redis.ClusterOptions{Addrs: addr, Password: pwd})
}
func (this *RedisSessionStore) Init(addr string, db int, pwd string, expire int64) {
	//防止多次初始化
	if this.client != nil {
		return
	}

	if expire > 0 {
		this.expire = expire
	}
	if addr == "" {
		return
	}

	if db < 0 || db > 16 {
		db = 0
	}
	var option redis.Options
	//if strings.TrimSpace(pwd)=="" {
	//	option=redis.Options{
	//		Addr:     addr,
	//		DB:       db,  // use default DB
	//	}
	//}else {
	option = redis.Options{
		Addr:     addr,
		Password: strings.TrimSpace(pwd), // no password set
		DB:       db,                     // use default DB
	}
	//}
	this.client = redis.NewClient(&option)
}

func (this *RedisSessionStore) Exist(id string) bool {
	result := this.client.Exists(id)
	if result != nil {
		count, _ := result.Result()
		if count >= 1 {
			return true
		}
	}
	return false
}

func (this *RedisSessionStore) Create(id string) {
	this.client.HSetNX(id, "_id_", id)
	if this.expire > 0 {
		this.client.Expire(id, ToSecond(this.expire))
	}

}
func (this *RedisSessionStore) GetValue(id, key string) interface{} {
	if id != "" && key != "" {
		result := this.client.HGet(id, key)
		if result != nil {
			v, err := result.Result()
			if err != nil {
				return nil
			}
			return v
		}
	}
	return nil
}
func (this *RedisSessionStore) SetValue(id, key string, value interface{}) {
	if id != "" && key != "" {
		if value != nil {
			this.client.HSet(id, key, value.(string))
		} else {
			this.client.HDel(id, key)
		}

	}
}
func (this *RedisSessionStore) Touch(id string) {
	if this.expire > 0 {
		this.client.Expire(id, ToSecond(this.expire))
	}
}

func (this *RedisSessionStore) Delete(id string) {
	this.client.Del(id)
}
