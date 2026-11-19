# F5: The CSV export is wrong; add XML

**Area:** Go. **Type:** fix, then build. **Points:** 10. Optional.

## The report

From a customer:

> I opened the CSV export in a spreadsheet. On some rows the columns are shifted: part of the title is in the
> severity column. The JSON export of the same findings looks right.

## Then, the request

From an auditor who works with one of our customers:

> Our tool reads XML. Can the export come as XML as well?

## What to do

1. Find out why the CSV export differs from the JSON export, and fix it.
2. Add `GET /api/export?format=xml`, in this shape:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<findings>
  <finding id="the finding's id">
    <title>Telnet service enabled at boot</title>
    <severity>critical</severity>
    <component>telnetd-lite</component>
    <product>Kestrel Dome 2</product>
  </finding>
</findings>
```

A finding with no known component has an empty `component` element.

## Done means

- All three formats hold the same findings, in the same order, with exactly the same values.
- The CSV starts with the header row `id,title,severity,component,product` and reads correctly in any standard
  CSV reader, whatever characters a title holds.
- The XML is well formed, whatever characters a title holds, and is served with an XML content type.
- An unknown format is still refused with 400, and an export still needs a session.
- There are tests that fail on the code as you received it and pass on yours.
