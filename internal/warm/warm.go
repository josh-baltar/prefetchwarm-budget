// Package warm computes per-key prefetch warm-up top-ups for a shared
// build-artifact cache.
//
// When space frees up in the shared cache, the warmer decides how many bytes
// to prefetch (top up) for each key out of the freed space. Each key gets a
// proportional slice of the free space (ShareNum/ShareDen). A key can only
// absorb bytes it does not already hold resident (Size minus Resident). A key
// that is served heavily out of live traffic keeps itself warm on its own, so
// the prefetch budget is skewed toward keys the live traffic is not already
// keeping resident; that skew is figured from the key's hit density
// (Hits/Lookups).
package warm

// Budget describes the freed space to distribute this pass.
type Budget struct {
	Free     int // freed bytes available to prefetch this pass
	ShareNum int // this key's share numerator
	ShareDen int // this key's share denominator
}

// Key describes one cache key being considered for a prefetch top-up.
type Key struct {
	Size     int // total bytes the key would occupy fully warm
	Resident int // bytes already resident in cache for this key
	Hits     int // recent lookups that were served from cache
	Lookups  int // recent lookups for this key
}

// WarmTopUp returns how many bytes to prefetch for this key this pass.
func WarmTopUp(b Budget, k Key) int {
	topup := b.Free * b.ShareNum / b.ShareDen

	warm := 0
	if k.Lookups > 0 {
		warm = (b.Free * b.ShareNum / b.ShareDen) * (k.Lookups - k.Hits) / k.Lookups
	}

	return topup - warm
}
