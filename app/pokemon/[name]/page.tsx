import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { getPokemon, getPokemonSpecies } from "@/app/lib/pokeapi";
import { TYPE_COLORS } from "@/app/types/pokemon";

interface Props {
  params: Promise<{ name: string }>;
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { name } = await params;
  return {
    title: `${name.charAt(0).toUpperCase() + name.slice(1)} | PokéDex`,
    description: `Datos completos de ${name}: estadísticas, tipos, habilidades, movimientos y más.`,
  };
}

const STAT_LABELS: Record<string, string> = {
  hp: "HP",
  attack: "Ataque",
  defense: "Defensa",
  "special-attack": "Sp. Ataque",
  "special-defense": "Sp. Defensa",
  speed: "Velocidad",
};

const STAT_COLORS: Record<string, string> = {
  hp: "#ef4444",
  attack: "#f97316",
  defense: "#eab308",
  "special-attack": "#3b82f6",
  "special-defense": "#8b5cf6",
  speed: "#22c55e",
};

function StatBar({ name, value }: { name: string; value: number }) {
  const pct = Math.min(100, (value / 255) * 100);
  const color = STAT_COLORS[name] ?? "#ef4444";
  return (
    <div className="flex items-center gap-4">
      <span className="w-28 text-right text-sm text-gray-400 shrink-0">
        {STAT_LABELS[name] ?? name}
      </span>
      <span className="w-8 text-sm font-bold text-gray-200 shrink-0">{value}</span>
      <div className="flex-1 rounded-full h-2"
        style={{ background: "rgba(255,255,255,0.08)" }}>
        <div
          className="h-full rounded-full transition-all duration-1000"
          style={{
            width: `${pct}%`,
            background: color,
            boxShadow: `0 0 8px ${color}88`,
          }}
        />
      </div>
    </div>
  );
}

