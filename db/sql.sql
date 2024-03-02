CREATE TABLE users (
	id serial PRIMARY KEY,
	name varchar (50) NOT NULL,
	surname varchar (25) NOT NULL,
	login varchar (25) NOT NULL,
	pass varchar (30) NOT NULL,
	email varchar (30) NOT NULL
	
);

CREATE TABLE expense_type (
	id serial PRIMARY KEY,
	users_id integer NOT NULL references users(id),  
	type_expenses varchar (50) NOT NULL
	
);

CREATE TABLE expense (
id serial,
expense_type_id integer NOT NULL references expense_type(id),  
reated_at timestamp default now() NOT NULL,
	spent_money numeric(16,2) NOT NULL
);



