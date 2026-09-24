#!/usr/bin/env python3
"""Write an empty AVB vbmeta image with verification and dm-verity disabled.

Same result as `avbtool make_vbmeta_image --flags 3 --padding_size 4096`,
without needing the AOSP avbtool checkout.
"""
import struct
import sys

FLAGS_DISABLE_VERITY_AND_VERIFICATION = 3


def build() -> bytes:
    header = b"AVB0"
    header += struct.pack(">II", 1, 0)            # required libavb version major.minor
    header += struct.pack(">QQ", 0, 0)            # auth + aux data block sizes
    header += struct.pack(">I", 0)                # algorithm: NONE
    header += struct.pack(">10Q", *([0] * 10))    # hash/signature/key/metadata/descriptor offsets+sizes
    header += struct.pack(">Q", 0)                # rollback index
    header += struct.pack(">II", FLAGS_DISABLE_VERITY_AND_VERIFICATION, 0)
    header += b"avbtool 1.1.0".ljust(48, b"\0")   # release string
    header += b"\0" * 80                          # reserved
    assert len(header) == 256
    return header.ljust(4096, b"\0")


if __name__ == "__main__":
    out = sys.argv[1] if len(sys.argv) > 1 else "vbmeta_disabled.img"
    with open(out, "wb") as f:
        f.write(build())
    print(f"wrote {out}")
