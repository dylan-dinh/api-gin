# api-gin

### How to use

You need to provide configuration file to run the program otherwise it won't run.

You have two examples of configuration file, one for running the program itself and one for the tests.

All section and properties need to be fulfilled otherwise it won't run :

- `database` : 
    - Host
    - Port
    - User
    - Password
    - Name


- `site` :
    - MaxPower
    - SiteName

You need to create a database (PostgresSQL) of the name of the section `database` :
- `CREATE DATABASE $name WITH OWNER $user_name`


Only the first argument will be used to open the configuration file.
Then run :
- `go run main.go -c=confFile.ini`

Model `site` and `engine` implement CRUD operations.