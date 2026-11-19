"""Command line: python -m engine.analyse <image.dsim> <advisories.json>"""

import json
import sys

from .advisories import load_advisories, match
from .dsimg import ImageError, parse

USAGE = "usage: python -m engine.analyse <image.dsim> <advisories.json>"


def main(argv) -> int:
    if len(argv) != 2:
        print(USAGE, file=sys.stderr)
        return 64
    with open(argv[0], "rb") as handle:
        data = handle.read()
    try:
        components = parse(data)
    except ImageError as error:
        print(f"error: {error}", file=sys.stderr)
        return 2
    result = match(components, load_advisories(argv[1]))
    report = {
        "components": [
            {"name": item.name, "version": item.version, "payloadSize": item.payload_size}
            for item in components
        ],
        "affected": result["affected"],
        "unverified": result["unverified"],
    }
    print(json.dumps(report, indent=2))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
