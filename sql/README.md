Fact time tracking: 4h

PostgreSQL limitations:
1) EXT tablefunc with crosstab:
mandatory of specification the column names during column creation
2) simple queries and views: no dynamic columns names (rows/corteges only)

Possible solutions:
1) (using others DBMS) pivots https://learn.microsoft.com/en-us/sql/t-sql/queries/from-using-pivot-and-unpivot?view=sql-server-ver16
2) (PostgreSQL) convert into json and represent as a view

No judgement here, but the task itself produces a lot of questions (why? what the purpose to retrieve this report in such way?)
The trooth to be told, the achievement of creation such solution looks like academic learning program.

I did find aboved limitations and took the liberty to leave explanation in that file.
IMHO, here another layer is requiered, but most likely it is additional time/others needs.