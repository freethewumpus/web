package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
	"strings"
)

func GetUser(auth string) (string, string, string, bool) {
	RedisUserFmt := fmt.Sprintf("u:%s", auth)
	cache, err := RedisConnection.Get(RedisUserFmt).Result()
	if err == redis.Nil {
		cursor, err := r.Table("tokens").Get(auth).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
		if cursor.IsNil() {
			return "", "", "", false
		}
		var token Token
		err = cursor.One(&token)
		if err != nil {
			panic(err)
		}
		err = cursor.Close()
		if err != nil {
			panic(err)
		}

		cursor, err = r.Table("users").Get(token.Uid).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
		if cursor.IsNil() {
			return "", "", "", false
		}
		var info User
		err = cursor.One(&info)
		if err != nil {
			panic(err)
		}
		err = cursor.Close()
		if err != nil {
			panic(err)
		}
		RedisConnection.Set(RedisUserFmt, fmt.Sprintf("%s|%s|%s", info.Domain, token.Uid, info.NamingScheme), 0)
		return info.Domain, token.Uid, info.NamingScheme, true
	} else if err == nil {
		parts := strings.Split(cache, "|")
		domain := parts[0]
		uid := parts[1]
		NamingScheme := parts[2]
		return domain, uid, NamingScheme, true
	} else {
		panic(err)
	}
}

func GetDomain(DomainId string) (Domain, bool) {
	RedisDomainFmt := fmt.Sprintf("d:%s", DomainId)
	cache, err := RedisConnection.Get(RedisDomainFmt).Result()
	var domain Domain
	if err == redis.Nil {
		cursor, err := r.Table("domains").Get(DomainId).Run(RethinkConnection)
		if err != nil {
			panic(err)
		}
		if cursor.IsNil() {
			return domain, false
		}
		err = cursor.One(&domain)
		if err != nil {
			panic(err)
		}
		err = cursor.Close()
		if err != nil {
			panic(err)
		}
		res, err := json.Marshal(domain)
		if err != nil {
			panic(err)
		}
		RedisConnection.Set(RedisDomainFmt, string(res), 0)
		return domain, true
	} else if err == nil {
		err := json.Unmarshal([]byte(cache), &domain)
		if err != nil {
			panic(err)
		}
		return domain, true
	} else {
		panic(err)
	}
}
