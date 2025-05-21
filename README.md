# Chirpy

## What is Chirpy
Chirpy is a guided project for learning HTTP servers.
It uses a combination of static content and APi's to serve this purpose.
## API Commands

### User Related Commands

#### /api/users

1. ```POST /api/users```
	accepts a user in the form of an email and a password.
	```
	{
		"email":"user@example.com"
		"password":"plain text password"
	}
	```
	and returns a 201 with the user in the format
	```
	{
		"id":"UUID"
		"created_at":"date"
		"updated_at":"date"
		"email":"example@domain.com"
		"is_chirpy_red":"False"
	}
	```
2. ```PUT /api/users```
	This call will reset the users email and password. It requres an auth token in the header as well as a body in the format
	```
	{
		"email":"user@example.com"
		"password":"plain text password"
	}
	```
3. ```POST /api/login```
   	accepts a login and returns a refresh token.
4. ```POST /api/refresh```
	Accepts a refresh token and returns an access token.
5. ```POST /api/revoke```
	This will revoke a refresh token

### Chirp-related Commands
The **/api/chirps** path will be used to interact with chirps

1. ```GET /api/chirps```
	This API call will return all the Chirps in the system.
	It accepts the following queries
		sort: can be asc or desc and will order the chirps by date.
	 	author_id: The UUID of a chirpy user. When included, it will retrieve Chirps for that user.

	Examples of Valid URLs
	```GET http://localhost:8080/api/chirps?sort=asc
	GET http://localhost:8080/api/chirps?sort=desc
	GET http://localhost:8080/api/chirps
	GET http://localhost:8080/api/chirps?author_id=000000-0000000-0000000
	```
2. ```GET /api/chirps/{chirpID}```
	fetch a chirp by its ID
3. ```POST /api/chirps```
	Add a Chirp, accepts the chirp in the following format, The chirp will be addded for the currtently logged in user.
	```{
		"body":"this is a chirp"
	}
	```
4. ```DELETE /api/chirps/{chirpID}```
	Allow a user to delete a Chirp he owns, {chirpID} will be the uuid of the chirp the user wishes to delete

### Admin Commands

1. GET /api/healthz
	health check to see if site is ready to receive.
2. GET /admin/metrics
	show the server statistics
3. POST /admin/reset
	reset metrics and clean the Database
4. POST /api/polka/webhooks
	webhook endpoint for Polka
