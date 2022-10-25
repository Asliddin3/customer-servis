CREATE TABLE if NOT exists address(
  id serial PRIMARY key,
  district VARCHAR(100),
  street VARCHAR(100)
);


CREATE Table if NOT exists customer(
  id serial PRIMARY KEY,
  firstname VARCHAR(200),
  lastname VARCHAR(200),
  bio TEXT,
  email VARCHAR(200),
  phonenumber VARCHAR(200),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

CREATE Table if NOT exists customer_address(
  customer_id int REFERENCES customer(id),
  address_id int REFERENCES address(id)
);