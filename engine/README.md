# engine

Reads a DSIM firmware image, lists the components inside it and matches them against advisories.
Python 3.10 or later, standard library only. Everything here is made up, the format included:
see [FORMAT.md](FORMAT.md).

| File | What |
|---|---|
| `dsimg.py` | `parse(data)` returns the components of an image, or raises `ImageError` |
| `advisories.py` | `load_advisories(path)` and `match(components, advisories)` |
| `analyse.py` | The command line |
| `samples/` | Sample images and advisories |
| `tests/` | The tests |

Run everything from the repository root.

Tests:

```
python -m unittest discover -s engine/tests -t .
```

Analyse an image:

```
python -m engine.analyse engine/samples/good-router.dsim engine/samples/advisories.json
```

It prints JSON and exits 0. A malformed image prints `error: <message>` to standard error and exits 2.
Wrong arguments print the usage to standard error and exit 64.
