package app

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"encoding/json"
	"errors"
	"github.com/mattheath/base62"
	"github.com/speps/go-hashids"
)

const(
	//全局自增器
	URLIdKey = "next.url.id"
	//短地址和地址的映射
	ShortLinkKey = "shortlink:%s:url"
	//地址hash和短地址的映射
	URLHashKey = "urlhash:%s:url"
	//短地址和详情的映射
	ShortLinkDetailKey = "shortlink:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL 		string	`json:"url"`
	CreateAt 	string 	`json:"create_at"`
	ExpirationInMinutes time.Duration 	`json:"expiration_in_minutes"`
}

func NewRedisCli(addr string, password string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr: addr,
		Password: password,
		DB: db,
	})
	if _,err := c.Ping().Result();err != nil{
		panic(err)
	}
	return &RedisCli{Cli:c}
}

func (r *RedisCli) Shorten(url string, exp int64) (string, error) {
	h := toSha1(url)

	d,err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()
	if err == redis.Nil{

	}else if err != nil{
		return "",err
	}else{
		if d == "{}" {

		}else{
			return d,nil
		}
	}

	err = r.Cli.Incr(URLIdKey).Err()
	if err != nil {
		return "",err
	}
	if err != nil{
		return "",err
	}

	id,err := r.Cli.Get(URLIdKey).Int64()
	if err != nil {
		return "",err
	}
	eid := base62.EncodeInt64(id)

	if err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, eid), url,
		time.Minute * time.Duration(exp)).Err();err != nil {
		return "",err
	}

	if err = r.Cli.Set(fmt.Sprintf(URLHashKey, h), eid,
		time.Minute * time.Duration(exp)).Err();err != nil {
		return "",err
	}

	detail,err := json.Marshal(&URLDetail{
		URL: url,
		CreateAt: time.Now().String(),
		ExpirationInMinutes: time.Duration(exp),
	})
	if err != nil{
		return "",err
	}

	err = r.Cli.Set(fmt.Sprintf(ShortLinkDetailKey, eid), detail,
		time.Minute * time.Duration(exp)).Err()
	if err != nil{
		return "",err
	}

	return eid,nil
}

func (r *RedisCli)ShortLinkInfo(eid string)(interface{}, error){
	d,err := r.Cli.Get(fmt.Sprintf(ShortLinkDetailKey,eid)).Result()
	if err == redis.Nil{
		return "",StatusError{404,errors.New("Unknow short URL")}
	}else if err != nil {
		return "",err
	}else{
		var res URLDetail
		if err = json.Unmarshal([]byte(d), &res);err != nil {
			return "",err
		}
		return res,nil
	}
}

func (r *RedisCli)UnShorten(eid string)(string, error)  {
	url,err := r.Cli.Get(fmt.Sprintf(ShortLinkKey,eid)).Result()
	if err == redis.Nil{
		return "",StatusError{404,err}
	}else if err != nil {
		return "",err
	}else{
		return url,nil
	}
}

func toSha1(url string) string {
	hd := hashids.NewData()
	hd.Salt = url
	hd.MinLength = 0
	h, _ := hashids.NewWithData(hd)
	r, _ := h.Encode([]int{45,434,1313,99})
	return r
}

