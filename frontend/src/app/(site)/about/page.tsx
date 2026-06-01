import Image from "next/image";
import { PageTransition } from "@/components/PageTransition";
import { SITE } from "@/lib/constants";
import { serverFetch } from "@/lib/server-api";
import type { PublicSiteData } from "@/lib/types";
import { DEFAULT_CONTACT, DEFAULT_SITE_DATA } from "@/lib/constants";
import { ContactStrip } from "@/components/layout/ContactStrip";

export const metadata = { title: "About" };

async function getSite(): Promise<PublicSiteData> {
  const data = await serverFetch<PublicSiteData>("/api/site", 120);
  return data ?? DEFAULT_SITE_DATA;
}

export default async function AboutPage() {
  const site = await getSite();

  return (
    <PageTransition>
      <section className="flex min-h-[50vh] items-end bg-bg-card pb-16 pt-32">
        <div className="section-pad mx-auto w-full max-w-7xl">
          <p className="text-xs tracking-[0.4em] text-gold uppercase">Our Heart</p>
          <h1 className="mt-4 font-display text-5xl md:text-6xl">The Soul Behind {SITE.name}</h1>
        </div>
      </section>

      <section className="section-pad mx-auto max-w-5xl">
        <div className="grid gap-12 lg:grid-cols-[280px_1fr] lg:items-start">
          <div className="relative mx-auto aspect-[3/4] w-full max-w-[280px] overflow-hidden rounded-sm border border-gold/25 bg-bg-muted shadow-[0_12px_40px_rgba(6,20,40,0.5)]">
            {site.owner_portrait_url ? (
              <Image
                src={site.owner_portrait_url}
                alt="Founder of Surya Photography"
                fill
                className="object-cover"
                unoptimized
              />
            ) : (
              <div className="flex h-full flex-col items-center justify-center p-6 text-center text-sm text-text-muted">
                <span className="font-display text-2xl text-gold/40">S</span>
                <p className="mt-4">A portrait of our founder will grace this space soon.</p>
              </div>
            )}
          </div>

          <div>
            <p className="text-lg leading-relaxed text-text-muted">
              At Surya Photography, we don&apos;t merely click shutters—we hold space for the laughter, the tears, the
              glances that say &ldquo;I love you&rdquo; without words. From the first garland to the last farewell at
              the mandapam, we walk beside your family with patience, respect, and an artist&apos;s eye.
            </p>
            <p className="mt-6 text-lg leading-relaxed text-text-muted">
              Based in Rasipuram, we have been blessed to witness countless weddings, puberty ceremonies, baby showers,
              and intimate outdoor stories across Tamil Nadu. Every album we design and every film we edit carries the
              same promise: your memories deserve to feel as beautiful as you felt on that day.
            </p>
          </div>
        </div>

        <div className="mt-20 border-t border-gold/10 pt-16">
          <h2 className="font-display text-2xl text-gold">Reach Us With Love</h2>
          <div className="mt-6 max-w-lg">
            <ContactStrip contact={site.contact} />
          </div>
        </div>

        <div className="mt-16 grid gap-8 border-t border-gold/10 pt-16 sm:grid-cols-3">
          {[
            { label: "Celebrations We've Cherished", value: site.events_covered },
            { label: "Years Devoted to This Craft", value: site.years_of_craft },
            { label: "Families Who Trust Us", value: site.happy_families },
          ].map((s) => (
            <div key={s.label} className="text-center">
              <p className="font-display text-4xl text-gold">{s.value}</p>
              <p className="mt-2 text-xs tracking-widest text-text-muted uppercase">{s.label}</p>
            </div>
          ))}
        </div>
      </section>
    </PageTransition>
  );
}
