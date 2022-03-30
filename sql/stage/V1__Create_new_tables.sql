CREATE TABLE cats (
                      id UUID CONSTRAINT cats_primary_key PRIMARY KEY ,
                      name VARCHAR,
                      date_birth DATE,
                      vaccinated BOOLEAN,
                      image_path VARCHAR
);

CREATE TABLE users (
                       id UUID CONSTRAINT users_primary_key PRIMARY KEY ,
                       name VARCHAR,
                       email VARCHAR UNIQUE,
                       password VARCHAR,
                       refresh_token VARCHAR
);