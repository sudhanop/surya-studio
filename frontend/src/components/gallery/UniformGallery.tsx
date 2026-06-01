"use client";

import { useState } from "react";
import type { PortfolioMedia } from "@/lib/types";

export function UniformGallery({ items, onOpen }: { items: PortfolioMedia[]; onOpen?: (index: number) => void }) {
  const photos = items.filter((i) => i.media_type === "photo" && i.url);
  const [index, setIndex] = useState(0);

  if (!photos.length) {
    return (
      <p className="py-20 text-center text-text-muted">
        Our gallery is being filled with love—beautiful frames from your celebrations will appear here soon.
      </p>
    );
  }

  const current = photos[index];
  const prev = () => setIndex((i) => (i - 1 + photos.length) % photos.length);
  const next = () => setIndex((i) => (i + 1) % photos.length);

  return (
    <div className="space-y-6">
      <div className="relative mx-auto max-w-5xl overflow-hidden rounded-sm border border-gold/20 bg-bg-muted">
        <button type="button" className="block w-full" onClick={() => onOpen?.(index)}>
          <div className="relative aspect-[16/10] max-h-[72vh] w-full bg-bg-muted">
            <img
              src={current.url!}
              alt=""
              className="absolute inset-0 h-full w-full scale-110 object-cover opacity-30 blur-2xl"
              aria-hidden="true"
            />
            <img
              key={current.id}
              src={current.url!}
              alt={current.title || "Gallery photo"}
              className="absolute inset-0 h-full w-full object-cover object-center"
              loading="eager"
            />
          </div>
        </button>
        {onOpen && (
          <button
            type="button"
            onClick={() => onOpen(index)}
            className="absolute bottom-3 left-3 rounded-full border border-gold/30 bg-bg-deep/70 px-4 py-2 text-xs tracking-widest text-gold backdrop-blur hover:bg-gold/10"
          >
            Full screen
          </button>
        )}
        {photos.length > 1 && (
          <>
            <button
              type="button"
              onClick={prev}
              className="absolute left-3 top-1/2 -translate-y-1/2 rounded-full border border-gold/40 bg-bg-deep/80 px-4 py-2 text-sm text-gold backdrop-blur hover:bg-gold/10"
              aria-label="Previous photo"
            >
              ←
            </button>
            <button
              type="button"
              onClick={next}
              className="absolute right-3 top-1/2 -translate-y-1/2 rounded-full border border-gold/40 bg-bg-deep/80 px-4 py-2 text-sm text-gold backdrop-blur hover:bg-gold/10"
              aria-label="Next photo"
            >
              →
            </button>
            <p className="absolute bottom-3 right-4 text-xs tracking-widest text-gold/90">
              {index + 1} / {photos.length}
            </p>
          </>
        )}
        {current.title && (
          <p className="border-t border-gold/10 bg-bg-card px-4 py-3 text-center font-display text-lg text-gold">
            {current.title}
          </p>
        )}
      </div>

      {photos.length > 1 && (
        <div className="flex gap-3 overflow-x-auto pb-2">
          {photos.map((item, i) => (
            <button
              key={item.id}
              type="button"
              onClick={() => setIndex(i)}
              className={`relative h-24 w-24 flex-shrink-0 overflow-hidden rounded-sm border-2 transition ${
                i === index ? "border-gold" : "border-gold/15 opacity-70 hover:opacity-100"
              }`}
            >
              <img src={item.url!} alt="" className="h-full w-full object-cover" loading="lazy" />
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
