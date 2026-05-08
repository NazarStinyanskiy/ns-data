# ns-data
Relational Database Management System

## Compilation

<code>make build</code> - will compile a program and put binary into ./bin/nsdata

## Run

### Help
<code>./bin/nsdata --help</code>

### Create a DB and connect
<code>./bin/nsdata db init {db_name}</code>
<br>
<code>./bin/nsdata db connect {db_name}</code>

### Available commands
<code>CREATE TABLE ( {column_name} {column_type} {column_modifier} , ... )</code> - Please, separate all parts with spaces. Will be adjusted in the future.


<code>LIST TABLES</code>


<code>DESCRIBE TABLE {table_name}</code>
