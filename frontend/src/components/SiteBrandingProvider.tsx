"use client";

import { createContext, useContext } from "react";

type Branding = {
  logoUrl?: string;
};

const SiteBrandingContext = createContext<Branding>({});

export function SiteBrandingProvider({
  logoUrl,
  children,
}: {
  logoUrl?: string;
  children: React.ReactNode;
}) {
  return (
    <SiteBrandingContext.Provider value={{ logoUrl }}>{children}</SiteBrandingContext.Provider>
  );
}

export function useSiteBranding() {
  return useContext(SiteBrandingContext);
}
