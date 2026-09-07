package domain

import "fmt"

const idempotencyKeyPrefix = "payment:idempotency:"

// BuildIdempotencyKey - собирает реальный Redis-ключ из сырого значения,
// присланного клиентом (например, из заголовка Idempotency-Key).
// Один источник форматирования - чтобы префикс не разъехался между вызовами.
func BuildIdempotencyKey(rawKey string) string {
	return fmt.Sprintf("%s%s", idempotencyKeyPrefix, rawKey)
}
