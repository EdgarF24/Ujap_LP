"use client";

import { useState, useCallback } from "react";
import Image from "next/image";
import Link from "next/link";
import { getPokemon, getPokemonList } from "@/app/lib/pokeapi";
import { Pokemon, TYPE_COLORS } from "@/app/types/pokemon";
import LoadingSpinner from "@/app/components/LoadingSpinner";

export default function BuscadorPage() {
  const [query, setQuery] = useState("");
  const [results, setResults] = useState<Pokemon[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);
  const [error, setError] = useState("");

  const handleSearch = useCallback(async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = query.trim().toLowerCase();
    if (!trimmed) return;

    setLoading(true);
    setError("");
    setResults([]);
    setSearched(true);

    try {
      // Try exact match first
      try {
        const exact = await getPokemon(trimmed);
        setResults([exact]);
        setLoading(false);
        return;
      } catch {
        // Not an exact match, search by name substring
      }

      // Fuzzy search from full list
      const list = await getPokemonList(1302, 0);
      const matches = list.results
        .filter((p) => p.name.includes(trimmed))
        .slice(0, 20);

      if (matches.length === 0) {
        setError(`No se encontró ningún Pokémon con "${query}"`);
        setLoading(false);
        return;
      }

      const details = await Promise.allSettled(matches.map((p) => getPokemon(p.name)));
      const valid = details
        .filter((r) => r.status === "fulfilled")
        .map((r) => (r as PromiseFulfilledResult<Pokemon>).value);

      setResults(valid);
    } catch {
      setError("Error al buscar. Intenta de nuevo.");
    } finally {
      setLoading(false);
    }
  }, [query]);

  return (
    <div style={{ minHeight: "100vh", background: "#0d0d1a" }}>
      {/* Header */}
      <div className="grid-bg py-16 px-4 text-center"
        style={{ borderBottom: "1px solid rgba(255,255,255,0.05)" }}>
        <h1 className="text-4xl md:text-5xl font-extrabold mb-3" style={{ color: "#f1f5f9" }}>
          🔍 <span className="text-red-500">Buscador</span> Pokémon
        </h1>
        <p className="text-gray-400 text-lg mb-10">
          Busca por nombre o número de Pokédex
        </p>

        {/* Search form */}
        <form onSubmit={handleSearch} className="max-w-xl mx-auto flex gap-3">
          <input
            id="pokemon-search-input"
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Ej: pikachu, charizard, 25..."
            className="search-input flex-1 px-5 py-4 text-base"
          />
          <button
            id="pokemon-search-btn"
            type="submit"
            disabled={loading || !query.trim()}
            className="px-6 py-4 rounded-xl font-semibold text-white transition-all duration-200 hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed"
            style={{
              background: "linear-gradient(135deg, #ef4444, #dc2626)",
              boxShadow: "0 4px 20px rgba(239,68,68,0.4)",
              whiteSpace: "nowrap",
            }}
          >
            {loading ? (
              <svg className="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
              </svg>
            ) : (
              "Buscar"
            )}
          </button>
        </form>
      </div>

      <div className="max-w-6xl mx-auto px-4 py-10">
        {/* Loading */}
        {loading && <LoadingSpinner message="Buscando Pokémon..." />}

        {/* Error */}
        {!loading && error && (
          <div className="text-center py-16">
            <div className="text-6xl mb-4">😞</div>
            <p className="text-gray-400 text-lg">{error}</p>
          </div>
        )}

        {/* Initial state */}
        {!loading && !searched && (
          <div className="text-center py-16">
            <div className="text-7xl mb-6">🔮</div>
            <h2 className="text-2xl font-bold text-gray-300 mb-3">¿Qué Pokémon buscas?</h2>
            <p className="text-gray-500">
              Escribe el nombre o número y presiona buscar
            </p>
          </div>
        )}

        {/* Results */}
        {!loading && results.length > 0 && (
          <>
            <p className="text-gray-500 text-sm mb-6">
              {results.length} resultado{results.length !== 1 ? "s" : ""} para &quot;{query}&quot;
            </p>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-5">
              {results.map((pokemon, i) => {
                const img =
                  pokemon.sprites.other["official-artwork"].front_default ??
                  pokemon.sprites.front_default;
                const primaryType = pokemon.types[0]?.type.name ?? "normal";
                const typeColor = TYPE_COLORS[primaryType] ?? "#A8A878";

                return (
                  <Link
                    key={pokemon.id}
                    href={`/pokemon/${pokemon.name}`}
                    id={`search-result-${pokemon.name}`}
                    className="pokemon-card glass-card p-5 flex gap-4 items-center animate-fade-in-up"
                    style={{ animationDelay: `${i * 50}ms` }}
                  >
                    <div className="relative w-20 h-20 flex-shrink-0 rounded-xl overflow-hidden"
                      style={{ background: `${typeColor}22` }}>
                      {img && (
                        <Image
                          src={img}
                          alt={pokemon.name}
                          fill
                          sizes="80px"
                          className="object-contain p-1"
                        />
                      )}
                    </div>
                    <div>
                      <p className="text-xs text-gray-500 font-mono">
                        #{String(pokemon.id).padStart(3, "0")}
                      </p>
                      <h3 className="font-bold capitalize text-gray-100 text-lg">
                        {pokemon.name.replace(/-/g, " ")}
                      </h3>
                      <div className="flex gap-1 mt-1 flex-wrap">
                        {pokemon.types.map(({ type }) => (
                          <span
                            key={type.name}
                            className="type-badge"
                            style={{ background: TYPE_COLORS[type.name] ?? "#A8A878" }}
                          >
                            {type.name}
                          </span>
                        ))}
                      </div>
                    </div>
                  </Link>
                );
              })}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
