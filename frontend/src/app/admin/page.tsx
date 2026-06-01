import { redirect } from "next/navigation";

/** /admin is not advertised — use hidden footer access to reach /admin/login */
export default function AdminIndexPage() {
  redirect("/");
}
