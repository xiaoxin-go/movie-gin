#### mysql
```
create database movie charset=utf8;
create table t_actress (
    id int primary key auto_increment, 
    name varchar(50) not null unique,
    age int not null,
    height varchar(20),
    cup varchar(10),
    birthday datetime
);
    
create table t_film(
    id int primary key auto_increment, 
    name varchar(20) not null unique,
    title varchar(200) not null,
    release_date datetime not null,
    length varchar(20) not null,
    image varchar(100)
);

create table t_actress_film(
    id int primary key auto_increment,
    actress_id int,
    film_id int
);
    
create table t_genre(
    id int primary key auto_increment,
    name varchar(50) not null unique
);
    

create table t_genre_film(
    id int primary key auto_increment,
    film_id int,
    genre_id int
);  

    
create table t_link(
    id int primary key auto_increment,
    film_id int,
    magnet varchar(100) not null,
    name varchar(100) not null,
    size varchar(20) not null,
    share_date datetime
);
    
create table t_image(
    id int primary key auto_increment,
    film_id int,
    name varchar(100)
);


##### get movie
title .container>h3
info .container>.movie[0]>.info
release_date datetime not null,
length varchar(20) not null,
actress_id int

image .container>.movie>.screencap>.bigImage href

link #magnet-table>tr>td

image-small #sample-waterfall>a href