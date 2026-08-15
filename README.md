# prefetchwarm-budget

A small helper that decides per-key prefetch warm-up top-ups for a shared
build-artifact cache. When space frees up, each key gets a proportional slice
of the freed space, bounded by what the key can still absorb, and skewed by how
much the key already keeps itself warm from live traffic.
