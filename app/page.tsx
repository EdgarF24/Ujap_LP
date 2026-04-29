import Link from "next/link";
import Image from "next/image";
import type { Metadata } from "next";
import { getPokemon } from "@/app/lib/pokeapi";

export const metadata: Metadata = {
  title: "PokéDex - Inicio | Explora el Mundo Pokémon",
  description:
    "Tu PokéDex completa. Explora más de 1000 Pokémon, busca por nombre o número, descubre stats, habilidades y evoluciones.",
};

const FEATURED_POKEMON = ["pikachu", "charizard", "mewtwo", "eevee", "gengar", "snorlax"];

const STATS = [
  { label: "Pokémon", value: "1000+" },
  { label: "Tipos", value: "18" },
  { label: "Generaciones", value: "9" },
  { label: "Habilidades", value: "300+" },
];

export default async function HomePage() {
  const featured = await Promise.allSettled(
    FEATURED_POKEMON.map((name) => getPokemon(name))
  );

  const featuredPokemon = featured
    .filter((r) => r.status === "fulfilled")
    .map((r) => (r as PromiseFulfilledResult<Awaited<ReturnType<typeof getPokemon>>>).value);

  return (
    <div className="hero-gradient grid-bg min-h-screen">
      {/* Hero */}
      <section className="relative overflow-hidden py-24 px-4">
        {/* Decorative circles */}
        <div className="absolute top-20 right-10 w-72 h-72 rounded-full opacity-5"
          style={{ background: "radial-gradient(circle, #ef4444 0%, transparent 70%)" }} />
        <div className="absolute bottom-10 left-10 w-48 h-48 rounded-full opacity-5"
          style={{ background: "radial-gradient(circle, #f59e0b 0%, transparent 70%)" }} />

        <div className="max-w-4xl mx-auto text-center">
          {/* Badge */}
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full mb-6"
            style={{ background: "rgba(239,68,68,0.1)", border: "1px solid rgba(239,68,68,0.3)" }}>
            <span className="w-2 h-2 rounded-full bg-red-500 animate-pulse" />
            <span className="text-red-400 text-sm font-medium">Pokédex Online</span>
          </div>

          <h1 className="text-5xl md:text-7xl font-extrabold mb-6 leading-tight">
            Explora el{" "}
            <span className="relative">
              <span className="text-transparent bg-clip-text"
                style={{ backgroundImage: "linear-gradient(135deg, #ef4444, #f59e0b)" }}>
                Mundo Pokémon
              </span>
            </span>
          </h1>

          <p className="text-xl text-gray-400 mb-10 max-w-2xl mx-auto leading-relaxed">
            Descubre información detallada de todos los Pokémon: stats, tipos, habilidades,
            movimientos y mucho más. Tu Pokédex definitiva.
          </p>

          <div className="flex flex-wrap gap-4 justify-center">
            <Link
              href="/mundo"
              id="hero-cta-mundo"
              className="px-8 py-4 rounded-xl font-semibold text-white transition-all duration-200 hover:scale-105 hover:shadow-2xl"
              style={{
                background: "linear-gradient(135deg, #ef4444, #dc2626)",
                boxShadow: "0 8px 32px rgba(239,68,68,0.4)",
              }}
            >
              🌍 Mundo Pokémon
            </Link>
            <Link
              href="/buscador"
              id="hero-cta-buscador"
              className="px-8 py-4 rounded-xl font-semibold transition-all duration-200 hover:scale-105"
              style={{
                background: "rgba(255,255,255,0.05)",
                border: "1px solid rgba(255,255,255,0.1)",
                color: "#f1f5f9",
              }}
            >
              🔍 Buscar Pokémon
            </Link>
          </div>
        </div>
      </section>

      {/* Stats bar */}
      <section className="py-12 px-4 border-y"
        style={{ borderColor: "rgba(255,255,255,0.05)", background: "rgba(0,0,0,0.2)" }}>
        <div className="max-w-4xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-6 text-center">
          {STATS.map((stat) => (
            <div key={stat.label}>
              <div className="text-3xl font-extrabold text-red-500">{stat.value}</div>
              <div className="text-sm text-gray-400 mt-1">{stat.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* Featured Pokemon */}
      <section className="py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl font-bold mb-3" style={{ color: "#f1f5f9" }}>
              Pokémon Destacados
            </h2>
            <p className="text-gray-400">Los favoritos de todos los entrenadores</p>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            {featuredPokemon.map((pokemon) => {
              const img =
                pokemon.sprites.other["official-artwork"].front_default ??
                pokemon.sprites.front_default;
              return (
                <Link
                  key={pokemon.name}
                  href={`/pokemon/${pokemon.name}`}
                  id={`featured-${pokemon.name}`}
                  className="pokemon-card glass-card p-4 flex flex-col items-center gap-3 text-center"
                >
                  {img && (
                    <div className="relative w-20 h-20">
                      <Image
                        src={img}
                        alt={pokemon.name}
                        fill
                        sizes="80px"
                        className="object-contain animate-float"
                        style={{ filter: "drop-shadow(0 4px 12px rgba(0,0,0,0.5))" }}
                      />
                    </div>
                  )}
                  <span className="text-sm font-semibold capitalize text-gray-300">
                    {pokemon.name}
                  </span>
                </Link>
              );
            })}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 px-4">
        <div className="max-w-3xl mx-auto text-center glass-card p-12 animate-pulse-glow">
          <div className="text-5xl mb-4">⚡</div>
          <h2 className="text-3xl font-bold mb-4" style={{ color: "#f1f5f9" }}>
            ¿Listo para explorar?
          </h2>
          <p className="text-gray-400 mb-8">
            Más de 1000 Pokémon te esperan. Descubre sus habilidades, tipos y estadísticas únicas.
          </p>
          <Link
            href="/mundo"
            id="bottom-cta-mundo"
            className="inline-block px-10 py-4 rounded-xl font-bold text-white transition-all duration-200 hover:scale-105"
            style={{
              background: "linear-gradient(135deg, #ef4444, #dc2626)",
              boxShadow: "0 8px 32px rgba(239,68,68,0.4)",
            }}
          >
            Ver todos los Pokémon →
          </Link>
        </div>
      </section>
    </div>
  );
}
