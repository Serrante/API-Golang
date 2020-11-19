create database blockcoin;
create user blockcoin with encrypted password 'bl0ckc01n';
-- ALTER TABLE public.transactions OWNER TO blockcoin;
-- ALTER TABLE public.wallets OWNER TO blockcoin;
-- ALTER TABLE public.users OWNER TO blockcoin;

-- grant all privileges on database blockcoin to blockcoin ;
-- grant all privileges on table users to blockcoin ;
-- grant all privileges on table wallets to blockcoin ;
-- grant all privileges on table transactions to blockcoin ;

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	nickname varchar(15) NOT NULL,
	email varchar(40) NOT NULL,
	"password" varchar(100) NOT NULL,
	status bpchar(1) NULL DEFAULT '0'::bpchar,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_id_key UNIQUE (id),
	CONSTRAINT users_nickname_key UNIQUE (nickname)
);

-- Drop table

-- DROP TABLE public.wallets;

CREATE TABLE public.wallets (
	public_key uuid NOT NULL DEFAULT uuid_generate_v4(),
	"user" uuid NOT NULL,
	balance float4 NULL DEFAULT 0.0,
	updated_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT wallets_public_key_key UNIQUE (public_key),
	CONSTRAINT wallet_user_fk FOREIGN KEY ("user") REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);


-- Drop table

-- DROP TABLE public.transactions;

CREATE TABLE public.transactions (
	id uuid NOT NULL DEFAULT uuid_generate_v4(),
	origin uuid NOT NULL,
	target uuid NOT NULL,
	cash float4 NOT NULL,
	message varchar(255) NULL,
	created_at timestamp NOT NULL DEFAULT now(),
	updated_at timestamp NOT NULL DEFAULT now(),
	CONSTRAINT transactions_id_key UNIQUE (id),
	CONSTRAINT transaction_origin_fk FOREIGN KEY (origin) REFERENCES wallets(public_key) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT transaction_targert_fk FOREIGN KEY (target) REFERENCES wallets(public_key) ON UPDATE CASCADE ON DELETE CASCADE
);