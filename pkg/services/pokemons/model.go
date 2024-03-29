package pokemons

// Pokemon is the nominal object used for interacting with pokemons.
// This represents what is stored in the database.
type Pokemon struct {
	Name            string      `json:"name" binding:"required"`
	ImageUrl        string      `json:"image_url"`
	EvolutionTreeID int         `json:"evo_tree_id"`
	Weight          int         `json:"weight"`
	Height          int         `json:"height"`
	Stat            PokemonStat `json:"stat"`
	Types           []Type      `json:"types"`
}

type PokemonLight struct {
	Name      string `json:"name" binding:"required"`
	ImageUrl  string `json:"image_url"`
	EvoTreeId int    `json:"evo_tree_id"`
}

// Pokemons is used to present a list of pokemons in a JSON response.
type Pokemons struct {
	Pokemons []PokemonLight `json:"pokemons"`
}

type PokemonTypeCreateAssign struct {
	Name string `json:"name" binding:"required"`
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

type DeleteEvoData struct {
	Pokemons []Pokemon `json:"pokemons" binding:"required"`
}

type EvolutionCreate struct {
	EvolutionCreateData []EvolutionCreateData `json:"evolution_create"`
}

type Type struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type Types struct {
	Types []Type `json:"types"`
}

type TypePokemon struct {
	Name   string `json:"name" binding:"required"`
	TypeId int    `json:"type_id" binding:"required"`
}

type TypePokemonResponse struct {
	Pokemons []TypePokemon `json:"pokemons"`
}
