"use client";

import { useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import Image from "next/image";
import type { PortfolioMedia } from "@/lib/types";

interface LightboxProps {
  items: PortfolioMedia[];
  index: number | null;
  onClose: () => void;
  onPrev: () => void;
  onNext: () => void;
  onSelect?: (index: number) => void;
}

export function Lightbox({ items, index, onClose, onPrev, onNext, onSelect }: LightboxProps) {
  const item = index !== null ? items[index] : null;

  useEffect(() => {
    if (!item) return;
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
      if (e.key === "ArrowLeft") onPrev();
      if (e.key === "ArrowRight") onNext();
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [item, onClose, onPrev, onNext]);

  return (
    <AnimatePresence>
      {item && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          className="fixed inset-0 z-[100] flex items-center justify-center bg-black/95 p-4"
          onClick={onClose}
        >
          <button
            type="button"
            className="absolute right-4 top-4 z-10 text-2xl text-gold"
            onClick={onClose}
            aria-label="Close"
          >
            ×
          </button>
          {index !== null && (
            <div className="absolute left-4 top-4 z-10 rounded-full border border-gold/20 bg-bg-deep/70 px-4 py-2 text-xs tracking-widest text-gold backdrop-blur">
              {index + 1} / {items.length}
            </div>
          )}
          {items.length > 1 && (
            <>
              <button
                type="button"
                className="absolute left-4 z-10 rounded-full border border-gold/30 bg-bg-deep/60 px-4 py-2 text-gold backdrop-blur hover:bg-gold/10"
                onClick={(e) => { e.stopPropagation(); onPrev(); }}
                aria-label="Previous"
              >
                ‹
              </button>
              <button
                type="button"
                className="absolute right-16 z-10 rounded-full border border-gold/30 bg-bg-deep/60 px-4 py-2 text-gold backdrop-blur hover:bg-gold/10"
                onClick={(e) => { e.stopPropagation(); onNext(); }}
                aria-label="Next"
              >
                ›
              </button>
            </>
          )}
          <motion.div
            initial={{ scale: 0.95, opacity: 0 }}
            animate={{ scale: 1, opacity: 1 }}
            exit={{ scale: 0.95, opacity: 0 }}
            className="relative max-h-[90vh] w-full max-w-6xl"
            onClick={(e) => e.stopPropagation()}
          >
            {item.media_type === "video" ? (
              <video src={item.url} controls className="max-h-[85vh] w-full" />
            ) : (
              <Image
                src={item.url || "/placeholder.jpg"}
                alt={item.title || "Portfolio"}
                width={1400}
                height={900}
                className="max-h-[85vh] w-auto object-contain"
                unoptimized
              />
            )}
            {(item.title || item.caption) && (
              <div className="mx-auto mt-4 max-w-3xl text-center">
                {item.title && <p className="font-display text-xl text-gold">{item.title}</p>}
                {item.caption && <p className="mt-2 text-sm text-text-muted">{item.caption}</p>}
              </div>
            )}

            {onSelect && items.length > 1 && (
              <div className="mt-6 overflow-x-auto pb-2">
                <div className="mx-auto flex w-max gap-2 px-2">
                  {items.map((it, i) => {
                    const thumb = it.thumbnail_url || it.url;
                    const active = i === index;
                    return (
                      <button
                        key={it.id}
                        type="button"
                        onClick={() => onSelect(i)}
                        className={`h-16 w-16 flex-shrink-0 overflow-hidden rounded-sm border-2 transition ${
                          active ? "border-gold" : "border-gold/20 opacity-70 hover:opacity-100"
                        }`}
                        aria-label={`Open item ${i + 1}`}
                      >
                        {thumb ? (
                          <img src={thumb} alt="" className="h-full w-full object-cover" loading="lazy" />
                        ) : (
                          <div className="h-full w-full bg-bg-muted" />
                        )}
                      </button>
                    );
                  })}
                </div>
              </div>
            )}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
