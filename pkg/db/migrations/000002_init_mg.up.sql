
INSERT INTO pokemons (name, image_url, evo_tree_id, height, weight, hp, atk, def, sa, sd, spd)
VALUES
    ('ini', 'ini_image_url', NULL, 10, 30, 50, 12, 21, 34, 2, 10),
    ('ibu', 'ibu_image_url', NULL, 2, 31, 12, 2, 80, 45, 3, 11),
    ('budi', 'budi_image_url', NULL, 89, 68, 89, 12, 31, 31, 45, 23);

WITH new_tree_id AS (
    SELECT new_tree() AS id
)
INSERT INTO evo_tree (id, level, pokemon_name)
SELECT id, level, pokemon_name
FROM (
    VALUES
        ((SELECT id FROM new_tree_id), 1, 'ini'),
        ((SELECT id FROM new_tree_id), 2, 'ibu'),
        ((SELECT id FROM new_tree_id), 2, 'budi')
) AS data(id, level, pokemon_name);

INSERT INTO type (name) VALUES ('grass'), ('fire'), ('water');

INSERT INTO pokemon_type (pokemon, type_id) VALUES ('ini', 1), ('ibu', 1), ('budi', 1), ('budi', 2);
