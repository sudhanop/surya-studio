import type { StudioFunction } from "./types";

function eventDateStrings(f: StudioFunction): string[] {
  if (f.event_dates?.length) {
    return f.event_dates.map((d) => d.event_date.slice(0, 10));
  }
  if (f.function_date) {
    return [f.function_date.slice(0, 10)];
  }
  return [];
}

/** Normalize function payload for API PUT/POST */
export function functionPayload(f: StudioFunction, overrides: Partial<StudioFunction> = {}) {
  const merged = { ...f, ...overrides };
  return {
    customer_name: merged.customer_name,
    phone_number: merged.phone_number,
    address: merged.address,
    function_type: merged.function_type,
    function_date: merged.function_date?.slice(0, 10),
    event_dates: eventDateStrings(merged),
    total_amount: merged.total_amount,
    advance_paid: merged.advance_paid,
    assigned_editor: merged.assigned_editor,
    assigned_date: merged.assigned_date?.slice(0, 10),
    album_status: merged.album_status,
    video_status: merged.video_status,
    delivery_status: merged.delivery_status,
    overall_status: merged.overall_status,
    customer_booking_notes: merged.customer_booking_notes,
    services: merged.services ?? [],
    complimentary: merged.complimentary ?? [],
    admin_notes: merged.admin_notes,
    drive_links: merged.drive_links,
    inquiry_id: merged.inquiry_id,
  };
}

export function formatEventDates(f: StudioFunction): string {
  const dates = eventDateStrings(f);
  if (!dates.length) return "—";
  return dates.join(", ");
}
