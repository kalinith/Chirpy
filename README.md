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

GET /api/chirps", apiCfg.getChirps) //fetch all chirps

Examples of Valid URLs
GET http://localhost:8080/api/chirps?sort=asc
GET http://localhost:8080/api/chirps?sort=desc
GET http://localhost:8080/api/chirps



GET /api/chirps/{chirpID}", apiCfg.getChirp) //fetch one chirp by ID
POST /api/chirps", apiCfg.addChirp) //add a Chirp
DELETE /api/chirps/{chirpID}", apiCfg.deleteChirp)//Allow a user to delete a Chirp he owns

### Admin Commands

GET /api/healthz", health) //health check to see if site is ready to receive.
GET /admin/metrics", apiCfg.metrics) //show the server statistics
POST /admin/reset", apiCfg.reset) //reset metrics
POST /api/polka/webhooks", apiCfg.polkaWebhooks)//revoke refresh token
