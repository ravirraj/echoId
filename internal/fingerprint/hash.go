package fingerprint

func HashFingerprint(freq1, freq2, delta int) uint64 {

	var hash uint64

	hash |= uint64(freq1&0xFFFF) << 32
	hash |= uint64(freq2&0xFFFF) << 16
	hash |= uint64(delta & 0xFFFF)

	return hash
}
