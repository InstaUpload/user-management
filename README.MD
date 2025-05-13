# User management microservice.
---
This repo/module deals with everyting releated to users of InstaUpload.

## Note:
- *For config structs, no need to send pointer* if config struct is not expected to update(mutate) the values. *Only send struct as is*, even if it means creating a shallow copy in memory.
- On the *authenticated function user is expected in context*, Make sure that full user object is passed on as in functions like `SendVerification` user information is required *This is for service layer*.
- It seems context with custome value is not passed down from gRPC client to gRPC server, so we need to pass the user object in every function that is expecting user in context. As well as we need to convert `pb.AuthUserResponse` to `type.User` on every request that is expecting user in context.

## TODO:

- ### In handler layer, For every authenticated function *in* parameter should take in *AuthUserResponse and this user will be stored in context which then will be passed on to service layer.
	- [x] Implement this in SendVerification function. 
	- [ ] Implement this in SendEditorRequest function. 
	- [ ] Implement this in AddEditor function. 

- ### Maybe need to convert pb.AuthUserResponse to type.User on every request that is expecting user in context.
	- [x] Implement this in SendVerification function.
	- [ ] Implement this in SendEditorRequest function. 
	- [ ] Implement this in AddEditor function. 

- ### Test.
	- [x] Create a global context that can be used by all the test functions.
	- [x] Add MockService in global context.
	- [x] Add Create user's user in the global context.
	- [x] Use user from global context to login and add JWT token in global context.
	- [x] User JWT token from global context to test AuthUser function.
	- [ ] Add Fail test for `SendVerification` function.
	- [x] Add Fail test for `Verify` function.
	- [ ] Create test_mock.go in root of project to test kafka and mock email service.
	- [ ] Create test for kafka functions as well?

- ### Database.
	- [x] Add PostgreSQL database for this microservice.
	- [x] Add migration for users table.
	- [x] Call `database.Setup()` in main.go to migrate database whenever server starts.

- ### Create user.
	- [x] Add test to test the functionality to create user.
	- [x] Add code to add user in database.
	- [x] Add function to create superadmin users, and update the role to admin.
	- [ ] After user is successfuly created maybe in gateway service redirect user to `login` then to `SendVerification` endpoint.

- ### Login User.
	- [x] Add test to test the functionality to login user.
	- [x] Add code to check if correct user credential are passed.
	- [x] Add code to genrate JWT token and return the JWT token.

- ### Update User Role.
	- [x] Add test to test the functionality to update the user role.
	- [x] Add database migration files to add and remove roles table and alter users table.
	- [x] Update Auth function to return user role as well.
	- [x] Update gRPC Auth function response to send role as well.
	- [x] Add function in storage layer to update the role based on the given string.

- ### Message Broker.
	- [x] Add Welcome email to be send when user creates an account.
		- [x] Add functionality to send Email.
		- [x] Add welcome message template.
	- [x] Add functionality to send varification mail to user on visiting `SendVerification` function.
		- [x] Add welcome message template.

- ### AddEditor User.
	- [ ] Re-write function to add user, Basically creater should click on add button which will then send notification to User(editor) that creater __*name*__ has requested to add you as a editor.
		- [x] Add functionality to genrate token.
		- [ ] Add functionality to send email to editor.

- ### Refactor.
	- [ ] Have JWT interface at service struct and pass it to User service.
