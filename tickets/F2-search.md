# F2: Search repeats a product and breaks on an apostrophe

**Area:** Go. **Type:** fix. **Points:** 10. Core.

## The reports

From a customer, Kestrel Cameras:

> On the Products page, the last product on page 1 shows up again at the top of page 2. Page 1 also looks
> longer than the other pages.

From another customer:

> I searched for O'Brien and got "something went wrong".

## Done means

- A page holds at most 5 products, and `hasNext` is true exactly when another page follows.
- Going through the pages shows every product of the organisation once, and no product twice.
- Whatever a person types in the search box is treated as text to look for, and never causes an error.
- Searching still finds products by a part of their name, in any letter case.
- There are tests that fail on the code as you received it and pass on yours.
