"""Tests for the advisories."""

import unittest
from pathlib import Path

from engine.advisories import load_advisories

SAMPLES = Path(__file__).resolve().parent.parent / "samples"


class LoadAdvisories(unittest.TestCase):
    def test_sample_file(self):
        advisories = load_advisories(SAMPLES / "advisories.json")
        self.assertEqual([item["id"] for item in advisories], [f"DVS-2026-000{n}" for n in range(1, 7)])
        self.assertEqual(advisories[0], {
            "id": "DVS-2026-0001", "component": "fernhttpd",
            "introduced": "1.2", "fixed": "1.10", "severity": "high",
        })
        self.assertNotIn("fixed", advisories[1])
        self.assertIsNone(advisories[2]["introduced"])


if __name__ == "__main__":
    unittest.main()
