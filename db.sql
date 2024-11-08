--Buat table transaction categories
CREATE TABLE public.transaction_categories (
	transaction_category_id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	"name" varchar NULL,
	CONSTRAINT transaction_categories_pk PRIMARY KEY (transaction_category_id)
);

-- public.accounts definition

-- Drop table

-- DROP TABLE public.accounts;

CREATE TABLE public.accounts (
	account_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	name varchar NOT NULL,
	balance int8 NULL,
	CONSTRAINT account_pk PRIMARY KEY (account_id)
);

-- public.auths definition

-- Drop table

-- DROP TABLE public.auths;

CREATE TABLE public.auths (
	auth_id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	account_id int8 NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	CONSTRAINT auths_pk PRIMARY KEY (auth_id),
	CONSTRAINT auths_unique UNIQUE (account_id),
	CONSTRAINT auths_unique_username UNIQUE (username)
);

CREATE TABLE public."transaction" (
    transaction_id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
    transaction_category_id int8 NULL,
    account_id int8 NULL,
    from_account_id int8 NULL,
    to_account_id int8 NULL,
    amount int8 DEFAULT 0 NULL,
    transaction_date timestamp NULL,
    CONSTRAINT transaction_pk PRIMARY KEY (transaction_id),
    CONSTRAINT fk_transaction_category FOREIGN KEY (transaction_category_id) REFERENCES public.transaction_category (transaction_category_id)
);

