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

### Refactor permission and routing system

- Make custom router struct that uses Mux under the hood YML/JSON file for every route and their params
  - Add route of custom router checks that route is saved in YML/JSON config and panics if not
- AuthMiddleware
  - Match route with conf
  - Do perm and auth checks after
 
