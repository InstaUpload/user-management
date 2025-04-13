# User management microservice.
---
This repo/module deals with everyting releated to users of InstaUpload.

## TODO:
- ### Test.
	[x] Create a global context that can be used by all the test functions.
	[x] Add MockService in global context.
	[x] Add Create user's user in the global context.
	[x] Use user from global context to login and add JWT token in global context.
	[x] User JWT token from global context to test AuthUser function.
	[] Add Fail test for `SendVerification` function.
	[] Add Fail test for `Verify` function.
- ### Database.
	[x] Add PostgreSQL database for this microservice.
	[x] Add migration for users table.
	[ ] Call `database.Setup()` in main.go to migrate database whenever server starts.
- ### Create user.
	[x] Add test to test the functionality to create user.
	[x] Add code to add user in database.
	[ ] Post user creation an varification email has to be send to user.
	[ ] Add function to create superadmin users, and update the role to admin.
- ### Login User.
	[x] Add test to test the functionality to login user.
	[x] Add code to check if correct user credential are passed.
	[x] Add code to genrate JWT token and return the JWT token.
- ### Update User Role.
	[x] Add test to test the functionality to update the user role.
	[x] Add database migration files to add and remove roles table and alter users table.
	[ ] Update Auth function to return user role as well.
	[ ] Update gRPC Auth function response to send role as well.
	[x] Add function in storage layer to update the role based on the given string.
