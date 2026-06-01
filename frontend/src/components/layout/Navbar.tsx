"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Logo } from "@/components/ui/Logo";

const links = [
  { href: "/", label: "Home" },
  { href: "/portfolio", label: "Portfolio" },
  { href: "/about", label: "About" },
  { href: "/contact", label: "Contact" },
];

export function Navbar() {
  const pathname = usePathname();
  const [scrolled, setScrolled] = useState(false);
  const [open, setOpen] = useState(false);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 40);
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  return (
    <header
      className={`fixed inset-x-0 top-0 z-50 transition-all duration-500 ${
        scrolled ? "glass py-3 shadow-lg" : "bg-transparent py-5"
      }`}
    >
      <motion.div layout className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <motion.div layout className="flex items-center justify-between">
          <Logo />
          <nav className="hidden items-center gap-8 md:flex">
            {links.map((l) => (
              <Link
                key={l.href}
                href={l.href}
                className={`text-xs tracking-[0.2em] uppercase transition ${
                  pathname === l.href ? "text-gold" : "text-text-muted hover:text-gold"
                }`}
              >
                {l.label}
              </Link>
            ))}
          </nav>
          <button
            type="button"
            className="text-gold md:hidden"
            onClick={() => setOpen(!open)}
            aria-label="Menu"
          >
            <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              {open ? (
                <path strokeLinecap="square" strokeWidth={1.5} d="M6 6l12 12M6 18L18 6" />
              ) : (
                <path strokeLinecap="square" strokeWidth={1.5} d="M4 8h16M4 16h16" />
              )}
            </svg>
          </button>
        </motion.div>
        <AnimatePresence>
          {open && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: "auto" }}
              exit={{ opacity: 0, height: 0 }}
              className="mt-3 overflow-hidden border-t border-gold/10 md:hidden"
            >
              <nav className="flex flex-col gap-4 py-6">
                {links.map((l) => (
                  <Link
                    key={l.href}
                    href={l.href}
                    onClick={() => setOpen(false)}
                    className="text-sm tracking-widest uppercase text-text-muted hover:text-gold"
                  >
                    {l.label}
                  </Link>
                ))}
              </nav>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </header>
  );
}
