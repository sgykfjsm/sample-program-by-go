# Parse SQL file

Read SQL file and format lines to execute the SQL statements one by one.

For example, if the input file contains the following SQL statements,
```sql
CREATE DATABASE IF NOT EXISTS test;
USE test;
SELECT NOW(); SHOW STATUS;
CREATE TABLE IF NOT EXISTS users (
    id BIGINT NOT NULL AUTO_RANDOM PRIMARY KEY,
    name VARCHAR(64) NOT NULL
);
-- EOF
```

output should be below
```sql
CREATE DATABASE IF NOT EXISTS test;
USE test;
SELECT NOW();
SHOW STATUS;
CREATE TABLE IF NOT EXISTS users ( id BIGINT NOT NULL AUTO_RANDOM PRIMARY KEY, name VARCHAR(64) NOT NULL );
-- EOF
```
