DROP TABLE dokument CASCADE;
DROP TABLE foto CASCADE;
DROP TABLE enthalten CASCADE;
DROP TABLE bild CASCADE;
DROP TABLE serie CASCADE;
DROP TABLE ausstellung CASCADE;
DROP TABLE katalog CASCADE;


CREATE TABLE katalog (
    -- id         serial  PRIMARY KEY,
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- code       text    NOT NULL DEFAULT '', -- UNIQUE,      -- secondary key
    jahr       integer NOT NULL DEFAULT 0,
    titel      text    NOT NULL DEFAULT '',
    untertitel text    NOT NULL DEFAULT '',
    kommentar  text    NOT NULL DEFAULT ''
);


CREATE TABLE ausstellung (
    -- id         serial  PRIMARY KEY,
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- code       text    NOT NULL DEFAULT '', -- UNIQUE,      -- secondary key
    titel      text    NOT NULL DEFAULT '',
    untertitel text    NOT NULL DEFAULT '',
    typ        text    NOT NULL DEFAULT '',
    jahr       integer NOT NULL DEFAULT 0,
    von        date    NULL,
    bis        date    NULL,
    ort        text    NOT NULL DEFAULT '',
    venue      text    NOT NULL DEFAULT '',
    kommentar  text    NOT NULL DEFAULT ''
);


CREATE TABLE serie (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- id          serial  PRIMARY KEY,
    slug        text    NOT NULL UNIQUE,      -- secondary key
    jahr        integer NOT NULL DEFAULT 0,
    jahrbis     integer NOT NULL DEFAULT 0,
    titel       text    NOT NULL DEFAULT '',
    untertitel  text    NOT NULL DEFAULT '',
    anzahl      integer NOT NULL DEFAULT 0,
    technik     text    NOT NULL DEFAULT '',
    traeger     text    NOT NULL DEFAULT '',
    hoehe       integer NOT NULL DEFAULT 0,
    breite      integer NOT NULL DEFAULT 0,
    tiefe       integer NOT NULL DEFAULT 0,
    phase       text    NOT NULL DEFAULT '',
    anmerkungen text    NOT NULL DEFAULT '',
    kommentar   text    NOT NULL DEFAULT ''
);

CREATE TABLE bild (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- id          serial  PRIMARY KEY,
    dir         text    NOT NULL DEFAULT '',
    jahr        integer NOT NULL DEFAULT 0,
    phase       text    NOT NULL DEFAULT '',
    titel       text    NOT NULL DEFAULT '',
    foto_id     integer NOT NULL DEFAULT 0,   
    serie_id    integer REFERENCES serie(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    serie_nr    integer NOT NULL DEFAULT 0,
    teile       integer NOT NULL DEFAULT 0,
    technik     text    NOT NULL DEFAULT '',
    traeger     text    NOT NULL DEFAULT '',
    hoehe       integer NOT NULL DEFAULT 0,
    breite      integer NOT NULL DEFAULT 0,
    tiefe       integer NOT NULL DEFAULT 0,
    flaeche     double precision NOT NULL DEFAULT 0.0,
    anmerkungen text    NOT NULL DEFAULT '',
    kommentar   text    NOT NULL DEFAULT '',
    modified    timestamp with time zone not null default now()
);

CREATE TABLE enthalten (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- id             serial  PRIMARY KEY,
    bild_id        integer REFERENCES bild(id)        ON UPDATE CASCADE ON DELETE CASCADE,
    katalog_id     integer REFERENCES katalog(id)     ON UPDATE CASCADE ON DELETE CASCADE,
    ausstellung_id integer REFERENCES ausstellung(id) ON UPDATE CASCADE ON DELETE CASCADE
);



CREATE TABLE foto (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
   -- id        serial      PRIMARY KEY,
    bild_id   integer     REFERENCES bild(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    -- serie_id  integer     REFERENCES serie(slug) ON UPDATE CASCADE ON DELETE RESTRICT
    index     integer     NOT NULL DEFAULT 0,

    name      text        NOT NULL DEFAULT '',
    size      integer     NOT NULL DEFAULT 0,

    uploaded  timestamptz NOT NULL DEFAULT Now(),
    path      text        NOT NULL DEFAULT '',
    format    text        NOT NULL DEFAULT '',
    width     integer     NOT NULL DEFAULT 0,
    height    integer     NOT NULL DEFAULT 0,
    taken     timestamptz NOT NULL,
    caption   text        NOT NULL DEFAULT '',
    kommentar text        NOT NULL DEFAULT ''
);

CREATE TABLE dokument (
    id int GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    -- id        serial     PRIMARY KEY,
    name text,
    format text
);


insert into bild (id, dir, jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified) select id,concat('bilder/', id), jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified from bild2 where id in (8);
insert into bild (id, dir, jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified) select id,concat('bilder/', id), jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified from bild2 where id in (9,10,11,13,14,17,18);
insert into bild (id, dir, jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified) select id,concat('bilder/', id), jahr,phase,titel,foto_id,teile,technik,traeger,hoehe,breite,tiefe,anmerkungen,kommentar,modified from bild2 where id in (105,106,107,108,109,110,111,112,113,115,116,119,120,123,124,125,126,129,130,132,133,134,135,138,140,141,142,144,145,146,148,150,172);



insert into foto (id, bild_id, index, name,size,uploaded,path,format,width,height,taken,caption,kommentar) select id,bild_id,index,name,size,uploaded,path,format,width,height,taken,caption,kommentar from foto2 where bild_id in (8,9,10,11,13,14,17,18,105,106,107,108,109,110,111,112,113,115,116,119,120,123,124,125,126,129,130,132,133,134,135,138,140,141,142,144,145,146,148,150,172);

insert into ausstellung (id, titel,untertitel,typ,jahr,von,bis,ort,venue,kommentar) select id, titel,untertitel,typ,jahr,von,bis,ort,venue,kommentar from ausstellung2 ;
