# Simple Backend Using Go (Golang)

## Folder Structure
- Database: contains connection to remote PostgreSQL and DDL
- handlers: contains several handlers for application
- middleware: contain authorization for JWT Token
- model: contains Database Schema
- utils: extra simple checker for math and string (_just for fun_)

## Tech
- REST API with Gin
- JWT Token for authorization
- JSON request & response
- PostgreSQL as RDBMS
- GORM as database management
- CRUD

## Database
- auth
- account
- transaction
- transaction_category

## API Service
- /auth/login -> Auth Service Auth/Login
- /auth/signup -> Auth Service Auth/Signup
- /account/create
- /account/read
- /account/update
- /account/delete
- /account/list
- /account/my -> Middleware Validate Token to Auth Service Auth/Validate

## Created By
Rizky Indrabayu