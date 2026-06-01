"use client";

import Image from "next/image";
import Link from "next/link";
import { useSiteBranding } from "@/components/SiteBrandingProvider";

type LogoProps = {
  className?: string;
  href?: string;
  /** Admin shell passes logo directly (no site provider) */
  logoUrl?: string;
  compact?: boolean;
};

export function Logo({ className = "", href = "/", logoUrl: logoUrlProp, compact = false }: LogoProps) {
  const { logoUrl: ctxLogo } = useSiteBranding();
  const logoUrl = logoUrlProp ?? ctxLogo;

  return (
    <Link href={href} className={`group flex items-center gap-3 ${className}`}>
      {logoUrl ? (
        <span className="relative flex h-10 w-10 flex-shrink-0 items-center justify-center overflow-hidden rounded-full border border-gold/40 bg-bg-card">
          <Image src={logoUrl} alt="Surya Photography" fill className="object-contain p-1" unoptimized />
        </span>
      ) : (
        <span className="flex h-10 w-10 flex-shrink-0 items-center justify-center rounded-full border border-gold/40 bg-bg-card text-lg font-display font-semibold text-gold transition group-hover:border-gold">
          S
        </span>
      )}
      {!compact && (
        <span className="flex flex-col leading-tight">
          <span className="font-display text-xl tracking-[0.2em] text-text-main uppercase">Surya</span>
          <span className="text-[10px] tracking-[0.35em] text-gold uppercase">Photography</span>
        </span>
      )}
    </Link>
  );
}
