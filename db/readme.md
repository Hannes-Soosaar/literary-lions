The db is create with golang.
The db is opened with the golang application
The sql script is called from within the golang application

To run SQL on the SQLite instance in the lions.db 

navigate to the folder that contains the lions.db file

in the terminal run 
$ sqlite3 lions.db

after

run the following code in the running instance
sqlite> .read test.sql

this will execute the SQL script in the test.sql file.