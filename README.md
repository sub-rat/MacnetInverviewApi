# Macnet Social Network Api

## Run By Docker
```
docker-compose up -d
```
remove -d if you don't want to run in background


## Run by Makefile
Dependency:

You need to install postgresql in you computer or by docker using the docker-compose file by commenting the api section and enabling the databse only.

```
cp .env.examples .env
```
make neccessary changes in the environment
```
make run
```
This will build the project and run though binary file 

## Postman Collection 
postman collection is in postman_collection folder import it to postman for api document