export default async function PokemonDetailPage({ params }: Props) {
  const { name } = await params;

  let pokemon;
  try {
    pokemon = await getPokemon(name);
  } catch {
    notFound();
  }

  let species = null;
  try {
    species = await getPokemonSpecies(pokemon.id);
  } catch {
    // species optional
  }

  const primaryType = pokemon.types[0]?.type.name ?? "normal";
  const typeColor = TYPE_COLORS[primaryType] ?? "#A8A878";
  const img =
    pokemon.sprites.other["official-artwork"].front_default ??
    pokemon.sprites.other["official-artwork"].front_shiny ??
    pokemon.sprites.front_default;

  const description = species?.flavor_text_entries
    .find((e) => e.language.name === "es")
    ?.flavor_text?.replace(/\f/g, " ") ??
    species?.flavor_text_entries
    .find((e) => e.language.name === "en")
    ?.flavor_text?.replace(/\f/g, " ");

  const genus = species?.genera.find((g) => g.language.name === "es")?.genus ??
    species?.genera.find((g) => g.language.name === "en")?.genus;

  const totalStats = pokemon.stats.reduce((sum, s) => sum + s.base_stat, 0);

  return (
    <div style={{ minHeight: "100vh", background: "#0d0d1a" }}>
      {/* Back button */}
      <div className="max-w-5xl mx-auto px-4 pt-8">
        <Link
          href="/mundo"
          id="back-to-mundo"
          className="inline-flex items-center gap-2 text-gray-400 hover:text-white transition-colors text-sm"
        >
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Volver al Mundo Pokémon
        </Link>
      </div>

      {/* Hero */}
      <div
        className="relative py-16 px-4 overflow-hidden"
        style={{
          background: `linear-gradient(135deg, ${typeColor}22 0%, transparent 60%)`,
        }}
      >
        {/* BG pattern */}
        <div className="absolute inset-0 opacity-5 pointer-events-none"
          style={{
            backgroundImage: `radial-gradient(circle at 80% 50%, ${typeColor} 0%, transparent 50%)`,
          }} />

        <div className="max-w-5xl mx-auto flex flex-col md:flex-row items-center gap-10">
          {/* Image */}
          <div className="relative w-64 h-64 flex-shrink-0">
            <div className="absolute inset-0 rounded-full opacity-20"
              style={{ background: typeColor, filter: "blur(40px)" }} />
            {img ? (
              <Image
                src={img}
                alt={pokemon.name}
                fill
                sizes="256px"
                priority
                className="object-contain animate-float"
                style={{ filter: "drop-shadow(0 8px 32px rgba(0,0,0,0.6))" }}
              />
            ) : (
              <div className="w-full h-full rounded-full bg-gray-800 flex items-center justify-center text-6xl">?</div>
            )}
          </div>

          {/* Info */}
          <div className="flex-1 text-center md:text-left">
            <p className="text-gray-500 text-lg font-mono mb-1">
              #{String(pokemon.id).padStart(3, "0")}
            </p>
            <h1 className="text-5xl font-extrabold capitalize mb-2" style={{ color: "#f1f5f9" }}>
              {pokemon.name.replace(/-/g, " ")}
            </h1>
            {genus && (
              <p className="text-gray-400 mb-4 italic">{genus}</p>
            )}

            {/* Types */}
            <div className="flex gap-2 mb-5 flex-wrap justify-center md:justify-start">
              {pokemon.types.map(({ type }) => (
                <span
                  key={type.name}
                  className="type-badge text-sm px-4 py-1"
                  style={{ background: TYPE_COLORS[type.name] ?? "#A8A878" }}
                >
                  {type.name}
                </span>
              ))}
            </div>

            {/* Description */}
            {description && (
              <p className="text-gray-300 text-base leading-relaxed max-w-lg">
                {description}
              </p>
            )}

            {/* Badges */}
            <div className="flex gap-3 mt-5 flex-wrap justify-center md:justify-start">
              {species?.is_legendary && (
                <span className="px-3 py-1 rounded-full text-xs font-bold"
                  style={{ background: "rgba(234,179,8,0.2)", color: "#fbbf24", border: "1px solid rgba(234,179,8,0.4)" }}>
                  ⭐ Legendario
                </span>
              )}
              {species?.is_mythical && (
                <span className="px-3 py-1 rounded-full text-xs font-bold"
                  style={{ background: "rgba(139,92,246,0.2)", color: "#a78bfa", border: "1px solid rgba(139,92,246,0.4)" }}>
                  ✨ Mítico
                </span>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Details grid */}
      <div className="max-w-5xl mx-auto px-4 pb-16 grid grid-cols-1 lg:grid-cols-2 gap-6 mt-8">
        {/* Quick data */}
        <div className="glass-card p-6">
          <h2 className="text-lg font-bold mb-5" style={{ color: "#f1f5f9" }}>📋 Datos Básicos</h2>
          <div className="grid grid-cols-2 gap-4">
            {[
              { label: "Altura", value: `${(pokemon.height / 10).toFixed(1)} m` },
              { label: "Peso", value: `${(pokemon.weight / 10).toFixed(1)} kg` },
              { label: "Experiencia Base", value: pokemon.base_experience ?? "?" },
              { label: "Generación", value: species?.generation.name.replace("generation-", "Gen ").toUpperCase() ?? "?" },
              { label: "Hábitat", value: species?.habitat?.name ?? "Desconocido" },
            ].map(({ label, value }) => (
              <div key={label} className="p-3 rounded-xl"
                style={{ background: "rgba(255,255,255,0.04)" }}>
                <p className="text-xs text-gray-500 mb-1">{label}</p>
                <p className="font-semibold text-gray-200 capitalize">{String(value)}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Abilities */}
        <div className="glass-card p-6">
          <h2 className="text-lg font-bold mb-5" style={{ color: "#f1f5f9" }}>⚡ Habilidades</h2>
          <div className="flex flex-col gap-3">
            {pokemon.abilities.map(({ ability, is_hidden }) => (
              <div key={ability.name}
                className="flex items-center justify-between p-3 rounded-xl"
                style={{ background: "rgba(255,255,255,0.04)" }}>
                <span className="capitalize font-medium text-gray-200">
                  {ability.name.replace(/-/g, " ")}
                </span>
                {is_hidden && (
                  <span className="text-xs px-2 py-1 rounded-full"
                    style={{ background: "rgba(139,92,246,0.2)", color: "#a78bfa" }}>
                    Oculta
                  </span>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Stats */}
        <div className="glass-card p-6 lg:col-span-2">
          <div className="flex items-center justify-between mb-5">
            <h2 className="text-lg font-bold" style={{ color: "#f1f5f9" }}>📊 Estadísticas Base</h2>
            <span className="text-sm text-gray-400">
              Total: <strong className="text-white">{totalStats}</strong>
            </span>
          </div>
          <div className="flex flex-col gap-4">
            {pokemon.stats.map((s) => (
              <StatBar key={s.stat.name} name={s.stat.name} value={s.base_stat} />
            ))}
          </div>
        </div>

        {/* Moves (first 20) */}
        <div className="glass-card p-6 lg:col-span-2">
          <h2 className="text-lg font-bold mb-5" style={{ color: "#f1f5f9" }}>
            🥊 Movimientos ({pokemon.moves.length})
          </h2>
          <div className="flex flex-wrap gap-2">
            {pokemon.moves.slice(0, 40).map(({ move }) => (
              <span
                key={move.name}
                className="px-3 py-1 rounded-full text-xs capitalize"
                style={{
                  background: `${typeColor}22`,
                  color: "#cbd5e1",
                  border: `1px solid ${typeColor}44`,
                }}
              >
                {move.name.replace(/-/g, " ")}
              </span>
            ))}
            {pokemon.moves.length > 40 && (
              <span className="px-3 py-1 rounded-full text-xs text-gray-500">
                +{pokemon.moves.length - 40} más...
              </span>
            )}
          </div>
        </div>

        {/* Sprites */}
        <div className="glass-card p-6 lg:col-span-2">
          <h2 className="text-lg font-bold mb-5" style={{ color: "#f1f5f9" }}>🖼️ Sprites</h2>
          <div className="flex flex-wrap gap-6 justify-center items-center">
            {[
              { src: pokemon.sprites.front_default, label: "Normal" },
              { src: pokemon.sprites.front_shiny, label: "Shiny" },
              { src: pokemon.sprites.other.dream_world.front_default, label: "Dream World" },
            ].filter(s => s.src).map(({ src, label }) => (
              <div key={label} className="flex flex-col items-center gap-2">
                <div className="relative w-24 h-24 rounded-xl"
                  style={{ background: "rgba(255,255,255,0.04)" }}>
                  <Image
                    src={src!}
                    alt={label}
                    fill
                    sizes="96px"
                    className="object-contain p-2"
                  />
                </div>
                <span className="text-xs text-gray-500">{label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
