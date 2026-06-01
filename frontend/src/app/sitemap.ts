import type { MetadataRoute } from "next";

const base = process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:3000";

const slugs = [
  "wedding", "baby-shower", "puberty", "birthday", "ear-piercing",
  "couple-shoots", "maternity", "outdoor", "temple-functions",
];

export default function sitemap(): MetadataRoute.Sitemap {
  const routes = ["", "/portfolio", "/about", "/contact", ...slugs.map((s) => `/portfolio/${s}`)];
  return routes.map((path) => ({
    url: `${base}${path}`,
    lastModified: new Date(),
    changeFrequency: path === "" ? "weekly" : "monthly",
    priority: path === "" ? 1 : 0.8,
  }));
}
