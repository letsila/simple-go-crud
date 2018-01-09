# README #

### What is this repository for? ###

This is a simple Go webapp for my candidature at sphere.ms for the role of Go developer.

### Dependencies ###

You need to install the go-sql-driver from inside your workspace with the following command.

```bash
$ go get -u github.com/go-sql-driver/mysql
```

### Database setup ###

This project uses MySQL. Then, you need to specify your database parameters in the `config.json` file at the root of the project folder.
Create a table named "todos" with the following structure.

```sql
CREATE TABLE `todos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT '',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
```

### Running ###

From the root of the project, run the following command

```bash
$ go run server.go
```
