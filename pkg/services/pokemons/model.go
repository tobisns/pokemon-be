package pokemons

// Pokemon is the nominal object used for interacting with pokemons.
// This represents what is stored in the database.
type Pokemon struct {
	Name            string      `json:"name"`
	ImageUrl        string      `json:"image_url"`
	EvolutionTreeID int         `json:"evo_tree_id"`
	Weight          int         `json:"weight"`
	Height          int         `json:"height"`
	Stat            PokemonStat `json:"stat"`
}

// Pokemons is used to present a list of pokemons in a JSON response.
type Pokemons struct {
	Pokemons []Pokemon `json:"pokemons"`
}

type PokemonStat struct {
	HealthPoint    int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}

type PokemonDeleteResponse struct {
	Name    string `json:"name" binding:"required"`
	Message string `json:"message"`
}

// ArticleCreateUpdate is the request body that is
// accepted for create and updates to articles.
type PokemonCreateUpdate struct {
	Name     string      `json:"name" binding:"required"`
	ImageUrl string      `json:"image_url"`
	Weight   int         `json:"weight"`
	Height   int         `json:"height"`
	Stat     PokemonStat `json:"stat"`
}

type EvolutionData struct {
	ID          int    `json:"id" binding:"required"`
	Level       int    `json:"level"`
	PokemonName string `json:"pokemon_name"`
}

type EvolutionCreateData struct {
	Level       int    `json:"level" binding:"required"`
	PokemonName string `json:"pokemon_name" binding:"required"`
}

type EvolutionTree struct {
	EvolutionData []EvolutionData `json:"evolution_data"`
}

type EvolutionCreate struct {
	EvolutionCreateData []EvolutionCreateData `json:"evolution_create"`
}
