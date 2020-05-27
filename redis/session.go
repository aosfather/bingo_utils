package redis

import "github.com/go-redis/redis"

//分布式session存储实现
type RedisSessionStore struct {
	Addr   string //地址
	Db     int    //数据库
	Pwd    string //密码
	client *redis.Client
	Expire int64 //超期时间
}

func (this *RedisSessionStore) Init() {
	//防止多次初始化
	if this.client != nil {
		return
	}

	if this.Db < 0 || this.Db > 16 {
		this.Db = 0
	}

	this.client = redis.NewClient(&redis.Options{
		Addr:     this.Addr,
		Password: this.Pwd, // no password set
		DB:       this.Db,  // use default DB
	})
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
	if this.Expire > 0 {
		this.client.Expire(id, ToSecond(this.Expire))
	}

}
func (this *RedisSessionStore) GetValue(id, key string) interface{} {
	if id != "" && key != "" {
		result := this.client.HGet(id, key)
		if result != nil {
			return result.String()
		}
	}
	return nil
}
func (this *RedisSessionStore) SetValue(id, key string, value interface{}) {
	if id != "" && key != "" {
		this.client.HSet(id, key, value.(string))
	}
}
func (this *RedisSessionStore) Touch(id string) {
	if this.Expire > 0 {
		this.client.Expire(id, ToSecond(this.Expire))
	}
}

func (this *RedisSessionStore) Delete(id string) {
	this.client.Del(id)
}
