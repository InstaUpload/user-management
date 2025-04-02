# User management microservice.
---
This repo/module deals with everyting releated to users of InstaUpload.

## TODO:
- ### Database.
	[x] Add PostgreSQL database for this microservice.
	[x] Add migration for users table.
	[] Call `database.Setup()` in main.go to migrate database whenever server starts.
- ### Create user.
	[x] Add test to test the functionality to create user.
	[x] Add code to add user in database.
	[] Post user creation an varification email has to be send to user.
- ### Login User.
	[] Add test to test the functionality to login user.
	[] Add code to check if correct user credential are passed.
	[] Add code to genrate JWT token and return the JWT token.
