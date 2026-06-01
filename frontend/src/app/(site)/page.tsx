import Link from "next/link";
import { serverFetch } from "@/lib/server-api";
import type { Category, PortfolioMedia, PublicSiteData } from "@/lib/types";
import { HeroSlideshow } from "@/components/home/HeroSlideshow";
import { MasonryGallery } from "@/components/gallery/MasonryGallery";
import { PageTransition } from "@/components/PageTransition";
import { DEFAULT_CONTACT, DEFAULT_SITE_DATA, DEFAULT_TESTIMONIALS, SITE } from "@/lib/constants";

async function getHomeData() {
  const [featured, latest, categories, site] = await Promise.all([
    serverFetch<PortfolioMedia[]>("/api/portfolio/featured", 120),
    serverFetch<PortfolioMedia[]>("/api/portfolio/latest", 120),
    serverFetch<Category[]>("/api/categories", 300),
    serverFetch<PublicSiteData>("/api/site", 120),
  ]);
  return {
    featured: featured ?? [],
    latest: latest ?? [],
    categories: categories ?? [],
    site: site ?? DEFAULT_SITE_DATA,
  };
}

export default async function HomePage() {
  const { featured, latest, categories, site } = await getHomeData();
  const testimonials =
    site.testimonials?.length > 0 ? site.testimonials : DEFAULT_TESTIMONIALS;
  const wa = DEFAULT_CONTACT.whatsapp
    ? `https://wa.me/${DEFAULT_CONTACT.whatsapp.replace(/\D/g, "")}`
    : "#";

  return (
    <PageTransition>
      <HeroSlideshow slides={featured} />

      <section className="section-pad mx-auto max-w-7xl text-center">
        <p className="text-xs tracking-[0.4em] text-gold uppercase">Welcome, with warmth</p>
        <h2 className="mt-4 font-display text-4xl md:text-5xl">Where Your Love Becomes Legacy</h2>
        <p className="mx-auto mt-6 max-w-2xl leading-relaxed text-text-muted">{SITE.description}</p>
      </section>

      <section className="section-pad bg-bg-card">
        <div className="mx-auto max-w-7xl">
          <h2 className="text-center font-display text-3xl md:text-4xl">Featured Categories</h2>
          <div className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {categories.slice(0, 6).map((cat) => (
              <Link
                key={cat.id}
                href={`/portfolio/${cat.slug}`}
                className="group relative overflow-hidden rounded-sm border border-gold/10 bg-bg-muted p-8 transition hover:border-gold/30"
              >
                {cat.cover_image && (
                  <>
                    <img
                      src={cat.cover_image}
                      alt=""
                      className="absolute inset-0 h-full w-full scale-110 object-cover opacity-25 blur-2xl"
                      aria-hidden="true"
                    />
                    <img
                      src={cat.cover_image}
                      alt=""
                      className="absolute inset-0 h-full w-full object-cover object-center opacity-70 transition duration-700 group-hover:scale-[1.02] group-hover:opacity-80"
                      loading="lazy"
                    />
                    <div className="absolute inset-0 bg-gradient-to-b from-bg-deep/10 via-bg-deep/55 to-bg-deep/80" />
                  </>
                )}
                <div className="relative">
                  <h3 className="font-display text-2xl text-gold">{cat.name}</h3>
                  <p className="mt-2 line-clamp-2 text-sm text-text-muted">{cat.description}</p>
                  <span className="mt-4 inline-block text-xs tracking-widest text-gold uppercase opacity-0 transition group-hover:opacity-100">
                    View Gallery
                  </span>
                </div>
              </Link>
            ))}
          </div>
        </div>
      </section>

      <section className="section-pad mx-auto max-w-7xl">
        <h2 className="text-center font-display text-3xl">Latest Works</h2>
        <div className="mt-12">
          <MasonryGallery items={latest} animate={false} />
        </div>
        <div className="mt-12 text-center">
          <Link href="/portfolio" className="btn-outline">
            Explore Full Portfolio
          </Link>
        </div>
      </section>

      <section className="section-pad bg-bg-card">
        <div className="mx-auto max-w-7xl">
          <h2 className="text-center font-display text-3xl">Customer Reviews</h2>
          <p className="mx-auto mt-4 max-w-2xl text-center text-sm text-text-muted">
            Real words from families who trusted us with their once-in-a-lifetime celebrations.
          </p>
          <div className="mt-12 grid gap-8 md:grid-cols-3">
            {testimonials.map((t) => (
              <blockquote key={t.name} className="glass rounded-sm p-8">
                <p className="italic text-text-muted">&ldquo;{t.text}&rdquo;</p>
                <footer className="mt-4 text-sm text-gold">{t.name}</footer>
              </blockquote>
            ))}
          </div>
        </div>
      </section>

      <section className="section-pad mx-auto max-w-4xl text-center">
        <h2 className="font-display text-4xl">Ready to Celebrate With Us?</h2>
        <p className="mt-4 text-text-muted">Call or message us — we&apos;d be honoured to be part of your story.</p>
        <div className="mt-8 flex flex-wrap justify-center gap-4">
          <a href={wa} target="_blank" rel="noopener noreferrer" className="btn-primary">
            WhatsApp
          </a>
          <Link href="/contact" className="btn-outline">
            Contact Details
          </Link>
        </div>
      </section>

      <section className="border-t border-gold/10 bg-bg-muted py-10">
        <div className="mx-auto flex max-w-4xl flex-col items-center gap-3 px-4 text-center">
          <p className="text-xs tracking-[0.4em] text-gold uppercase">Instagram</p>
          <a
            href={DEFAULT_CONTACT.instagram_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-3 font-display text-lg text-gold hover:underline"
          >
            <svg width="22" height="22" viewBox="0 0 24 24" fill="none" aria-hidden="true">
              <path
                d="M7 2h10a5 5 0 015 5v10a5 5 0 01-5 5H7a5 5 0 01-5-5V7a5 5 0 015-5Z"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinejoin="round"
              />
              <path
                d="M12 16a4 4 0 100-8 4 4 0 000 8Z"
                stroke="currentColor"
                strokeWidth="1.8"
                strokeLinejoin="round"
              />
              <path d="M17.5 6.5h.01" stroke="currentColor" strokeWidth="2.6" strokeLinecap="round" />
            </svg>
            <span>@{DEFAULT_CONTACT.instagram_handle}</span>
          </a>
          <p className="max-w-2xl text-sm text-text-muted">
            Explore more weddings, portraits, and cinematic stories — follow our latest work on Instagram.
          </p>
        </div>
      </section>
    </PageTransition>
  );
}
