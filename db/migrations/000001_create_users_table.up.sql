CREATE TABLE IF NOT EXISTS files(
    id serial PRIMARY KEY,
    expense_id int NOT NULL REFERENCES expense (id) ON DELETE CASCADE,
    path_file varchar(300) UNIQUE NOT NULL,
    type_file varchar (50) NOT NULL
);