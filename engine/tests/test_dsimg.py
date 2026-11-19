"""Tests for the DSIM parser."""

import unittest
from pathlib import Path

from engine.dsimg import Component, parse

SAMPLES = Path(__file__).resolve().parent.parent / "samples"


def sample(name):
    return (SAMPLES / name).read_bytes()


class GoodImages(unittest.TestCase):
    def test_router(self):
        self.assertEqual(parse(sample("good-router.dsim")), [
            Component("fernhttpd", "1.9.4", 64),
            Component("linkcfgd", "2.4.1", 0),
            Component("bootwarden", "0.8", 45),
            Component("upnp-relay", "3.0.2", 0),
        ])

    def test_camera(self):
        self.assertEqual(parse(sample("good-camera.dsim")), [
            Component("rtsp-courier", "2.1-rc1", 32),
            Component("libcrypt-lite", "1.2.0", 0),
            Component("linkcfgd", "2.0", 0),
        ])

    def test_empty(self):
        self.assertEqual(parse(sample("empty.dsim")), [])

    def test_example_from_format_md(self):
        data = bytes.fromhex(
            "44 53 49 4D 01 00 03 00 26 00 00 00 01 08 00 6C 69 6E 6B"
            "63 66 67 64 02 05 00 32 2E 34 2E 31 03 04 00 44 53 49 4D"
        )
        self.assertEqual(parse(data), [Component("linkcfgd", "2.4.1", 4)])


if __name__ == "__main__":
    unittest.main()
