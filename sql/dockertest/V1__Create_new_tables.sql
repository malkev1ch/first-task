CREATE TABLE cats (
                      id UUID CONSTRAINT cats_primary_key PRIMARY KEY ,
                      name VARCHAR NOT NULL,
                      date_birth DATE NOT NULL,
                      vaccinated BOOLEAN NOT NULL,
                      image_path VARCHAR
);

CREATE TABLE users (
                      id UUID CONSTRAINT users_primary_key PRIMARY KEY ,
                      name VARCHAR NOT NULL,
                      email VARCHAR UNIQUE NOT NULL,
                      password VARCHAR NOT NULL,
                      refresh_token VARCHAR
);