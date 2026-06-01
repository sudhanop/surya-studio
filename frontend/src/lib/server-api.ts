/** Server-side fetch for RSC — uses revalidation for public pages */
const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function serverFetch<T>(
  path: string,
  revalidateSeconds = 60
): Promise<T | null> {
  try {
    const res = await fetch(`${API_URL}${path}`, {
      next: { revalidate: revalidateSeconds },
    });
    if (!res.ok) return null;
    return res.json() as Promise<T>;
  } catch {
    return null;
  }
}
