## TTLCache - an in-memory LRU cache with expiration

TTLCache is a minimal wrapper over a map with no values (simple boolean lookup table) in golang, entries of which are

1. Thread-safe
2. Auto-Expiring after a certain time
3. Auto-Extending expiration on `Get`s

#### Usage
```go
import (
  "time"
  "github.com/oplehto/ttlcache"
)

func main () {
  cache := ttlcache.NewCache(time.Second)
  cache.Set("key")
  value, exists := cache.Get("key")
  count := cache.Count()
}
```
