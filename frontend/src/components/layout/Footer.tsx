import Link from "next/link";
import { Logo } from "@/components/ui/Logo";
import { HiddenAdminAccess } from "@/components/layout/HiddenAdminAccess";
import { serverFetch } from "@/lib/server-api";
import type { ContactInfo } from "@/lib/types";
import { DEFAULT_CONTACT, mergeContact, SITE } from "@/lib/constants";

async function getContact(): Promise<ContactInfo> {
  const raw = await serverFetch<ContactInfo>("/api/contact-info", 300);
  return mergeContact(raw);
}

export async function Footer() {
  const year = new Date().getFullYear();
  const c = await getContact();

  return (
    <footer className="border-t border-gold/10 bg-bg-card">
      <div className="section-pad mx-auto max-w-7xl">
        <div className="grid gap-12 md:grid-cols-3">
          <div>
            <HiddenAdminAccess>
              <Logo />
            </HiddenAdminAccess>
            <p className="mt-4 max-w-xs text-sm leading-relaxed text-text-muted">
              {SITE.tagline} Every celebration deserves to be remembered with heart.
            </p>
          </div>
          <div>
            <h4 className="mb-4 text-xs tracking-[0.3em] text-gold uppercase">Explore</h4>
            <ul className="space-y-2 text-sm text-text-muted">
              <li>
                <Link href="/portfolio" className="hover:text-gold">
                  Our Portfolio
                </Link>
              </li>
              <li>
                <Link href="/about" className="hover:text-gold">
                  Our Story
                </Link>
              </li>
              <li>
                <Link href="/contact" className="hover:text-gold">
                  Reach Out
                </Link>
              </li>
            </ul>
          </div>
          <div>
            <h4 className="mb-4 text-xs tracking-[0.3em] text-gold uppercase">Studio</h4>
            <p className="text-sm leading-relaxed text-text-muted">
              {c.address}
              {c.pincode ? ` — ${c.pincode}` : ""}
            </p>
            <a href={`mailto:${c.contact_email}`} className="mt-2 block text-sm text-gold hover:underline">
              {c.contact_email}
            </a>
            <a href={`tel:${c.phone_number}`} className="mt-1 block text-sm text-text-muted hover:text-gold">
              {c.phone_number} · {c.phone_secondary}
            </a>
          </div>
        </div>
        <div className="mt-12 border-t border-gold/10 pt-8 text-center text-xs text-text-muted">
          <HiddenAdminAccess>
            <span>© {year} {SITE.name}. Crafted with devotion.</span>
          </HiddenAdminAccess>
        </div>
      </div>
    </footer>
  );
}
