import { Pokemon, PokemonListResponse, PokemonSpecies } from "@/app/types/pokemon";

const BASE_URL = "https://pokeapi.co/api/v2";

export async function getPokemonList(limit = 151, offset = 0): Promise<PokemonListResponse> {
  const res = await fetch(`${BASE_URL}/pokemon?limit=${limit}&offset=${offset}`, {
    next: { revalidate: 3600 },
  });
  if (!res.ok) throw new Error("Failed to fetch pokemon list");
  return res.json();
}

export async function getPokemon(nameOrId: string | number): Promise<Pokemon> {
  const res = await fetch(`${BASE_URL}/pokemon/${nameOrId}`, {
    next: { revalidate: 3600 },
  });
  if (!res.ok) throw new Error(`Failed to fetch pokemon: ${nameOrId}`);
  return res.json();
}

export async function getPokemonSpecies(nameOrId: string | number): Promise<PokemonSpecies> {
  const res = await fetch(`${BASE_URL}/pokemon-species/${nameOrId}`, {
    next: { revalidate: 3600 },
  });
  if (!res.ok) throw new Error(`Failed to fetch species: ${nameOrId}`);
  return res.json();
}

export async function searchPokemon(query: string): Promise<Pokemon | null> {
  try {
    return await getPokemon(query.toLowerCase().trim());
  } catch {
    return null;
  }
}
