"""Load advisories and match them against the components of an image."""

import json


def load_advisories(path) -> list[dict]:
    """Read a JSON list of {"id", "component", "introduced", "fixed", "severity"}.

    `introduced` and `fixed` may be missing or null.
    """
    with open(path, encoding="utf-8") as handle:
        return json.load(handle)


def match(components, advisories) -> dict:
    """Match the components of an image against advisories.

    Returns {"affected": [...], "unverified": [...]}.

    Rules:
    - Component names match case-insensitively and exactly. There are no
      partial matches.
    - A version is dotted integers such as 1.36.1: one or more groups of
      digits separated by dots. Versions compare numerically, group by
      group, with missing groups as 0. So 1.10 > 1.9 and 1.2 == 1.2.0.
    - A component is affected when introduced <= version < fixed. No
      `introduced` (missing or null) means from the beginning. No `fixed`
      (missing or null) means every later version.
    - Each affected item is {"advisory", "component", "version", "severity"}:
      the advisory id, the name and version as written in the image, and
      the severity of the advisory.
    - A component whose name matches an advisory but whose version is not
      dotted integers (for example 2.1-rc1, v3 or empty) cannot be compared.
      It goes to `unverified` as {"advisory", "component", "version"} and
      never to `affected`.
    - An advisory whose own `introduced` or `fixed` is not dotted integers
      raises ValueError.
    - Both lists are sorted by component name (lower-cased), then advisory id.
    """
    raise NotImplementedError("ticket F6")
