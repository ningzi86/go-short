package storage

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"time"
	"fmt"
	"github.com/mattheath/base62"
	"crypto/sha256"
	"encoding/hex"
	"go-short/serror"
	"errors"
)

type Storage interface {
	Shorten(url string, exp int64) (string, error)
	ShortenInfo(eid string) (interface{}, error)
	Unshorten(eid string) (string, error)
}

const (
	URLIDKEY           = "next.url.id"
	ShortLinkKey       = "shortLink:%s:url"
	URLHashKey         = "urlhash:%s:url"
	ShortlinkDetailKey = "shortlink:%s:detail"
)

type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL      string        `json:"url"`
	CreateAt string        `json:"create_at"`
	Expired  time.Duration `json:"expired"`
}

func NewRedisClient(addr string, pass string, db int) *RedisCli {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}
	return &RedisCli{Cli: c}
}

func toSha1(url string) string {
	//使用sha256哈希函数
	h := sha256.New()
	h.Write([]byte(url))
	sum := h.Sum(nil)

	//由于是十六进制表示，因此需要转换
	s := hex.EncodeToString(sum)
	return string(s)
}

func (r *RedisCli) Shorten(url string, exp int64) (string, error) {

	h := toSha1(url)
	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()

	if err != nil {
		if err != redis.Nil {
			return "", err
		}
	}
	if d != "" {
		return d, nil
	}

	err = r.Cli.Incr(URLIDKEY).Err()
	if err != nil {
		return "", err
	}

	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}

	eid := base62.EncodeInt64(id)
	err = r.Cli.Set(fmt.Sprintf(ShortLinkKey, eid), url, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(URLHashKey, h), eid, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(&URLDetail{
		URL:      url,
		CreateAt: time.Now().String(),
		Expired:  time.Duration(exp),
	})
	if err != nil {
		return "", err
	}

	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, eid), detail, time.Minute*time.Duration(exp)).Err()
	if err != nil {
		return "", err
	}
	return eid, nil

}

func (r *RedisCli) ShortenInfo(eid string) (interface{}, error) {

	key  := fmt.Sprintf(ShortlinkDetailKey, eid)
	d, err := r.Cli.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", serror.StatusError{404, errors.New("Unknown short URL")}
		}
		return "", err
	}
	return d, nil

}

func (r *RedisCli) Unshorten(eid string) (string, error) {

	url, err := r.Cli.Get(fmt.Sprintf(ShortLinkKey, eid)).Result()
	if err != nil {
		return "", err
	}
	return url, nil

}

