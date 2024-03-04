
DROP TRIGGER IF EXISTS delete_evo_tree_id ON evo_tree;
DROP TRIGGER IF EXISTS update_evo_tree_id ON evo_tree;

DROP FUNCTION IF EXISTS set_evo_tree_id();
DROP FUNCTION IF EXISTS new_tree();

DROP TABLE IF EXISTS pokemon_type;
DROP TABLE IF EXISTS type;
DROP TABLE IF EXISTS evo_tree;
DROP TABLE IF EXISTS pokemons;
DROP TABLE IF EXISTS users;

DROP SEQUENCE IF EXISTS tree_id_seq;