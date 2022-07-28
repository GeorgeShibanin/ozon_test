package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const namespace = "rl"

//Принимает клиент, название действия, временное окно в рамках которого ограничивается количесто действий и максимальное количество действий в рамках окна
func NewLimiter(client *redis.Client, action string, period time.Duration, limit int64) *Limiter {

	return &Limiter{client, action, period, limit}
}

type Limiter struct {
	client *redis.Client

	action string
	period time.Duration
	limit  int64
}

//go:embed incr_expirenx.lua
var incrExpireLua string
var incrExpireScript = redis.NewScript(incrExpireLua)

//Возвращает True или False в зависимости от того можно или нельзя совершать действие
func (l *Limiter) CanDoAt(ctx context.Context, ts time.Time) (bool, error) {
	key := l.key(ts)
	ttlMs := l.period.Milliseconds()

	//lua скрипт проверяет наличие ключа и в зависимости есть ключ или нет делает команду PEXP
	rawCount, err := incrExpireScript.Run(ctx, l.client, []string{key}, ttlMs).Result()
	if err != nil {
		return false, err
	}
	count := rawCount.(int64)

	return count <= l.limit, nil
}

func (l *Limiter) key(ts time.Time) string {
	//количество наносекунд с начада эпохи и делим на количество секунд в нашем периоде - получим номер периода от начала Unix эпохи
	interval := ts.UTC().UnixNano() / l.period.Nanoseconds()
	//ко всем ключам ratelimit добавляем префик "rl", чтобы они не перепутались с сокращёнными ссылками с префиксом "surl"
	//в ключ кодируем название действия и интервал
	return fmt.Sprintf("%s:%s:%x", namespace, l.action, interval)
}
