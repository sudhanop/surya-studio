import { serverFetch } from "@/lib/server-api";
import type { CategoryMediaResponse } from "@/lib/types";
import { CategoryShowcase } from "@/components/gallery/CategoryShowcase";
import { PageTransition } from "@/components/PageTransition";
import { notFound } from "next/navigation";

export async function generateMetadata({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await serverFetch<CategoryMediaResponse>(`/api/categories/${slug}/media`, 120);
  return { title: data?.category.name ?? "Portfolio" };
}

export default async function CategoryPage({ params }: { params: Promise<{ slug: string }> }) {
  const { slug } = await params;
  const data = await serverFetch<CategoryMediaResponse>(`/api/categories/${slug}/media`, 120);
  if (!data) notFound();

  const { category, photos, videos, has_videos } = data;

  return (
    <PageTransition>
      <section className="flex min-h-[40vh] items-end bg-bg-card pb-12 pt-32">
        <div className="section-pad mx-auto w-full max-w-7xl">
          <p className="text-xs tracking-[0.4em] text-gold uppercase">Collection</p>
          <h1 className="mt-4 font-display text-5xl">{category.name}</h1>
          {category.description && (
            <p className="mt-4 max-w-xl text-text-muted">{category.description}</p>
          )}
        </div>
      </section>

      <section className="section-pad mx-auto max-w-7xl">
        <CategoryShowcase photos={photos} videos={has_videos ? videos : []} />
      </section>
    </PageTransition>
  );
}
