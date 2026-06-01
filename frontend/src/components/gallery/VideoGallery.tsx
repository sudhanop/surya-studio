"use client";

import { useState } from "react";
import type { PortfolioMedia } from "@/lib/types";
import { Lightbox } from "./Lightbox";

export function VideoGallery({ items }: { items: PortfolioMedia[] }) {
  const videos = items.filter((i) => i.media_type === "video" && i.url);
  const [index, setIndex] = useState<number | null>(null);

  if (!videos.length) {
    return (
      <p className="py-16 text-center text-text-muted">
        Films will appear here soon.
      </p>
    );
  }

  return (
    <>
      <div className="grid gap-6 md:grid-cols-2">
        {videos.map((v, i) => (
          <button
            key={v.id}
            type="button"
            onClick={() => setIndex(i)}
            className="group overflow-hidden rounded-sm border border-gold/10 bg-bg-card text-left transition hover:border-gold/35"
          >
            <div className="relative aspect-video overflow-hidden bg-bg-muted">
              {v.thumbnail_url ? (
                <img src={v.thumbnail_url} alt={v.title || "Film"} className="h-full w-full object-cover transition duration-700 group-hover:scale-105" />
              ) : (
                <div className="h-full w-full bg-gradient-to-br from-bg-deep via-accent-blue/30 to-bg-muted" />
              )}
              <div className="absolute inset-0 bg-black/35 transition group-hover:bg-black/20" />
              <div className="absolute inset-0 flex items-center justify-center">
                <span className="inline-flex h-14 w-14 items-center justify-center rounded-full border border-gold/40 bg-bg-deep/70 text-gold backdrop-blur">
                  ▶
                </span>
              </div>
            </div>
            <div className="p-5">
              <p className="font-display text-lg text-gold">{v.title || "Cinematic Film"}</p>
              <p className="mt-1 text-sm text-text-muted">Watch in full screen</p>
            </div>
          </button>
        ))}
      </div>

      <Lightbox
        items={videos}
        index={index}
        onClose={() => setIndex(null)}
        onPrev={() => setIndex((i) => (i === null ? null : (i - 1 + videos.length) % videos.length))}
        onNext={() => setIndex((i) => (i === null ? null : (i + 1) % videos.length))}
        onSelect={(i) => setIndex(i)}
      />
    </>
  );
}
