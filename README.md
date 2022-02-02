#### - Install Golang engine 1.17
https://go.dev/doc/install
***
##### - Install migration tools
https://github.com/golang-migrate/migrate

For linux:
~~~~
curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
apt-get update
apt-get install -y migrate
~~~~

###### - Install Postgres
https://www.digitalocean.com/community/tutorials/how-to-install-postgresql-on-ubuntu-20-04-quickstart-ru

###### **Add user to postgres**

https://medium.com/coding-blocks/creating-user-database-and-adding-access-on-postgresql-8bfcd2f4a91e

###### Migrations

migrate create -ext sql -dir migrations create_users

migrate -path migrations -database "postgresql://root:1@localhost:5432/venera" up