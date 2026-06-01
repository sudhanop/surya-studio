"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import type { PortfolioMedia } from "@/lib/types";
import { Lightbox } from "./Lightbox";

export function MasonryGallery({ items, animate = true }: { items: PortfolioMedia[]; animate?: boolean }) {
  const photos = items.filter((i) => i.media_type === "photo" && i.url);
  const [lightbox, setLightbox] = useState<number | null>(null);
  const Button: any = animate ? motion.button : "button";

  if (!photos.length) {
    return (
      <p className="py-20 text-center text-text-muted">
        Gallery coming soon. Check back after our next shoot.
      </p>
    );
  }

  return (
    <>
      <div className="masonry">
        {photos.map((item, i) => (
          <Button
            key={item.id}
            type="button"
            {...(animate
              ? {
                  initial: { opacity: 0, y: 20 },
                  whileInView: { opacity: 1, y: 0 },
                  viewport: { once: true },
                  transition: { delay: (i % 6) * 0.05 },
                }
              : {})}
            className="masonry-item group relative w-full overflow-hidden rounded-sm border border-gold/10 bg-bg-muted shadow-[0_10px_40px_rgba(6,20,40,0.35)]"
            onClick={() => setLightbox(i)}
          >
            <img
              src={item.url!}
              alt={item.title || "Photo"}
              className="block h-auto w-full"
              loading="lazy"
            />
            <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/55 via-black/10 to-transparent opacity-0 transition group-hover:opacity-100" />
          </Button>
        ))}
      </div>
      <Lightbox
        items={photos}
        index={lightbox}
        onClose={() => setLightbox(null)}
        onPrev={() => setLightbox((i) => (i === null ? null : (i - 1 + photos.length) % photos.length))}
        onNext={() => setLightbox((i) => (i === null ? null : (i + 1) % photos.length))}
        onSelect={(i) => setLightbox(i)}
      />
    </>
  );
}
