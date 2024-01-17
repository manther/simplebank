Simple bank is just a practice project.

# Data
It has a simple db schema created at first on https://dbdiagram.io/
The schema was modeled there, and exported to a postgres file called Simple Bank app.sql.

Postgres16 was used locally to develop against.

Sqlc was the ORM used to generate model structs, the beginnings of some basic crud go functions, and a set of DB migration sql files. 
Additional transactional, deadlocking optimaizations are made to the "query".sql files and then the makefile target "sqlc" can be executed to rebuild the orm operations.
```
make sqlc
```
Each crud orm file has a test file associated with it. The majority of operations should be show test coverage. 
There is a util folder with some testing utilities mainly around generating random values so far. 

The makefile also supports:
DB migration
```
make migrateup
make migratedown
make dbdrop
```
Launching Postgrest container locally
```
make postgres
```
Preparing DB
```
make createdb
make migrateup
```