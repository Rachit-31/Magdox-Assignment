"""Parser for DSIM firmware images. The format is described in FORMAT.md."""

import struct
from dataclasses import dataclass

MAGIC = b"DSIM"
NAME, VERSION, PAYLOAD = 0x01, 0x02, 0x03


class ImageError(ValueError):
    """The bytes are not a valid DSIM image."""


@dataclass(frozen=True)
class Component:
    name: str
    version: str
    payload_size: int


def parse(data: bytes) -> list[Component]:
    """Return the components of a DSIM image, or raise ImageError."""
    magic, version, flags, count, total = struct.unpack("<4sBBHI", data[:12])
    if magic != MAGIC:
        raise ImageError("not a DSIM image")

    names, versions, payloads = [], [], {}
    offset = 12
    for _ in range(count):
        kind = data[offset]
        (length,) = struct.unpack_from("<H", data, offset + 1)
        value = data[offset + 3:offset + 3 + length]
        offset += 3 + length
        if kind == NAME:
            names.append(value.decode("utf-8"))
        elif kind == VERSION:
            versions.append(value.decode("utf-8"))
        elif kind == PAYLOAD:
            payloads[len(names) - 1] = len(value)

    return [
        Component(name, version, payloads.get(index, 0))
        for index, (name, version) in enumerate(zip(names, versions))
    ]
