"use client";

import { useCallback, useEffect, useState } from "react";
import { api, asArray } from "@/lib/api";
import type { Inquiry } from "@/lib/types";

const statuses = [
  { key: "pending", label: "Pending" },
  { key: "spoken", label: "Spoken" },
  { key: "converted", label: "Converted" },
  { key: "cancelled", label: "Cancelled" },
  { key: "all", label: "All" },
] as const;

export default function AdminInquiriesPage() {
  const [items, setItems] = useState<Inquiry[]>([]);
  const [status, setStatus] = useState<(typeof statuses)[number]["key"]>("pending");
  const [saving, setSaving] = useState(false);

  const load = useCallback((s: (typeof statuses)[number]["key"]) => {
    const qs = s === "all" ? "" : `?status=${encodeURIComponent(s)}`;
    return api
      .getAuth<Inquiry[]>(`/api/admin/inquiries${qs}`)
      .then((data) => setItems(asArray(data)))
      .catch(console.error);
  }, []);

  useEffect(() => {
    load(status);
  }, [status, load]);

  async function updateStatus(id: number, next: string) {
    await api.putAuth(`/api/admin/inquiries/${id}/status`, { status: next });
    await load(status);
  }

  async function convert(id: number) {
    await api.postAuth(`/api/admin/inquiries/${id}/convert`, {});
    await load(status);
  }

  async function remove(id: number) {
    if (!confirm("Delete this inquiry?")) return;
    await api.deleteAuth(`/api/admin/inquiries/${id}`);
    await load(status);
  }

  async function create(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setSaving(true);
    const fd = new FormData(e.currentTarget);
    try {
      await api.postAuth("/api/admin/inquiries", {
        customer_name: fd.get("customer_name"),
        phone_number: fd.get("phone_number"),
        occasion_type: fd.get("occasion_type"),
        wanted_date: fd.get("wanted_date") || "",
        address: fd.get("address") || "",
        message: fd.get("message") || "",
        status: "pending",
      });
      e.currentTarget.reset();
      await load("pending");
      setStatus("pending");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div>
      <h1 className="font-display text-3xl text-gold">Inquiries</h1>

      <section className="mt-6 glass rounded-sm p-6">
        <h2 className="font-display text-xl text-gold">Add inquiry</h2>
        <form onSubmit={create} className="mt-4 grid gap-4 md:grid-cols-2">
          <label className="text-xs text-text-muted">
            Customer name
            <input
              name="customer_name"
              required
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <label className="text-xs text-text-muted">
            Phone number
            <input
              name="phone_number"
              required
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <label className="text-xs text-text-muted">
            Occasion
            <input
              name="occasion_type"
              required
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <label className="text-xs text-text-muted">
            Wanted date
            <input
              name="wanted_date"
              type="date"
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <label className="text-xs text-text-muted md:col-span-2">
            Address
            <input
              name="address"
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <label className="text-xs text-text-muted md:col-span-2">
            Notes
            <textarea
              name="message"
              rows={2}
              className="mt-1 w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60"
            />
          </label>
          <div className="md:col-span-2">
            <button type="submit" disabled={saving} className="btn-primary">
              {saving ? "Saving..." : "Add inquiry"}
            </button>
          </div>
        </form>
      </section>

      <div className="mt-8 flex flex-wrap gap-2">
        {statuses.map((s) => (
          <button
            key={s.key}
            type="button"
            onClick={() => setStatus(s.key)}
            className={`rounded-sm border px-4 py-2 text-xs tracking-widest uppercase ${
              status === s.key
                ? "border-gold bg-gold/10 text-gold"
                : "border-gold/20 text-text-muted hover:border-gold/50 hover:text-gold"
            }`}
          >
            {s.label}
          </button>
        ))}
      </div>

      <div className="mt-8 overflow-x-auto">
        <table className="w-full text-left text-sm">
          <thead className="border-b border-gold/20 text-xs tracking-widest text-gold uppercase">
            <tr>
              <th className="p-3">Name</th>
              <th className="p-3">Phone</th>
              <th className="p-3">Occasion</th>
              <th className="p-3">Date</th>
              <th className="p-3">Status</th>
              <th className="p-3">Actions</th>
            </tr>
          </thead>
          <tbody>
            {items.map((i) => (
              <tr key={i.id} className="border-b border-gold/5">
                <td className="p-3">{i.customer_name}</td>
                <td className="p-3">{i.phone_number}</td>
                <td className="p-3">{i.occasion_type}</td>
                <td className="p-3">{i.wanted_date?.slice(0, 10) || "—"}</td>
                <td className="p-3 capitalize">{i.status}</td>
                <td className="p-3 space-x-2">
                  {(i.status === "pending" || i.status === "new") && (
                    <button
                      type="button"
                      className="text-gold hover:underline"
                      onClick={() => updateStatus(i.id, "spoken")}
                    >
                      Mark spoken
                    </button>
                  )}
                  {(i.status === "spoken" || i.status === "pending" || i.status === "new") && (
                    <button
                      type="button"
                      className="text-gold hover:underline"
                      onClick={() => updateStatus(i.id, "cancelled")}
                    >
                      Cancel
                    </button>
                  )}
                  {i.status === "cancelled" && (
                    <button
                      type="button"
                      className="text-gold hover:underline"
                      onClick={() => updateStatus(i.id, "pending")}
                    >
                      Reopen
                    </button>
                  )}
                  {i.status !== "converted" && i.status !== "cancelled" && (
                    <button
                      type="button"
                      className="text-gold hover:underline"
                      onClick={() => convert(i.id)}
                    >
                      Convert to order
                    </button>
                  )}
                  <button
                    type="button"
                    className="text-red-400 hover:underline"
                    onClick={() => remove(i.id)}
                  >
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
