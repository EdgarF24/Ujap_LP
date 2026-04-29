"use client";

import { useState, useEffect, useCallback } from "react";
import { getPokemon, getPokemonList } from "@/app/lib/pokeapi";
import { Pokemon } from "@/app/types/pokemon";
import PokemonCard from "@/app/components/PokemonCard";
import LoadingSpinner from "@/app/components/LoadingSpinner";

const LIMIT = 24;

export default function MundoPage() {
  const [pokemonList, setPokemonList] = useState<Pokemon[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [offset, setOffset] = useState(0);
  const [total, setTotal] = useState(0);
  const [selectedType, setSelectedType] = useState<string>("all");

  const TYPES = [
    "all", "fire", "water", "grass", "electric", "psychic",
    "ice", "dragon", "dark", "fairy", "fighting", "poison",
    "ground", "flying", "bug", "rock", "ghost", "steel", "normal",
  ];

  const fetchPokemon = useCallback(async (currentOffset: number, reset = false) => {
    try {
      const list = await getPokemonList(LIMIT, currentOffset);
      if (reset) setTotal(list.count);

      const details = await Promise.allSettled(
        list.results.map((p) => getPokemon(p.name))
      );
      const valid = details
        .filter((r) => r.status === "fulfilled")
        .map((r) => (r as PromiseFulfilledResult<Pokemon>).value);

      setPokemonList((prev) => (reset ? valid : [...prev, ...valid]));
    } catch (err) {
      console.error(err);
    }
  }, []);

  useEffect(() => {
    setLoading(true);
    fetchPokemon(0, true).finally(() => setLoading(false));
  }, [fetchPokemon]);

  const handleLoadMore = async () => {
    const newOffset = offset + LIMIT;
    setOffset(newOffset);
    setLoadingMore(true);
    await fetchPokemon(newOffset);
    setLoadingMore(false);
  };

  const filtered = selectedType === "all"
    ? pokemonList
    : pokemonList.filter((p) =>
        p.types.some((t) => t.type.name === selectedType)
      );

  return (
    <div style={{ minHeight: "100vh", background: "#0d0d1a" }}>
      {/* Header */}
      <div className="grid-bg py-16 px-4 text-center"
        style={{ borderBottom: "1px solid rgba(255,255,255,0.05)" }}>
        <h1 className="text-4xl md:text-5xl font-extrabold mb-3" style={{ color: "#f1f5f9" }}>
          🌍 Mundo <span className="text-red-500">Pokémon</span>
        </h1>
        <p className="text-gray-400 text-lg">
          Explora {total > 0 ? total.toLocaleString() : "todos los"} Pokémon de todas las generaciones
        </p>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* Type filter */}
        <div className="mb-8 overflow-x-auto pb-2">
          <div className="flex gap-2 min-w-max">
            {TYPES.map((type) => (
              <button
                key={type}
                id={`filter-type-${type}`}
                onClick={() => setSelectedType(type)}
                className="px-4 py-2 rounded-full text-sm font-medium capitalize transition-all duration-200"
                style={{
                  background: selectedType === type
                    ? "#ef4444"
                    : "rgba(255,255,255,0.06)",
                  color: selectedType === type ? "#fff" : "#94a3b8",
                  border: "1px solid",
                  borderColor: selectedType === type ? "#ef4444" : "rgba(255,255,255,0.08)",
                }}
              >
                {type === "all" ? "Todos" : type}
              </button>
            ))}
          </div>
        </div>

        {/* Results count */}
        {!loading && (
          <p className="text-gray-500 text-sm mb-6">
            Mostrando {filtered.length} Pokémon
            {selectedType !== "all" ? ` de tipo ${selectedType}` : ""}
          </p>
        )}

        {/* Grid */}
        {loading ? (
          <LoadingSpinner message="Cargando Pokémon..." />
        ) : filtered.length === 0 ? (
          <div className="text-center py-20 text-gray-500">
            <div className="text-5xl mb-4">😞</div>
            <p>No hay Pokémon de tipo <strong>{selectedType}</strong> cargados aún.</p>
            <button
              onClick={() => setSelectedType("all")}
              className="mt-4 px-6 py-2 rounded-xl text-sm font-medium"
              style={{ background: "#ef4444", color: "#fff" }}
            >
              Ver todos
            </button>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
              {filtered.map((pokemon, i) => (
                <div
                  key={pokemon.id}
                  className="animate-fade-in-up"
                  style={{ animationDelay: `${(i % LIMIT) * 30}ms` }}
                >
                  <PokemonCard pokemon={pokemon} />
                </div>
              ))}
            </div>

            {/* Load more */}
            {pokemonList.length < total && selectedType === "all" && (
              <div className="text-center mt-12">
                <button
                  id="load-more-btn"
                  onClick={handleLoadMore}
                  disabled={loadingMore}
                  className="px-10 py-4 rounded-xl font-semibold text-white transition-all duration-200 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
                  style={{
                    background: loadingMore
                      ? "rgba(239,68,68,0.5)"
                      : "linear-gradient(135deg, #ef4444, #dc2626)",
                    boxShadow: "0 8px 32px rgba(239,68,68,0.3)",
                  }}
                >
                  {loadingMore ? (
                    <span className="flex items-center gap-2">
                      <svg className="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                      </svg>
                      Cargando...
                    </span>
                  ) : (
                    `Cargar más (${total - pokemonList.length} restantes)`
                  )}
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
