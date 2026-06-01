"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api, asArray } from "@/lib/api";
import type { DashboardStats, Inquiry, StudioFunction } from "@/lib/types";

export default function AdminDashboardPage() {
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [upcoming, setUpcoming] = useState<StudioFunction[]>([]);
  const [inquiries, setInquiries] = useState<Inquiry[]>([]);

  useEffect(() => {
    Promise.all([
      api.getAuth<DashboardStats>("/api/admin/dashboard"),
      api.getAuth<StudioFunction[]>("/api/admin/functions/upcoming"),
      api.getAuth<Inquiry[]>("/api/admin/inquiries?status=new"),
    ]).then(([s, u, i]) => {
      setStats(s);
      setUpcoming(asArray(u));
      setInquiries(asArray(i).slice(0, 5));
    }).catch(console.error);
  }, []);

  const cards = stats
    ? [
        { label: "Upcoming Functions", value: stats.upcoming_functions },
        { label: "New Inquiries", value: stats.recent_inquiries },
        { label: "Pending Albums", value: stats.pending_albums },
        { label: "Pending Video Edits", value: stats.pending_video_edits },
        { label: "Pending Deliveries", value: stats.pending_deliveries },
        { label: "Total Uploads", value: stats.total_uploads },
      ]
    : [];

  return (
    <div>
      <h1 className="font-display text-3xl text-gold">Dashboard</h1>
      <p className="mt-2 text-text-muted">Studio overview</p>

      <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {cards.map((c) => (
          <div key={c.label} className="glass rounded-sm p-6">
            <p className="text-3xl font-display text-gold">{c.value}</p>
            <p className="mt-1 text-xs tracking-widest text-text-muted uppercase">{c.label}</p>
          </div>
        ))}
      </div>

      <div className="mt-12 grid gap-8 lg:grid-cols-2">
        <section>
          <div className="mb-4 flex items-center justify-between">
            <h2 className="font-display text-xl">Upcoming Functions</h2>
            <Link href="/admin/functions" className="text-xs text-gold">View all</Link>
          </div>
          <div className="space-y-3">
            {upcoming.length === 0 && <p className="text-sm text-text-muted">No upcoming functions</p>}
            {upcoming.map((f) => (
              <div key={f.id} className="rounded-sm border border-gold/10 bg-bg-card p-4">
                <p className="font-medium">{f.customer_name}</p>
                <p className="text-sm text-text-muted">{f.function_type} · {f.function_date?.slice(0, 10)}</p>
              </div>
            ))}
          </div>
        </section>

        <section>
          <div className="mb-4 flex items-center justify-between">
            <h2 className="font-display text-xl">Recent Inquiries</h2>
            <Link href="/admin/inquiries" className="text-xs text-gold">View all</Link>
          </div>
          <div className="space-y-3">
            {inquiries.length === 0 && <p className="text-sm text-text-muted">No new inquiries</p>}
            {inquiries.map((i) => (
              <div key={i.id} className="rounded-sm border border-gold/10 bg-bg-card p-4">
                <p className="font-medium">{i.customer_name}</p>
                <p className="text-sm text-text-muted">{i.occasion_type} · {i.phone_number}</p>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
