import { serverFetch } from "@/lib/server-api";
import type { ContactInfo } from "@/lib/types";
import { PageTransition } from "@/components/PageTransition";
import { DEFAULT_CONTACT, mergeContact } from "@/lib/constants";
import { ContactStrip } from "@/components/layout/ContactStrip";

export const metadata = { title: "Contact" };

async function getContactInfo(): Promise<ContactInfo> {
  return mergeContact(await serverFetch<ContactInfo>("/api/contact-info", 600));
}

export default async function ContactPage() {
  const info = await getContactInfo();
  const wa = info.whatsapp ? `https://wa.me/${info.whatsapp.replace(/\D/g, "")}` : "#";

  return (
    <PageTransition>
      <section className="flex min-h-[40vh] items-end bg-bg-card pb-12 pt-32">
        <div className="section-pad mx-auto w-full max-w-7xl">
          <p className="text-xs tracking-[0.4em] text-gold uppercase">We&apos;d Love to Hear From You</p>
          <h1 className="mt-4 font-display text-5xl">Connect With Our Studio</h1>
          <p className="mt-4 max-w-xl text-text-muted">
            Reach us by phone or WhatsApp — we&apos;ll guide you with care for your celebration.
          </p>
        </div>
      </section>

      <section className="section-pad mx-auto max-w-2xl">
        <div className="space-y-8">
          <h2 className="font-display text-2xl text-gold">Call or message us</h2>
          <div className="flex flex-wrap gap-4">
            <a href={wa} target="_blank" rel="noopener noreferrer" className="btn-primary">
              WhatsApp Us
            </a>
            <a href={`tel:${info.phone_number || DEFAULT_CONTACT.phone_number}`} className="btn-outline">
              Call {info.phone_number || DEFAULT_CONTACT.phone_number}
            </a>
          </div>
          <ContactStrip contact={info} />
          {info.google_maps_embed && (
            <div className="aspect-video overflow-hidden rounded-sm border border-gold/10">
              <iframe
                src={info.google_maps_embed}
                className="h-full w-full border-0"
                loading="lazy"
                referrerPolicy="no-referrer-when-downgrade"
                title="Studio location"
              />
            </div>
          )}
        </div>
      </section>
    </PageTransition>
  );
}
