# Chirpy

## What is Chirpy
serveMux.Handle("/app/",apiCfg.middlewareMetricsInc(http.StripPrefix("/app",rootHandler))) //Static file content

## API Commands

### User Related Commands

#### /api/users

1. ```POST /api/users```
	accepts a user in the form of e-mail and password.
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

PUT /api/users", apiCfg.passwordUpdate)//revoke refresh token
POST /api/login", apiCfg.login)//login user will eventually return a token.

POST /api/refresh", apiCfg.refresh)//Refresh the access token.
POST /api/revoke", apiCfg.revoke)//revoke refresh token

### Chirp related Commands
the /api/chirps path will be used to interact with chirps

1. ```GET /api/chirps```
	This API call will return all the Chirps in the system.
	It accepts the following queries
		sort: can be asc or desc and will order the chirps by date.
	 	author_id: The uuid of a chirpy user, to only retrieve Chirps for that user.

	Examples of Valid URLs
	```GET http://localhost:8080/api/chirps?sort=asc
	GET http://localhost:8080/api/chirps?sort=desc
	GET http://localhost:8080/api/chirps
	GET http://localhost:8080/api/chirps?author_id=000000-0000000-0000000
```
2. ```GET /api/chirps/{chirpID}```
	fetch a chirp by its ID
3.	```POST /api/chirps```
	Add a Chirp, accepts the chirp in the following format, The chirp will be addded for the currtently logged in user.
	```{
		"body":"this is a chirp"
	}
	```
4. ```DELETE /api/chirps/{chirpID}```
	Allow a user to delete a Chirp he owns, {chirpID} will be the uuid of the chirp the user wishes to delete

### Admin Commands

GET /api/healthz", health) //health check to see if site is ready to receive.
GET /admin/metrics", apiCfg.metrics) //show the server statistics
POST /admin/reset", apiCfg.reset) //reset metrics
POST /api/polka/webhooks", apiCfg.polkaWebhooks)//revoke refresh token
