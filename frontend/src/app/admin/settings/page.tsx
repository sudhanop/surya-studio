"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import { api } from "@/lib/api";
import type { PublicSiteData, Testimonial } from "@/lib/types";
import { DEFAULT_SITE_DATA, DEFAULT_TESTIMONIALS } from "@/lib/constants";

const inputClass =
  "w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60";

export default function AdminSettingsPage() {
  const [data, setData] = useState<PublicSiteData | null>(null);
  const [testimonials, setTestimonials] = useState<Testimonial[]>([...DEFAULT_TESTIMONIALS]);
  const [saving, setSaving] = useState(false);
  const [uploadingPortrait, setUploadingPortrait] = useState(false);
  const [uploadingLogo, setUploadingLogo] = useState(false);

  const load = () =>
    api
      .getAuth<PublicSiteData>("/api/admin/settings")
      .then((site) => {
        setData(site);
        setTestimonials(site.testimonials?.length ? site.testimonials : [...DEFAULT_TESTIMONIALS]);
      })
      .catch(() => {
        setData(DEFAULT_SITE_DATA);
        setTestimonials([...DEFAULT_TESTIMONIALS]);
      });

  useEffect(() => {
    load();
  }, []);

  async function handleSave(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (!data) return;
    setSaving(true);
    const fd = new FormData(e.currentTarget);
    try {
      await api.putAuth("/api/admin/settings", {
        events_covered: fd.get("events_covered"),
        years_of_craft: fd.get("years_of_craft"),
        happy_families: fd.get("happy_families"),
        contact_email: fd.get("contact_email"),
        phone_primary: fd.get("phone_primary"),
        phone_secondary: fd.get("phone_secondary"),
        instagram_url: fd.get("instagram_url"),
        youtube_url: fd.get("youtube_url"),
        address: fd.get("address"),
        pincode: fd.get("pincode"),
        whatsapp: fd.get("whatsapp"),
        testimonials: testimonials.filter((t) => t.name.trim() && t.text.trim()),
      });
      load();
    } finally {
      setSaving(false);
    }
  }

  async function handlePortrait(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploadingPortrait(true);
    const fd = new FormData();
    fd.append("file", file);
    try {
      await api.uploadAuth("/api/admin/settings/owner-portrait", fd);
      load();
    } finally {
      setUploadingPortrait(false);
      e.target.value = "";
    }
  }

  async function handleLogo(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;
    setUploadingLogo(true);
    const fd = new FormData();
    fd.append("file", file);
    try {
      await api.uploadAuth("/api/admin/settings/logo", fd);
      load();
    } finally {
      setUploadingLogo(false);
      e.target.value = "";
    }
  }

  async function removePortrait() {
    if (!confirm("Remove the owner portrait from the About page?")) return;
    await api.deleteAuth("/api/admin/settings/owner-portrait");
    load();
  }

  async function removeLogo() {
    if (!confirm("Remove custom logo and use the default mark?")) return;
    await api.deleteAuth("/api/admin/settings/logo");
    load();
  }

  function updateTestimonial(i: number, field: keyof Testimonial, value: string) {
    setTestimonials((prev) => {
      const next = [...prev];
      next[i] = { ...next[i], [field]: value };
      return next;
    });
  }

  function addTestimonial() {
    setTestimonials((prev) => [...prev, { name: "", text: "" }]);
  }

  function removeTestimonial(i: number) {
    setTestimonials((prev) => prev.filter((_, j) => j !== i));
  }

  if (!data) {
    return <p className="text-gold">Loading studio settings…</p>;
  }

  const c = data.contact;

  return (
    <div className="max-w-2xl">
      <h1 className="font-display text-3xl text-gold">Studio Story & Branding</h1>
      <p className="mt-2 text-sm text-text-muted">
        Logo, kind words, contact details, and your portrait — everything families see on the website.
      </p>

      <section className="mt-10 glass rounded-sm p-6">
        <h2 className="font-display text-xl text-gold">Studio logo</h2>
        <p className="mt-2 text-sm text-text-muted">
          Shown in the navigation, admin sidebar, and browser tab icon after upload.
        </p>
        <div className="mt-6 flex flex-wrap items-start gap-6">
          <div className="relative flex h-20 w-20 items-center justify-center overflow-hidden rounded-full border border-gold/20 bg-bg-muted">
            {data.logo_url ? (
              <Image src={data.logo_url} alt="Studio logo" fill className="object-contain p-2" unoptimized />
            ) : (
              <span className="font-display text-2xl text-gold">S</span>
            )}
          </div>
          <div className="space-y-2">
            <label className="btn-outline inline-block cursor-pointer text-xs">
              {uploadingLogo ? "Uploading…" : "Upload logo"}
              <input type="file" accept="image/png" className="hidden" onChange={handleLogo} disabled={uploadingLogo} />
            </label>
            {data.logo_url && (
              <button type="button" onClick={removeLogo} className="block text-xs text-red-400 hover:underline">
                Remove logo
              </button>
            )}
          </div>
        </div>
      </section>

      <section className="mt-10 glass rounded-sm p-6">
        <h2 className="font-display text-xl text-gold">Owner portrait</h2>
        <p className="mt-2 text-sm text-text-muted">Shown on the About page.</p>
        <div className="mt-6 flex flex-wrap items-start gap-6">
          <div className="relative h-40 w-40 overflow-hidden rounded-sm border border-gold/20 bg-bg-muted">
            {data.owner_portrait_url ? (
              <Image src={data.owner_portrait_url} alt="Owner" fill className="object-cover" unoptimized />
            ) : (
              <div className="flex h-full items-center justify-center text-xs text-text-muted">No photo yet</div>
            )}
          </div>
          <div className="space-y-2">
            <label className="btn-outline inline-block cursor-pointer text-xs">
              {uploadingPortrait ? "Uploading…" : "Upload portrait"}
              <input
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handlePortrait}
                disabled={uploadingPortrait}
              />
            </label>
            {data.owner_portrait_url && (
              <button type="button" onClick={removePortrait} className="block text-xs text-red-400 hover:underline">
                Remove photo
              </button>
            )}
          </div>
        </div>
      </section>

      <section className="mt-10 glass rounded-sm p-6">
        <div className="flex items-center justify-between">
          <h2 className="font-display text-xl text-gold">Kind words</h2>
          <button type="button" onClick={addTestimonial} className="text-xs text-gold hover:underline">
            + Add review
          </button>
        </div>
        <p className="mt-2 text-sm text-text-muted">
          Couple or family names appear in gold on the home page, just as they do today.
        </p>
        <div className="mt-6 space-y-6">
          {testimonials.map((t, i) => (
            <div key={i} className="rounded-sm border border-gold/15 p-4">
              <label className="text-xs text-text-muted">
                Name (e.g. Priya & Arun)
                <input
                  value={t.name}
                  onChange={(e) => updateTestimonial(i, "name", e.target.value)}
                  className={`${inputClass} mt-1`}
                />
              </label>
              <label className="mt-3 block text-xs text-text-muted">
                Their words
                <textarea
                  value={t.text}
                  onChange={(e) => updateTestimonial(i, "text", e.target.value)}
                  rows={3}
                  className={`${inputClass} mt-1`}
                />
              </label>
              {testimonials.length > 1 && (
                <button
                  type="button"
                  onClick={() => removeTestimonial(i)}
                  className="mt-2 text-xs text-red-400 hover:underline"
                >
                  Remove
                </button>
              )}
            </div>
          ))}
        </div>
      </section>

      <form onSubmit={handleSave} className="mt-10 space-y-8">
        <section className="glass rounded-sm p-6">
          <h2 className="font-display text-xl text-gold">Hearts we&apos;ve touched</h2>
          <div className="mt-4 grid gap-4 sm:grid-cols-3">
            <label className="text-xs text-text-muted">
              Events covered
              <input name="events_covered" defaultValue={data.events_covered} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Years of craft
              <input name="years_of_craft" defaultValue={data.years_of_craft} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Happy families
              <input name="happy_families" defaultValue={data.happy_families} className={`${inputClass} mt-1`} />
            </label>
          </div>
        </section>

        <section className="glass rounded-sm p-6">
          <h2 className="font-display text-xl text-gold">Contact details</h2>
          <div className="mt-4 grid gap-4 sm:grid-cols-2">
            <label className="text-xs text-text-muted sm:col-span-2">
              Email (display only — no online enquiry form)
              <input name="contact_email" defaultValue={c.contact_email} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Primary mobile
              <input name="phone_primary" defaultValue={c.phone_number} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Secondary mobile
              <input name="phone_secondary" defaultValue={c.phone_secondary || ""} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              WhatsApp (with country code)
              <input name="whatsapp" defaultValue={c.whatsapp} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Instagram URL
              <input name="instagram_url" defaultValue={c.instagram_url} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted sm:col-span-2">
              YouTube URL
              <input name="youtube_url" defaultValue={c.youtube_url} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted sm:col-span-2">
              Studio address
              <textarea name="address" rows={3} defaultValue={c.address || ""} className={`${inputClass} mt-1`} />
            </label>
            <label className="text-xs text-text-muted">
              Pincode
              <input name="pincode" defaultValue={c.pincode || ""} className={`${inputClass} mt-1`} />
            </label>
          </div>
        </section>

        <button type="submit" disabled={saving} className="btn-primary">
          {saving ? "Saving…" : "Save all changes"}
        </button>
      </form>
    </div>
  );
}
