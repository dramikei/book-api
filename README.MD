# About

REST API that I built using Golang, Echo, GORM, MySQL while learning Go


# Installation

Clone into -

```
{GOPATH}/src/
```

Create MySQL DB -

```
CREATE database Library;
use Library;
create table Books(
	id int unsigned auto_increment primary key,
    name varchar(255),
    author varchar(255),
    qty int
);
```


