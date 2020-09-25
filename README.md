# BackAPI

## Introduction

The purpose of this repository is to be usable as a skeleton of Web API. It contains the main components initialized and running without any further complexity nor
any domain-specific code.

It can be used to start any project from the blog to a E-Commerce website or a SaaS.

## Content

Currently, the project contains :
- [x] **Global**
- Configuration handling with files and environment variables. uses [Viper](https://github.com/spf13/viper)
- A running CLI to start the project. uses [Cobra](https://github.com/spf13/cobra)
- [x] **HTTP**
- A server initialized and running. uses [Echo](https://github.com/labstack/echo)
- CRUD Endpoints for users
- Sessions system with login/logout
- Admin system
- Authentication, Authorization and Admin middlewares
- Payloads validation
- Response formatter
- [x] **Database**
- A Postgresql implementation running
- Migration system.
- A basic User model with all the required methods to interact with it
- A basic Session model with all the required methods to interact with it
- [x] **Misc**
- a JWT interface to create easily JWT token if needed

## How to add a new entity

Adding a new entity is really straightforward and should only require very specific independent steps : 
- Add the entity in the migration / model (under the PG folder.)
- Add the corresponding part in the root of the project (corresponding to the controller)
- Finally add the corresponding endpoints in the HTTP part
_More detailled version to come..._

## Other

The code is tested from end-to-end and should be usable as is to start a new project.
That's also why it will not implement that much more new features to keep it simple to run and understand.
