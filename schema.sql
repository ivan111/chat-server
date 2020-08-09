CREATE TABLE messages (
    message_id SERIAL,
    name varchar(16) NOT NULL,
    message varchar(140) NOT NULL,
    addr varchar(32) NOT NULL,
    color varchar(8) NOT NULL,
    time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (message_id)
);

CREATE TABLE names (
    name_id SERIAL,
    name varchar(16) NOT NULL,

    PRIMARY KEY (name_id)
);


/*
$ psql chat

# \copy names(name) from './names/myouji.tsv' delimiter '	' csv header;
# \copy names(name) from './names/mei.tsv' delimiter '	' csv header;
# \copy names(name) from './names/kana.tsv' delimiter '	' csv header;
# \copy names(name) from './names/adana.tsv' delimiter '	' csv header;
*/
