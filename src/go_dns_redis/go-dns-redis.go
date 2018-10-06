package godnsredis

import (
	"fmt"

	"github.com/go-redis/redis"
)

// CreateRedisDatabaseConnection : Create a New Client for Redis Server...
func CreateRedisDatabaseConnection() *redis.Client {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println("Redis Connection Ping Check : ", pong, err)
	// Output: PONG <nil>

	return redisClient
}

// SetRedisKey : Setting a new Key-Value pair in Redis
func SetRedisKey(redisClient *redis.Client, redisKey string, redisValue string) string {

	err := redisClient.Set(redisKey, redisValue, 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := redisClient.Get(redisKey).Result()
	if err != nil {
		return "SetRedisKeyError - " + err.Error()
	}

	fmt.Println("SetRedisKey returns ", redisKey, val)
	return redisKey

}

// GetRedisKey : Get the Value from Redis while passing an Redis key
func GetRedisKey(redisClient *redis.Client, redisKey string) string {

	fmt.Println("Redis Key Received - " + redisKey)

	val, err := redisClient.Get(redisKey).Result()

	if err != nil {
		panic(err)
	} else {
		fmt.Println("GetRedisKey returns ", redisKey, val)
		return val
	}

}

// IsRedisKey : Checks if Redis has a key with the given DNS Alias
func IsRedisKey(redisClient *redis.Client, dnsAlias string) bool {

	fmt.Println("Redis DNS Alias key received - \\" + dnsAlias + "\\")

	val, err := redisClient.Keys(dnsAlias).Result()

	if err != nil {
		panic(err)
	} else {
		if len(val) == 0 {
			fmt.Println("IsRedisKey returns ", dnsAlias, val)
			return false
		}
		fmt.Println("IsRedisKey returns ", dnsAlias, val)
		return true
	}

}
