-- Account Table
CREATE TABLE IF NOT EXISTS account
(
    account_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 0 MINVALUE 0 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    balance bigint NOT NULL,
    referral_account_id bigint,
    CONSTRAINT account_pkey PRIMARY KEY (account_id),
    CONSTRAINT account_referral_account_id_fkey FOREIGN KEY (referral_account_id)
        REFERENCES auth (account_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

-- Auth Table
CREATE TABLE IF NOT EXISTS auth
(
    auth_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    account_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    username character varying COLLATE pg_catalog."default" NOT NULL,
    password character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT auth_pkey PRIMARY KEY (auth_id),
    CONSTRAINT auth_account_id_key UNIQUE (account_id),
    CONSTRAINT auth_username_key UNIQUE (username)
)

-- Transaction Table
CREATE TABLE IF NOT EXISTS transaction
(
    transaction_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    transaction_category_id bigint,
    account_id bigint NOT NULL,
    from_account_id bigint,
    to_account_id bigint,
    amount bigint NOT NULL,
    transaction_date timestamp with time zone,
    CONSTRAINT transaction_pkey PRIMARY KEY (transaction_id),
    CONSTRAINT transaction_transaction_category_id_fkey FOREIGN KEY (transaction_category_id)
        REFERENCES transaction_category (transaction_category_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)

-- Transaction_Category Table
CREATE TABLE IF NOT EXISTS transaction_category
(
    transaction_category_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
    name character varying COLLATE pg_catalog."default",
    CONSTRAINT transaction_category_pkey PRIMARY KEY (transaction_category_id)
)