# F4: Filter the findings by severity

**Area:** React. **Type:** build. **Points:** 15. Optional.

## The request

From a customer:

> The findings list is long. I only want to see the critical ones first.

## What to build

In `web/src/FindingsTable.tsx`, above the table:

- A drop-down list (a `select`) with the visible label **Severity**, and the choices **All**, **Critical**,
  **High**, **Medium**, **Low**, in that order. It starts on All.
- Choosing a severity shows only the findings of that severity. Choosing All shows every finding again.
- A line that says how many are shown, exactly in this form: `Showing 3 of 12 findings`. Give it
  `role="status"`, so a screen reader announces the change.
- When no finding matches, the table shows one row that says `No findings match.`

## Done means

- It works with the keyboard alone: a person can reach the list with the Tab key and change it with the arrow
  keys, and the list keeps the focus after a change. A native `select` does this.
- The label is tied to the list, so assistive technology reads "Severity" for it.
- The component is still `FindingsTable`, still takes the `findings` prop, and does not change the array it is given.
- The table keeps its caption and its columns.
- There are tests for the filter in `web/src/FindingsTable.test.tsx`.
