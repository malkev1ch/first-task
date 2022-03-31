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

INSERT INTO cats(id, name, date_birth, vaccinated) VALUES ('9d9044a6-d8a8-4e8c-9132-e583d2ebd6c4', 'Some name', '2018-09-22', FALSE);