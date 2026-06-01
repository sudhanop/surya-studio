"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, clearToken, getToken } from "@/lib/api";
import { Logo } from "@/components/ui/Logo";
import type { PublicSiteData } from "@/lib/types";
import { DEFAULT_SITE_DATA } from "@/lib/constants";

const nav = [
  { href: "/admin/dashboard", label: "Dashboard" },
  { href: "/admin/inquiries", label: "Inquiries" },
  { href: "/admin/functions", label: "Functions" },
  { href: "/admin/portfolio", label: "Portfolio" },
  { href: "/admin/settings", label: "Studio" },
  { href: "/admin/account", label: "Account" },
];

export function AdminShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [ready, setReady] = useState(false);
  const [logoUrl, setLogoUrl] = useState<string | undefined>();

  useEffect(() => {
    if (pathname === "/admin/login") {
      setReady(true);
      return;
    }
    const token = getToken();
    if (!token) {
      router.replace("/admin/login");
      return;
    }
    Promise.all([
      api.getAuth("/api/admin/me"),
      api.getAuth<PublicSiteData>("/api/admin/settings").catch(() => DEFAULT_SITE_DATA),
    ])
      .then(([, site]) => {
        setLogoUrl(site.logo_url);
        setReady(true);
      })
      .catch(() => {
        clearToken();
        router.replace("/admin/login");
      });
  }, [pathname, router]);

  if (pathname === "/admin/login") return <>{children}</>;
  if (!ready) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-bg-deep text-gold">
        Loading...
      </div>
    );
  }

  return (
    <div className="flex min-h-screen bg-bg-deep">
      <aside className="hidden w-64 flex-shrink-0 border-r border-gold/10 bg-bg-card p-6 md:block">
        <Logo logoUrl={logoUrl} href="/admin/dashboard" />
        <nav className="mt-10 space-y-2">
          {nav.map((n) => (
            <Link
              key={n.href}
              href={n.href}
              className={`block rounded-sm px-4 py-2 text-sm tracking-wide ${
                pathname === n.href ? "bg-gold/10 text-gold" : "text-text-muted hover:text-gold"
              }`}
            >
              {n.label}
            </Link>
          ))}
        </nav>
        <button
          type="button"
          onClick={async () => {
            await api.postAuth("/api/admin/logout").catch(() => {});
            clearToken();
            router.push("/admin/login");
          }}
          className="mt-8 text-xs tracking-widest text-text-muted uppercase hover:text-gold"
        >
          Logout
        </button>
      </aside>
      <main className="flex-1 overflow-auto p-6 pb-24 md:p-10 md:pb-10">{children}</main>

      <nav className="fixed bottom-0 left-0 right-0 flex border-t border-gold/10 bg-bg-card md:hidden">
        {nav.map((n) => (
          <Link
            key={n.href}
            href={n.href}
            className={`flex-1 py-3 text-center text-[10px] tracking-wider uppercase ${
              pathname === n.href ? "text-gold" : "text-text-muted"
            }`}
          >
            {n.label}
          </Link>
        ))}
      </nav>
    </div>
  );
}
