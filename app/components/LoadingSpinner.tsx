export default function LoadingSpinner({ message = "Cargando..." }: { message?: string }) {
  return (
    <div className="flex flex-col items-center justify-center gap-6 py-20">
      <div className="pokeball-spinner" />
      <p className="text-gray-400 text-sm animate-pulse">{message}</p>
    </div>
  );
}
