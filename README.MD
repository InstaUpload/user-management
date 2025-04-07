# User management microservice.
---
This repo/module deals with everyting releated to users of InstaUpload.

## TODO:
- ### Test.
	[x] Create a global context that can be used by all the test functions.
	[x] Add MockService in global context.
	[] Add Create user's user in the global context.
	[] Use user from global context to login and add JWT token in global context.
	[x] User JWT token from global context to test AuthUser function.
- ### Database.
	[x] Add PostgreSQL database for this microservice.
	[x] Add migration for users table.
	[] Call `database.Setup()` in main.go to migrate database whenever server starts.
- ### Create user.
	[x] Add test to test the functionality to create user.
	[x] Add code to add user in database.
	[] Post user creation an varification email has to be send to user.
- ### Login User.
	[x] Add test to test the functionality to login user.
	[x] Add code to check if correct user credential are passed.
	[x] Add code to genrate JWT token and return the JWT token.
