# F6: Damaged firmware images; advisory matching (bonus)

**Area:** Python. **Type:** fix, then build. **Points:** 35 (part A 15, part B 20). Bonus.

`engine/` reads firmware images in a made-up format, DSIM, described in `engine/FORMAT.md`. An image lists the
software components inside a device. Sample images are in `engine/samples/`.

## Part A: the reader fails on damaged images (fix, 15 points)

From the support team:

> Customers upload images that were cut short during download, or damaged. The reader works on good images.
> On damaged ones it sometimes crashes with a Python error we cannot show to a customer, and sometimes it
> returns components that are not really there. Try the `broken-` samples.

Done means, for `parse(data)` in `engine/dsimg.py`:

- A good image gives the same components as today: `Component` objects with `name`, `version` and `payload_size`.
- Every image that breaks a rule in `engine/FORMAT.md` raises `ImageError`, with a short message that says
  what is wrong and where. No other kind of error ever comes out, whatever the bytes are.
- The reader never uses bytes beyond the end of the data, and never returns a value shorter than the image
  says it is.
- There are tests for the damaged cases.

## Part B: which components are affected by an advisory (build, 20 points)

`match(components, advisories)` in `engine/advisories.py` is not written yet. Its rules are in its docstring.
Write it, with tests. `engine/samples/advisories.json` holds made-up advisories to try it with.

When it works, this prints the components of an image and the advisories that affect them:

```
python -m engine.analyse engine/samples/good-router.dsim engine/samples/advisories.json
```

Done means:

- Every rule in the docstring holds. A version such as `1..2`, `1.2.`, `.5`, `1.x` or `2.1rc1` is not dotted integers.
- The command above exits with 0 for a good image and prints JSON with the keys `components` (each with `name`,
  `version` and `payloadSize`), `affected` and `unverified`.
- For a damaged image it exits with 2, prints `error: ` and the reason on standard error, and prints nothing
  on standard output.
- Called with the wrong arguments it exits with 64 and prints a usage line on standard error.

Standard library only. Keep the folder a package (`engine/__init__.py`), and run tests from the repository root.
