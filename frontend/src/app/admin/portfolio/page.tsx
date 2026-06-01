"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import { api, asArray } from "@/lib/api";
import type { Category, PortfolioMedia } from "@/lib/types";

export default function AdminPortfolioPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [categoryId, setCategoryId] = useState<number>(0);
  const [media, setMedia] = useState<PortfolioMedia[]>([]);
  const [uploading, setUploading] = useState(false);

  useEffect(() => {
    api.getAuth<Category[]>("/api/admin/categories").then((cats) => {
      const list = asArray(cats);
      setCategories(list);
      if (list[0]) setCategoryId(list[0].id);
    });
  }, []);

  useEffect(() => {
    if (!categoryId) return;
    api
      .getAuth<PortfolioMedia[]>(`/api/admin/media?category_id=${categoryId}`)
      .then((data) => setMedia(asArray(data)));
  }, [categoryId]);

  async function handleUpload(e: React.ChangeEvent<HTMLInputElement>) {
    const files = Array.from(e.target.files ?? []);
    if (!files.length || !categoryId) return;
    setUploading(true);
    try {
      for (const file of files) {
        const fd = new FormData();
        fd.append("file", file);
        fd.append("category_id", String(categoryId));
        fd.append("media_type", file.type.startsWith("video") ? "video" : "photo");
        await api.uploadAuth("/api/admin/upload", fd);
      }
      const list = await api.getAuth<PortfolioMedia[]>(`/api/admin/media?category_id=${categoryId}`);
      setMedia(asArray(list));
    } finally {
      setUploading(false);
      e.target.value = "";
    }
  }

  async function handleDelete(m: PortfolioMedia) {
    if (!confirm("Delete this media?")) return;
    await api.deleteAuth(`/api/admin/media/${m.id}`);
    setMedia((prev) => prev.filter((x) => x.id !== m.id));
  }

  return (
    <div>
      <h1 className="font-display text-3xl text-gold">Portfolio Management</h1>

      <div className="mt-6 flex flex-wrap gap-4">
        <select
          value={categoryId}
          onChange={(e) => setCategoryId(Number(e.target.value))}
          className="border border-gold/20 bg-bg-muted px-4 py-2 text-sm"
        >
          {categories.map((c) => (
            <option key={c.id} value={c.id}>{c.name}</option>
          ))}
        </select>
        <label className="btn-outline cursor-pointer">
          {uploading ? "Uploading..." : "Upload Media"}
          <input type="file" accept="image/*,video/*" multiple className="hidden" onChange={handleUpload} disabled={uploading} />
        </label>
      </div>

      <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {media.map((m) => (
          <div key={m.id} className="overflow-hidden rounded-sm border border-gold/10 bg-bg-card">
            <div className="relative aspect-square bg-bg-muted">
              {m.media_type === "photo" && m.url ? (
                <Image src={m.url} alt="" fill className="object-cover" unoptimized />
              ) : (
                <video src={m.url} className="h-full w-full object-cover" />
              )}
            </div>
            <div className="flex items-center justify-end p-3 text-xs">
              <button
                type="button"
                className="text-red-400"
                onClick={() => handleDelete(m)}
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
