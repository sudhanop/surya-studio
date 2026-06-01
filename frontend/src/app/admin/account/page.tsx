"use client";

import { useState } from "react";
import { api, setToken } from "@/lib/api";

const inputClass =
  "w-full border border-gold/25 bg-bg-muted px-4 py-3 text-sm text-text-main outline-none focus:border-gold/60";

export default function AdminAccountPage() {
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState("");
  const [ok, setOk] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNew, setShowNew] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError("");
    setOk("");

    if (newPassword && newPassword !== confirmPassword) {
      setError("New password and confirmation do not match");
      return;
    }

    const fd = new FormData(e.currentTarget);
    const currentPassword = String(fd.get("current_password") || "");
    const newUsername = String(fd.get("new_username") || "");

    if (!currentPassword) {
      setError("Current password is required");
      return;
    }
    if (!newUsername && !newPassword) {
      setError("Enter a new username and/or a new password");
      return;
    }

    setSaving(true);
    try {
      const res = await api.putAuth<{ token: string }>("/api/admin/account", {
        current_password: currentPassword,
        new_username: newUsername || "",
        new_password: newPassword || "",
      });
      if (res.token) setToken(res.token);
      e.currentTarget.reset();
      setNewPassword("");
      setConfirmPassword("");
      setOk("Updated successfully");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Update failed");
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="max-w-xl">
      <h1 className="font-display text-3xl text-gold">Admin Account</h1>
      <p className="mt-2 text-sm text-text-muted">Change your admin username and password.</p>

      <form onSubmit={handleSubmit} className="mt-8 space-y-6 glass rounded-sm p-6">
        <label className="block text-xs tracking-widest text-gold uppercase">
          Current password
          <div className="relative mt-2">
            <input name="current_password" type={showCurrent ? "text" : "password"} required className={`${inputClass} pr-12`} />
            <button
              type="button"
              onClick={() => setShowCurrent((v) => !v)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gold/80 hover:text-gold"
              aria-label={showCurrent ? "Hide password" : "Show password"}
            >
              {showCurrent ? (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M3 3l18 18" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
                  <path
                    d="M10.58 10.58A2 2 0 0012 16a2 2 0 001.42-3.42"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M9.9 5.12A10.94 10.94 0 0112 5c6 0 10 7 10 7a18.25 18.25 0 01-3.34 4.62"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M6.23 6.23A18.25 18.25 0 002 12s4 7 10 7c.72 0 1.41-.1 2.07-.27"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              ) : (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path
                    d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7-10-7-10-7Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M12 15a3 3 0 100-6 3 3 0 000 6Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              )}
            </button>
          </div>
        </label>

        <label className="block text-xs tracking-widest text-gold uppercase">
          New username (email)
          <input name="new_username" type="email" placeholder="new-admin@example.com" className={`${inputClass} mt-2`} />
        </label>

        <label className="block text-xs tracking-widest text-gold uppercase">
          New password
          <div className="relative mt-2">
            <input
              type={showNew ? "text" : "password"}
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              className={`${inputClass} pr-12`}
            />
            <button
              type="button"
              onClick={() => setShowNew((v) => !v)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gold/80 hover:text-gold"
              aria-label={showNew ? "Hide password" : "Show password"}
            >
              {showNew ? (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M3 3l18 18" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
                  <path
                    d="M10.58 10.58A2 2 0 0012 16a2 2 0 001.42-3.42"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M9.9 5.12A10.94 10.94 0 0112 5c6 0 10 7 10 7a18.25 18.25 0 01-3.34 4.62"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M6.23 6.23A18.25 18.25 0 002 12s4 7 10 7c.72 0 1.41-.1 2.07-.27"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              ) : (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path
                    d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7-10-7-10-7Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M12 15a3 3 0 100-6 3 3 0 000 6Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              )}
            </button>
          </div>
        </label>

        <label className="block text-xs tracking-widest text-gold uppercase">
          Confirm new password
          <div className="relative mt-2">
            <input
              type={showConfirm ? "text" : "password"}
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className={`${inputClass} pr-12`}
            />
            <button
              type="button"
              onClick={() => setShowConfirm((v) => !v)}
              className="absolute right-3 top-1/2 -translate-y-1/2 text-gold/80 hover:text-gold"
              aria-label={showConfirm ? "Hide password" : "Show password"}
            >
              {showConfirm ? (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M3 3l18 18" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" />
                  <path
                    d="M10.58 10.58A2 2 0 0012 16a2 2 0 001.42-3.42"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M9.9 5.12A10.94 10.94 0 0112 5c6 0 10 7 10 7a18.25 18.25 0 01-3.34 4.62"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M6.23 6.23A18.25 18.25 0 002 12s4 7 10 7c.72 0 1.41-.1 2.07-.27"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              ) : (
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path
                    d="M2 12s4-7 10-7 10 7 10 7-4 7-10 7-10-7-10-7Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                  <path
                    d="M12 15a3 3 0 100-6 3 3 0 000 6Z"
                    stroke="currentColor"
                    strokeWidth="1.8"
                    strokeLinecap="round"
                    strokeLinejoin="round"
                  />
                </svg>
              )}
            </button>
          </div>
        </label>

        {error && <p className="text-sm text-red-400">{error}</p>}
        {ok && <p className="text-sm text-green-300">{ok}</p>}

        <button type="submit" disabled={saving} className="btn-primary">
          {saving ? "Saving..." : "Save changes"}
        </button>
      </form>
    </div>
  );
}
