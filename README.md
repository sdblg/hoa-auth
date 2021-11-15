# hoa-auth

Docker commands
``` shell
make build env=dev
make run env=dev
```

Database install in Mac
1. `brew install --cask dbeaver-community`
2. `brew install postgresql`
3. `brew services start postgresql`
4. `psql postgres`
5. `CREATE ROLE newUser WITH LOGIN PASSWORD ‘password’;`
6. `ALTER ROLE newUser CREATEDB;`
7. `\q`
8. `psql postgres -U newuser`
9. `brew services stop postgresql`