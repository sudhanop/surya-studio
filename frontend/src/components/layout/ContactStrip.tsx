import Link from "next/link";
import { DEFAULT_CONTACT } from "@/lib/constants";
import type { ContactInfo } from "@/lib/types";

export function ContactStrip({ contact }: { contact: ContactInfo }) {
  const c = { ...DEFAULT_CONTACT, ...contact };
  const wa = c.whatsapp ? `https://wa.me/${c.whatsapp.replace(/\D/g, "")}` : null;

  return (
    <div className="space-y-2 text-sm text-text-muted">
      <p>
        <span className="text-gold">Email:</span>{" "}
        <a href={`mailto:${c.contact_email}`} className="hover:text-gold hover:underline">
          {c.contact_email}
        </a>
      </p>
      <p>
        <span className="text-gold">Call:</span>{" "}
        <a href={`tel:${c.phone_number}`} className="hover:text-gold">
          {c.phone_number}
        </a>
        {c.phone_secondary && (
          <>
            {" · "}
            <a href={`tel:${c.phone_secondary}`} className="hover:text-gold">
              {c.phone_secondary}
            </a>
          </>
        )}
      </p>
      <p>
        <span className="text-gold">Instagram:</span>{" "}
        <a href={c.instagram_url} target="_blank" rel="noopener noreferrer" className="hover:text-gold hover:underline">
          @{DEFAULT_CONTACT.instagram_handle}
        </a>
      </p>
      <p>
        <span className="text-gold">YouTube:</span>{" "}
        <a href={c.youtube_url} target="_blank" rel="noopener noreferrer" className="hover:text-gold hover:underline">
          @{DEFAULT_CONTACT.youtube_handle}
        </a>
      </p>
      <p className="leading-relaxed">
        <span className="text-gold">Studio:</span> {c.address}
        {c.pincode ? ` — ${c.pincode}` : ""}
      </p>
      {wa && (
        <Link href={wa} target="_blank" rel="noopener noreferrer" className="btn-outline mt-4 inline-block text-xs">
          Message on WhatsApp
        </Link>
      )}
    </div>
  );
}
