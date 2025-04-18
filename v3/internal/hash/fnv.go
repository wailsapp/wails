package hash

import "hash/fnv"

func Fnv(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s)) // Hash implementations never return errors (see https://pkg.go.dev/hash#Hash)
	return h.Sum32()
}
