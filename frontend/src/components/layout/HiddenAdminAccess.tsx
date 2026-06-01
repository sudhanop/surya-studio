"use client";

import { useRouter } from "next/navigation";
import { useRef } from "react";

/** Five quick clicks on the studio mark opens admin login (not linked in navigation). */
export function HiddenAdminAccess({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const clicks = useRef(0);
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  function handleClick() {
    clicks.current += 1;
    if (timer.current) clearTimeout(timer.current);
    timer.current = setTimeout(() => {
      clicks.current = 0;
    }, 2000);
    if (clicks.current >= 5) {
      clicks.current = 0;
      router.push("/admin/login");
    }
  }

  return (
    <span
      role="presentation"
      onClick={handleClick}
      onKeyDown={(e) => e.key === "Enter" && handleClick()}
      className="cursor-default select-none"
    >
      {children}
    </span>
  );
}
