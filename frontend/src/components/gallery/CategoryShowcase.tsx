"use client";

import { useMemo, useState } from "react";
import type { PortfolioMedia } from "@/lib/types";
import { PhotoShowcase } from "./PhotoShowcase";
import { VideoGallery } from "./VideoGallery";

type TabKey = "photos" | "films";

export function CategoryShowcase({ photos, videos }: { photos: PortfolioMedia[]; videos: PortfolioMedia[] }) {
  const [tab, setTab] = useState<TabKey>("photos");
  const photoItems = useMemo(() => photos.filter((p) => p.media_type === "photo" && p.url), [photos]);
  const videoItems = useMemo(() => videos.filter((v) => v.media_type === "video" && v.url), [videos]);

  return (
    <div className="space-y-10">
      <div className="sticky top-20 z-20 -mx-4 flex flex-wrap items-center justify-between gap-4 border-y border-gold/10 bg-bg-deep/70 px-4 py-3 backdrop-blur md:mx-0 md:rounded-sm md:border md:bg-bg-card">
        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={() => setTab("photos")}
            className={`rounded-sm px-4 py-2 text-xs tracking-widest uppercase transition ${
              tab === "photos" ? "bg-gold/10 text-gold" : "text-text-muted hover:text-gold"
            }`}
          >
            Photos ({photoItems.length})
          </button>
          <button
            type="button"
            onClick={() => setTab("films")}
            disabled={videoItems.length === 0}
            className={`rounded-sm px-4 py-2 text-xs tracking-widest uppercase transition ${
              tab === "films" ? "bg-gold/10 text-gold" : "text-text-muted hover:text-gold"
            } ${videoItems.length === 0 ? "opacity-40" : ""}`}
          >
            Films ({videoItems.length})
          </button>
        </div>
        <p className="text-xs text-text-muted">
          Tap any photo or film to view full screen.
        </p>
      </div>

      {tab === "photos" ? (
        <PhotoShowcase items={photoItems} />
      ) : (
        <section className="space-y-6">
          <div>
            <h2 className="font-display text-2xl text-gold">Cinematic films</h2>
            <p className="mt-2 text-sm text-text-muted">Play a film, then use arrows to continue.</p>
          </div>
          <VideoGallery items={videoItems} />
        </section>
      )}
    </div>
  );
}
