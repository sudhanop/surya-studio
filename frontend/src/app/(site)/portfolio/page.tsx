import Link from "next/link";
import { serverFetch } from "@/lib/server-api";
import type { Category } from "@/lib/types";
import { PageTransition } from "@/components/PageTransition";

export const metadata = { title: "Portfolio" };

export default async function PortfolioPage() {
  const categories = (await serverFetch<Category[]>("/api/categories", 300)) ?? [];

  return (
    <PageTransition>
      <section className="flex min-h-[50vh] items-end bg-bg-card pb-16 pt-32">
        <div className="section-pad mx-auto w-full max-w-7xl">
          <p className="text-xs tracking-[0.4em] text-gold uppercase">Our Work</p>
          <h1 className="mt-4 font-display text-5xl md:text-6xl">Portfolio</h1>
          <p className="mt-4 max-w-xl text-text-muted">
            Explore our collections across weddings, celebrations, and cinematic portraits.
          </p>
        </div>
      </section>

      <section className="section-pad mx-auto max-w-7xl">
        <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map((cat) => (
            <Link
              key={cat.id}
              href={`/portfolio/${cat.slug}`}
              className="group block overflow-hidden rounded-sm border border-gold/10 bg-bg-card transition hover:border-gold/40"
            >
              <div className="relative flex aspect-[4/3] items-center justify-center overflow-hidden bg-bg-muted">
                {cat.cover_image && (
                  <>
                    <img
                      src={cat.cover_image}
                      alt=""
                      className="absolute inset-0 h-full w-full scale-110 object-cover opacity-30 blur-2xl"
                      aria-hidden="true"
                    />
                    <img
                      src={cat.cover_image}
                      alt={cat.name}
                      className="absolute inset-0 h-full w-full object-cover object-center transition duration-700 group-hover:scale-[1.02]"
                      loading="lazy"
                    />
                  </>
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-bg-deep/95 via-bg-deep/55 to-transparent" />
                <div className="absolute left-4 top-4 flex items-center gap-2">
                  <span className="rounded-full border border-gold/25 bg-bg-deep/60 px-3 py-1 text-[10px] tracking-widest text-gold backdrop-blur">
                    {cat.photo_count ?? 0} Photos
                  </span>
                  <span className="rounded-full border border-gold/25 bg-bg-deep/60 px-3 py-1 text-[10px] tracking-widest text-gold backdrop-blur">
                    {cat.video_count ?? 0} Films
                  </span>
                </div>
                {!cat.cover_image && (
                  <span className="relative font-display text-5xl text-gold/40">{cat.name[0]}</span>
                )}
                <div className="relative mt-auto w-full p-6">
                  <h2 className="font-display text-3xl text-gold">{cat.name}</h2>
                  <p className="mt-2 line-clamp-2 text-sm text-text-muted">{cat.description || "View the full collection"}</p>
                  <p className="mt-4 text-xs tracking-widest text-gold uppercase opacity-80 group-hover:opacity-100">
                    Open collection →
                  </p>
                </div>
              </div>
            </Link>
          ))}
        </div>
      </section>
    </PageTransition>
  );
}
