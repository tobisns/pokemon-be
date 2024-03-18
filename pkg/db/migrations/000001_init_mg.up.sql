
CREATE TABLE pokemons (
    name VARCHAR(32) PRIMARY KEY,
    image_url TEXT,
    evo_tree_id INTEGER DEFAULT -1,
    height INTEGER DEFAULT 0,
    weight INTEGER DEFAULT 0,
    hp INTEGER DEFAULT 0,
    atk INTEGER DEFAULT 0,
    def INTEGER DEFAULT 0,
    sa INTEGER DEFAULT 0,
    sd INTEGER DEFAULT 0,
    spd INTEGER DEFAULT 0
);

CREATE TABLE evo_tree (
    id INTEGER NOT NULL,
    level INTEGER NOT NULL,
    pokemon_name VARCHAR(32),
    PRIMARY KEY (pokemon_name),
    FOREIGN KEY (pokemon_name) REFERENCES pokemons(name) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE type (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

CREATE TABLE pokemon_type (
    pokemon VARCHAR(32) REFERENCES pokemons(name) ON DELETE CASCADE ON UPDATE CASCADE,
    type_id INTEGER REFERENCES type(id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (pokemon, type_id)
);

CREATE TABLE users (
    username VARCHAR(32) PRIMARY KEY,
    password CHAR(64),
    admin BOOLEAN DEFAULT FALSE
);

CREATE OR REPLACE FUNCTION set_evo_tree_id()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE pokemons
        SET evo_tree_id = NEW.id
        WHERE name = NEW.pokemon_name;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE pokemons
        SET evo_tree_id = -1
        WHERE name = OLD.pokemon_name;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_evo_tree_id
AFTER INSERT ON evo_tree
FOR EACH ROW
EXECUTE FUNCTION set_evo_tree_id();

CREATE TRIGGER delete_evo_tree_id
AFTER DELETE ON evo_tree
FOR EACH ROW
EXECUTE FUNCTION set_evo_tree_id();

CREATE SEQUENCE tree_id_seq;

CREATE OR REPLACE FUNCTION new_tree()
RETURNS INTEGER AS $$
DECLARE
    next_id INTEGER;
BEGIN
    -- Get the next value from the sequence
    SELECT nextval('tree_id_seq') INTO next_id;
    
    RETURN next_id;
END;
$$ LANGUAGE plpgsql;
