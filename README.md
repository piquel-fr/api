# piquel.fr

The [piquel.fr](https://piquel.fr) website.

## Air does no work on Windows

## TODO

- Setup proper error handling
  - Send to user
  - Panic some
- Auth
  - Fix google not working as auth provider
  - make sure users cant access settings page and other logged in only pages
- Add 404 and 405 pages
- Setup Webauthn for passkeys
- Setup database connection pooling

### Refactor permission and routing system

- AuthMiddleware
  - Match route with conf
  - Do perm and auth checks after
 
