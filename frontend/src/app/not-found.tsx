import Link from "next/link";

export default function NotFound() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-bg-deep px-4 text-center">
      <p className="text-xs tracking-[0.4em] text-gold uppercase">404</p>
      <h1 className="mt-4 font-display text-5xl">Page Not Found</h1>
      <p className="mt-4 max-w-md text-text-muted">
        The moment you are looking for is not here — but we can capture yours.
      </p>
      <Link href="/" className="btn-primary mt-10">
        Return Home
      </Link>
    </div>
  );
}
