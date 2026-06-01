"use client";

import { useEffect, useMemo, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { AnimatePresence, motion } from "framer-motion";
import { useSiteBranding } from "@/components/SiteBrandingProvider";
import type { PortfolioMedia } from "@/lib/types";
import { SITE } from "@/lib/constants";

const fallbackSlides = [
  { url: "", title: "Timeless Weddings" },
  { url: "", title: "Cinematic Stories" },
  { url: "", title: "Cherished Moments" },
];

export function HeroSlideshow({ slides }: { slides: PortfolioMedia[] }) {
  const items = useMemo(
    () =>
      slides.length > 0
        ? slides.map((s) => ({ url: s.url || "", title: s.title || SITE.tagline }))
        : fallbackSlides,
    [slides]
  );

  const [index, setIndex] = useState(0);
  const { logoUrl } = useSiteBranding();

  useEffect(() => {
    if (items.length <= 1) return;
    const id = window.setInterval(() => {
      setIndex((i) => (i + 1) % items.length);
    }, 3200);
    return () => window.clearInterval(id);
  }, [items]);

  const hero = items[index] ?? items[0];

  return (
    <section className="relative flex min-h-screen items-center justify-center overflow-hidden">
      <div className="absolute inset-0">
        <AnimatePresence mode="wait">
          {hero.url ? (
            <motion.div
              key={`${index}-${hero.url}`}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.8 }}
              className="absolute inset-0"
            >
              <Image src={hero.url} alt="" fill priority className="object-cover" unoptimized />
            </motion.div>
          ) : (
            <motion.div
              key={`${index}-fallback`}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.8 }}
              className="h-full w-full bg-gradient-to-br from-bg-deep via-accent-blue/40 to-bg-muted"
            />
          )}
        </AnimatePresence>
        <div className="absolute inset-0 bg-gradient-to-b from-[#061428]/70 via-[#0d1f3c]/80 to-[#061428]/90" />
      </div>

      <div className="relative z-10 mx-auto max-w-4xl px-4 text-center">
        {logoUrl && (
          <div className="mb-6 flex justify-center">
            <Image
              src={logoUrl}
              alt="Surya Photography"
              width={520}
              height={240}
              className="h-24 w-auto object-contain sm:h-28 md:h-32"
              unoptimized
              priority
            />
          </div>
        )}
        <p className="mb-4 text-xs tracking-[0.4em] text-gold uppercase">
          {SITE.name}
        </p>
        <h1 className="font-display text-4xl leading-[0.95] text-text-main sm:text-5xl md:text-7xl lg:text-8xl">
          <AnimatePresence mode="wait">
            <motion.span
              key={`${index}-${hero.title}`}
              initial={{ opacity: 0, y: 14 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -14 }}
              transition={{ duration: 0.5 }}
              className="inline-block"
            >
              {hero.title}
            </motion.span>
          </AnimatePresence>
        </h1>
        <p className="mx-auto mt-6 max-w-xl text-lg text-text-muted">
          {SITE.tagline}
        </p>
        <div className="mt-10 flex flex-wrap justify-center gap-4">
          <Link href="/portfolio" className="btn-primary">View Portfolio</Link>
          <Link href="/contact" className="btn-outline">Book a Session</Link>
        </div>
      </div>
    </section>
  );
}
