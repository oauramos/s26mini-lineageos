package main

import "encoding/binary"

// disabledVBMeta builds an empty AVB vbmeta image with verification and dm-verity disabled.
// Same bytes as tools/make_vbmeta_disabled.py (`avbtool make_vbmeta_image --flags 3 --padding_size 4096`).
func disabledVBMeta() []byte {
	b := make([]byte, 4096)
	copy(b, "AVB0")
	binary.BigEndian.PutUint32(b[4:], 1) // required libavb major version (minor 0)
	// auth/aux sizes, algorithm, offsets, rollback index: all zero
	binary.BigEndian.PutUint32(b[120:], 3) // flags: disable verity + verification
	copy(b[128:], "avbtool 1.1.0")
	return b
}
