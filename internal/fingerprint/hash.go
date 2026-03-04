package fingerprint

func HashFingerprint(freq1, freq2, delta int) uint64 {

	var hash uint64

	hash |= uint64(freq1) << 32
	hash |= uint64(freq2) << 16
	hash |= uint64(delta)

	return hash
}