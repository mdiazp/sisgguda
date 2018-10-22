drop table users;

create table users(
  id         serial primary key,
  username   varchar(100) not null unique, 
  isAdmin    boolean not null default false,
  createdAt timestamp not null  
);