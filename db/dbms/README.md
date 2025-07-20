practically building an dbms:

module1: building query executor
- the node idea is pretty interesting:
  - use linked list data structure
  - read db line by line - perhaps this is how the db work
  - and there is order of execution - it makes sense to me know:
    - scan db first > filter > sort > select | limit > aggregation
    - there are some steps that need the copy of the rows, which feels expensive, but i think we can overcome that

module2: handling reading file, not loading everything from memory:
- the trick is to read line by line, row by row, and not dump everything into the memory
- tdb. note to self: its good that i sit in 2h in a row, but i could have done more, like archive a goal from the beginning
- it s a good start, tho. i just need a goal, and not stray away from it