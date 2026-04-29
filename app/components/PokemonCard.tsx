import Link from "next/link";
import Image from "next/image";
import { Pokemon, TYPE_COLORS } from "@/app/types/pokemon";

interface PokemonCardProps {
  pokemon: Pokemon;
}

export default function PokemonCard({ pokemon }: PokemonCardProps) {
  const primaryType = pokemon.types[0]?.type.name ?? "normal";
  const typeColor = TYPE_COLORS[primaryType] ?? "#A8A878";
  const imgSrc =
    pokemon.sprites.other["official-artwork"].front_default ??
    pokemon.sprites.front_default ??
    "/pokeball-placeholder.png";

  const paddedId = String(pokemon.id).padStart(3, "0");

  return (
    <Link href={`/pokemon/${pokemon.name}`} id={`pokemon-card-${pokemon.name}`}>
      <div
        className="pokemon-card glass-card p-4 flex flex-col items-center gap-3 cursor-pointer relative overflow-hidden"
        style={{ minHeight: "220px" }}
      >
        {/* Background blob */}
        <div
          className="absolute inset-0 opacity-10 rounded-2xl"
          style={{
            background: `radial-gradient(circle at center, ${typeColor} 0%, transparent 70%)`,
          }}
        />

        {/* ID */}
        <span
          className="absolute top-3 left-3 text-xs font-bold font-mono"
          style={{ color: "rgba(255,255,255,0.3)" }}
        >
          #{paddedId}
        </span>

        {/* Image */}
        <div className="relative w-24 h-24 mt-2">
          {imgSrc ? (
            <Image
              src={imgSrc}
              alt={pokemon.name}
              fill
              sizes="96px"
              className="object-contain drop-shadow-lg"
              style={{ filter: "drop-shadow(0 4px 12px rgba(0,0,0,0.5))" }}
            />
          ) : (
            <div className="w-full h-full rounded-full bg-gray-700 flex items-center justify-center">
              <span className="text-3xl">?</span>
            </div>
          )}
        </div>

        {/* Name */}
        <h3
          className="capitalize font-bold text-base text-center"
          style={{ color: "#f1f5f9" }}
        >
          {pokemon.name.replace(/-/g, " ")}
        </h3>

        {/* Types */}
        <div className="flex gap-2 flex-wrap justify-center">
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
}
