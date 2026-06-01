"use client";

import { useEffect, useState } from "react";
import { api, asArray } from "@/lib/api";
import type { StudioFunction } from "@/lib/types";
import {
  ALBUM_STATUSES,
  BOOKING_SERVICES,
  COMPLIMENTARY_ITEMS,
  FUNCTION_STATUSES,
  VIDEO_STATUSES,
} from "@/lib/constants";
import { formatEventDates, functionPayload } from "@/lib/function-utils";

const inputClass =
  "w-full border border-gold/25 bg-bg-muted px-3 py-2 text-sm text-text-main outline-none focus:border-gold/60";

function statusLabel(value: string) {
  return value.replaceAll("_", " ");
}

function emptyForm() {
  return {
    customer_name: "",
    phone_number: "",
    function_type: "",
    event_dates: [""] as string[],
    total_amount: 0,
    advance_paid: 0,
    customer_booking_notes: "",
    services: [] as string[],
    complimentary: [] as string[],
    album_status: "not_started",
    video_status: "not_started",
    overall_status: "upcoming",
  };
}

export default function AdminFunctionsPage() {
  const [items, setItems] = useState<StudioFunction[]>([]);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState(emptyForm());
  const [expandedId, setExpandedId] = useState<number | null>(null);

  const load = () =>
    api.getAuth<StudioFunction[]>("/api/admin/functions").then((data) => setItems(asArray(data))).catch(console.error);

  useEffect(() => {
    load();
  }, []);

  function toggleList(key: "services" | "complimentary", value: string) {
    setForm((f) => {
      const list = f[key];
      return {
        ...f,
        [key]: list.includes(value) ? list.filter((x) => x !== value) : [...list, value],
      };
    });
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    const dates = form.event_dates.map((d) => d.trim()).filter(Boolean);
    await api.postAuth("/api/admin/functions", {
      customer_name: form.customer_name,
      phone_number: form.phone_number,
      function_type: form.function_type,
      event_dates: dates.length ? dates : undefined,
      function_date: dates[0],
      total_amount: form.total_amount,
      advance_paid: form.advance_paid,
      customer_booking_notes: form.customer_booking_notes || undefined,
      services: form.services,
      complimentary: form.complimentary,
      album_status: form.album_status,
      video_status: form.video_status,
      overall_status: form.overall_status,
    });
    setShowForm(false);
    setForm(emptyForm());
    load();
  }

  async function handleDelete(f: StudioFunction) {
    if (!confirm(`Remove booking for ${f.customer_name}? This cannot be undone.`)) return;
    await api.deleteAuth(`/api/admin/functions/${f.id}`);
    load();
  }

  async function updateField(f: StudioFunction, overrides: Partial<StudioFunction>) {
    await api.putAuth(`/api/admin/functions/${f.id}`, functionPayload(f, overrides));
    load();
  }

  return (
    <div>
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="font-display text-3xl text-gold">Celebrations We&apos;re Honouring</h1>
          <p className="mt-2 text-sm text-text-muted">Track every booking, service, and delivery with care.</p>
        </div>
        <button type="button" onClick={() => setShowForm(!showForm)} className="btn-primary text-xs">
          {showForm ? "Close" : "Add Booking"}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleCreate} className="mt-6 glass space-y-6 rounded-sm p-6">
          <div className="grid gap-4 sm:grid-cols-2">
            <input
              placeholder="Beloved client's name"
              required
              value={form.customer_name}
              onChange={(e) => setForm({ ...form, customer_name: e.target.value })}
              className={inputClass}
            />
            <input
              placeholder="Phone number"
              required
              value={form.phone_number}
              onChange={(e) => setForm({ ...form, phone_number: e.target.value })}
              className={inputClass}
            />
            <input
              placeholder="Occasion (e.g. Wedding — 2 days)"
              required
              value={form.function_type}
              onChange={(e) => setForm({ ...form, function_type: e.target.value })}
              className={inputClass}
            />
            <div className="sm:col-span-2">
              <p className="mb-2 text-xs tracking-widest text-gold uppercase">Event dates (add each day)</p>
              {form.event_dates.map((d, i) => (
                <div key={i} className="mb-2 flex gap-2">
                  <input
                    type="date"
                    required={i === 0}
                    value={d}
                    onChange={(e) => {
                      const next = [...form.event_dates];
                      next[i] = e.target.value;
                      setForm({ ...form, event_dates: next });
                    }}
                    className={inputClass}
                  />
                  {form.event_dates.length > 1 && (
                    <button
                      type="button"
                      className="text-xs text-red-400"
                      onClick={() =>
                        setForm({
                          ...form,
                          event_dates: form.event_dates.filter((_, j) => j !== i),
                        })
                      }
                    >
                      Remove
                    </button>
                  )}
                </div>
              ))}
              <button
                type="button"
                className="text-xs text-gold hover:underline"
                onClick={() => setForm({ ...form, event_dates: [...form.event_dates, ""] })}
              >
                + Add another day
              </button>
            </div>
            <input
              type="number"
              placeholder="Total amount"
              value={form.total_amount || ""}
              onChange={(e) => setForm({ ...form, total_amount: Number(e.target.value) })}
              className={inputClass}
            />
            <input
              type="number"
              placeholder="Advance received"
              value={form.advance_paid || ""}
              onChange={(e) => setForm({ ...form, advance_paid: Number(e.target.value) })}
              className={inputClass}
            />
          </div>

          <div>
            <p className="mb-2 text-xs tracking-widest text-gold uppercase">What the family requested</p>
            <textarea
              placeholder="Notes from when booking was confirmed — traditions, wishes, special moments..."
              rows={3}
              value={form.customer_booking_notes}
              onChange={(e) => setForm({ ...form, customer_booking_notes: e.target.value })}
              className={inputClass}
            />
          </div>

          <div>
            <p className="mb-2 text-xs tracking-widest text-gold uppercase">Services booked</p>
            <div className="flex flex-wrap gap-2">
              {BOOKING_SERVICES.map((s) => (
                <label key={s} className="flex cursor-pointer items-center gap-2 rounded-sm border border-gold/20 px-3 py-2 text-xs">
                  <input
                    type="checkbox"
                    checked={form.services.includes(s)}
                    onChange={() => toggleList("services", s)}
                  />
                  {s}
                </label>
              ))}
            </div>
          </div>

          <div>
            <p className="mb-2 text-xs tracking-widest text-gold uppercase">Complimentary gifts</p>
            <div className="flex flex-wrap gap-2">
              {COMPLIMENTARY_ITEMS.map((s) => (
                <label key={s} className="flex cursor-pointer items-center gap-2 rounded-sm border border-gold/20 px-3 py-2 text-xs">
                  <input
                    type="checkbox"
                    checked={form.complimentary.includes(s)}
                    onChange={() => toggleList("complimentary", s)}
                  />
                  {s}
                </label>
              ))}
            </div>
          </div>

          <button type="submit" className="btn-primary w-full sm:w-auto">
            Save This Celebration
          </button>
        </form>
      )}

      <div className="mt-8 space-y-4">
        {items.map((f) => (
          <article key={f.id} className="rounded-sm border border-gold/10 bg-bg-card p-4">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <p className="font-medium text-text-main">{f.customer_name}</p>
                <p className="text-sm text-text-muted">
                  {f.function_type} · {formatEventDates(f)} · {f.phone_number}
                </p>
                <p className="mt-1 text-sm text-gold">Balance: ₹{f.balance_amount}</p>
              </div>
              <div className="flex flex-wrap gap-2">
                <button
                  type="button"
                  className="text-xs text-gold hover:underline"
                  onClick={() => setExpandedId(expandedId === f.id ? null : f.id)}
                >
                  {expandedId === f.id ? "Hide" : "Details"}
                </button>
                <button
                  type="button"
                  className="text-xs text-red-400 hover:underline"
                  onClick={() => handleDelete(f)}
                >
                  Delete
                </button>
              </div>
            </div>

            <div className="mt-4 grid gap-3 sm:grid-cols-3">
              <label className="text-xs text-text-muted">
                Overall
                <select
                  value={f.overall_status}
                  onChange={(e) => updateField(f, { overall_status: e.target.value })}
                  className="mt-1 block w-full bg-bg-muted text-sm"
                >
                  {FUNCTION_STATUSES.map((s) => (
                    <option key={s} value={s}>
                      {statusLabel(s)}
                    </option>
                  ))}
                </select>
              </label>
              <label className="text-xs text-text-muted">
                Album
                <select
                  value={f.album_status}
                  onChange={(e) => updateField(f, { album_status: e.target.value })}
                  className="mt-1 block w-full bg-bg-muted text-sm"
                >
                  {ALBUM_STATUSES.map((s) => (
                    <option key={s.value} value={s.value}>
                      {s.label}
                    </option>
                  ))}
                </select>
              </label>
              <label className="text-xs text-text-muted">
                Video
                <select
                  value={f.video_status}
                  onChange={(e) => updateField(f, { video_status: e.target.value })}
                  className="mt-1 block w-full bg-bg-muted text-sm"
                >
                  {VIDEO_STATUSES.map((s) => (
                    <option key={s.value} value={s.value}>
                      {s.label}
                    </option>
                  ))}
                </select>
              </label>
            </div>

            {expandedId === f.id && (
              <div className="mt-4 border-t border-gold/10 pt-4 text-sm text-text-muted">
                {f.customer_booking_notes && (
                  <p className="mb-2">
                    <span className="text-gold">Client wishes:</span> {f.customer_booking_notes}
                  </p>
                )}
                {f.services && f.services.length > 0 && (
                  <p className="mb-2">
                    <span className="text-gold">Services:</span> {f.services.join(", ")}
                  </p>
                )}
                {f.complimentary && f.complimentary.length > 0 && (
                  <p>
                    <span className="text-gold">Complimentary:</span> {f.complimentary.join(", ")}
                  </p>
                )}
              </div>
            )}
          </article>
        ))}
        {items.length === 0 && (
          <p className="py-12 text-center text-text-muted">No bookings yet — add your first celebration above.</p>
        )}
      </div>
    </div>
  );
}
