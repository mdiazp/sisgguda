drop table users cascade;
drop table groups cascade;
drop table group_specialists cascade;
drop table group_adusers cascade;

create table users(
  id          serial primary key,
  username    varchar(100) not null unique, 
  description varchar(300), 
  rol         varchar(100) not null, 
  createdAt   timestamp not null  
);

create table groups(
  id          serial primary key,
  name        varchar(100) not null unique, 
  description varchar(300),
  createdAt   timestamp not null  
);

create table group_specialists(
  group_id     integer not null references groups(id) ON DELETE CASCADE,
  user_id    integer not null references users(id) ON DELETE CASCADE,

  primary key(group_id,user_id)
);

create table group_adusers(
  group_id     integer not null references groups(id) ON DELETE CASCADE,
  ad_username    varchar(100) not null,

  primary key(group_id,ad_username)
);

INSERT INTO users 
    (username, description, rol, createdAt) 
VALUES 
    ('sisgguda','Usuario superadmin','SuperAdmin',current_timestamp),
    ('manuel.diaz','Desarrollador DI','Admin',current_timestamp);