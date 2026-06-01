import type { ContactInfo } from "./types";

export const SITE = {
  name: "Surya Photography",
  tagline: "Where every heartbeat becomes a frame you'll treasure forever.",
  description:
    "A Rasipuram studio devoted to wedding stories, family milestones, and cinematic memories—crafted with warmth, artistry, and soul.",
};

export const OCCASIONS = [
  "Wedding",
  "Baby Shower",
  "Puberty",
  "Birthday",
  "Ear Piercing",
  "Couple Shoot",
  "Maternity",
  "Outdoor",
  "Temple Function",
  "Other",
];

export const BOOKING_SERVICES = [
  "Traditional Photo",
  "Traditional Video",
  "Candid Photo",
  "Candid Video",
  "Outdoor Shoot",
  "Pre Wedding",
  "Post Wedding Shoot",
  "Drone",
  "Selfie Booth",
  "LED Wall",
] as const;

export const COMPLIMENTARY_ITEMS = [
  "Frame",
  "Cup",
  "Calendar",
] as const;

export const FUNCTION_STATUSES = [
  "upcoming",
  "shoot_completed",
  "photos_selected",
  "editing",
  "completed",
  "album_ready",
  "delivered",
] as const;

export const ALBUM_STATUSES = [
  { value: "not_started", label: "Not started" },
  { value: "in_progress", label: "Work started" },
  { value: "design_completed", label: "Design completed" },
  { value: "printed", label: "Printed" },
  { value: "delivered", label: "Delivered" },
] as const;

export const VIDEO_STATUSES = [
  { value: "not_started", label: "Not started" },
  { value: "in_progress", label: "Work started" },
  { value: "editing", label: "Editing" },
  { value: "completed", label: "Completed" },
  { value: "delivered", label: "Delivered" },
] as const;

/** Official studio contact — used when API is unavailable or DB is empty */
export const DEFAULT_CONTACT = {
  whatsapp: "919715241568",
  instagram_url: "https://www.instagram.com/surya_photography_nkl",
  instagram_handle: "surya_photography_nkl",
  facebook_url: "",
  youtube_url: "https://www.youtube.com/@suryaphotography4303",
  youtube_handle: "suryaphotography4303",
  contact_email: "suryaphotographyrsp@gmail.com",
  phone_number: "9715241568",
  phone_secondary: "8884897499",
  address:
    "Surya Photography, near DNC (Chamundi) Theater, opposite to Adhisindha Thirumana Mandabam, Pattanam road, Rasipuram",
  pincode: "637408",
  google_maps_embed: "",
};

export const DEFAULT_TESTIMONIALS = [
  {
    name: "Priya & Arun",
    text: "Every frame felt like a movie. Surya Photography captured our wedding with pure magic.",
  },
  {
    name: "Lakshmi Family",
    text: "Professional, warm, and incredibly talented. Our puberty function album is breathtaking.",
  },
  {
    name: "Divya",
    text: "The maternity shoot was dreamy. Highly recommend for anyone who wants cinematic quality.",
  },
] as const;

export const DEFAULT_SITE_DATA = {
  events_covered: "500+",
  years_of_craft: "10+",
  happy_families: "1000+",
  owner_portrait_url: "",
  logo_url: "",
  testimonials: [...DEFAULT_TESTIMONIALS],
  contact: DEFAULT_CONTACT,
};

/** Always show official contact even if API returns partial/old data */
export function mergeContact(info: Partial<ContactInfo> | null | undefined): ContactInfo {
  if (!info) return DEFAULT_CONTACT;
  return {
    ...DEFAULT_CONTACT,
    ...info,
    contact_email: info.contact_email || DEFAULT_CONTACT.contact_email,
    phone_number: info.phone_number || DEFAULT_CONTACT.phone_number,
    phone_secondary: info.phone_secondary || DEFAULT_CONTACT.phone_secondary,
    address: info.address || DEFAULT_CONTACT.address,
    pincode: info.pincode || DEFAULT_CONTACT.pincode,
    instagram_url: info.instagram_url || DEFAULT_CONTACT.instagram_url,
    youtube_url: info.youtube_url || DEFAULT_CONTACT.youtube_url,
    whatsapp: info.whatsapp || DEFAULT_CONTACT.whatsapp,
  };
}
