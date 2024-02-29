package pokemons

// Pokemon is the nominal object used for interacting with pokemons.
// This represents what is stored in the database.
type Pokemon struct {
	Name          string `json:"name"`
	ImageUrl      string `json:"image_url"`
	EvolutionTree int    `json:"evo_tree_id"`
}

// Pokemons is used to present a list of pokemons in a JSON response.
type Pokemons struct {
	Pokemons []Pokemon `json:"pokemons"`
}

// ArticleCreateUpdate is the request body that is
// accepted for create and updates to articles.
type PokemonCreateUpdate struct {
	Name          string `json:"name" binding:"required"`
	ImageUrl      string `json:"image_url"`
	EvolutionTree int    `json:"evo_tree_id"`
}
