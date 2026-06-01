"use client";

import { useMemo, useState } from "react";
import { motion } from "framer-motion";
import type { PortfolioMedia } from "@/lib/types";
import { Lightbox } from "./Lightbox";

export function PhotoShowcase({ items }: { items: PortfolioMedia[] }) {
  const photos = useMemo(() => items.filter((i) => i.media_type === "photo" && i.url), [items]);
  const [index, setIndex] = useState<number | null>(null);

  if (!photos.length) {
    return (
      <p className="py-20 text-center text-text-muted">
        Gallery coming soon. Check back after our next shoot.
      </p>
    );
  }

  return (
    <>
      <section className="space-y-6">
        <div className="flex flex-wrap items-end justify-between gap-6">
          <div>
            <p className="text-xs tracking-[0.4em] text-gold uppercase">Gallery</p>
            <h3 className="mt-3 font-display text-2xl text-gold">Full story</h3>
          </div>
          <p className="text-xs text-text-muted">Keyboard: Esc to close · ←/→ to navigate</p>
        </div>

        <div className="masonry">
          {photos.map((item, i) => (
            <motion.button
              key={item.id}
              type="button"
              initial={{ opacity: 0, y: 18 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ delay: (i % 8) * 0.035 }}
              className="masonry-item group relative w-full overflow-hidden rounded-sm border border-gold/10 bg-bg-muted shadow-[0_10px_40px_rgba(6,20,40,0.35)]"
              onClick={() => setIndex(i)}
            >
              <img
                src={item.url!}
                alt={item.title || "Photo"}
                className="block h-auto w-full"
                loading="lazy"
              />
              <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/55 via-black/10 to-transparent opacity-0 transition group-hover:opacity-100" />
              <div className="pointer-events-none absolute bottom-3 left-3 right-3 flex items-end justify-between gap-3 opacity-0 transition group-hover:opacity-100">
                <p className="line-clamp-2 text-left text-xs tracking-widest text-gold uppercase">
                  {item.title || "View"}
                </p>
                <span className="rounded-full border border-gold/30 bg-bg-deep/60 px-3 py-1 text-[10px] tracking-widest text-gold backdrop-blur">
                  Open
                </span>
              </div>
            </motion.button>
          ))}
        </div>
      </section>

      <Lightbox
        items={photos}
        index={index}
        onClose={() => setIndex(null)}
        onPrev={() => setIndex((i) => (i === null ? null : (i - 1 + photos.length) % photos.length))}
        onNext={() => setIndex((i) => (i === null ? null : (i + 1) % photos.length))}
        onSelect={(i) => setIndex(i)}
      />
    </>
  );
}
