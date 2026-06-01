import type { Metadata } from "next";
import { Navbar } from "@/components/layout/Navbar";
import { Footer } from "@/components/layout/Footer";
import { SiteBrandingProvider } from "@/components/SiteBrandingProvider";
import { serverFetch } from "@/lib/server-api";
import type { PublicSiteData } from "@/lib/types";
import { DEFAULT_SITE_DATA } from "@/lib/constants";

async function getSite() {
  return (await serverFetch<PublicSiteData>("/api/site", 120)) ?? DEFAULT_SITE_DATA;
}

export async function generateMetadata(): Promise<Metadata> {
  const site = await getSite();
  if (site.logo_url) {
    return {
      icons: {
        icon: site.logo_url,
        apple: site.logo_url,
      },
    };
  }
  return {};
}

export default async function SiteLayout({ children }: { children: React.ReactNode }) {
  const site = await getSite();

  return (
    <SiteBrandingProvider logoUrl={site.logo_url}>
      <Navbar />
      <main className="min-h-screen pt-0">{children}</main>
      <Footer />
    </SiteBrandingProvider>
  );
}
