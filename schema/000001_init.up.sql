create table if not exists wallets (
                                       id varchar(255) primary key,
                                       balance numeric not null
);

create table if not exists transactions (
                                            time timestamp not null,
                                            sender_id varchar(255) not null,
                                            receiver_id varchar(255) not null,
                                            amount numeric not null
);