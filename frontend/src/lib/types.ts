export interface Category {
  id: number;
  name: string;
  slug: string;
  description?: string;
  cover_image?: string;
  display_order: number;
  photo_count?: number;
  video_count?: number;
}

export interface PortfolioMedia {
  id: number;
  category_id: number;
  title?: string;
  caption?: string;
  media_type: "photo" | "video";
  file_path: string;
  url?: string;
  thumbnail_url?: string;
  is_featured: boolean;
  display_order?: number;
  is_published?: boolean;
  category_slug?: string;
  category_name?: string;
}

export interface CategoryMediaResponse {
  category: Category;
  photos: PortfolioMedia[];
  videos: PortfolioMedia[];
  has_videos: boolean;
}

export interface Inquiry {
  id: number;
  customer_name: string;
  phone_number: string;
  occasion_type: string;
  wanted_date?: string;
  address?: string;
  message?: string;
  status: string;
  created_at: string;
}

export interface FunctionEventDate {
  id?: number;
  function_id?: number;
  event_date: string;
  day_label?: string;
  sort_order?: number;
}

export interface StudioFunction {
  id: number;
  inquiry_id?: number;
  customer_name: string;
  phone_number: string;
  address?: string;
  function_type: string;
  function_date: string;
  event_dates?: FunctionEventDate[];
  total_amount: number;
  advance_paid: number;
  balance_amount: number;
  assigned_editor?: string;
  assigned_date?: string;
  album_status: string;
  video_status: string;
  delivery_status: string;
  overall_status: string;
  customer_booking_notes?: string;
  services?: string[];
  complimentary?: string[];
  admin_notes?: string;
  drive_links?: string;
}

export interface DashboardStats {
  upcoming_functions: number;
  recent_inquiries: number;
  pending_albums: number;
  pending_video_edits: number;
  pending_deliveries: number;
  total_uploads: number;
  latest_portfolio_updates: number;
}

export interface ContactInfo {
  whatsapp: string;
  instagram_url: string;
  facebook_url: string;
  youtube_url: string;
  contact_email: string;
  phone_number: string;
  phone_secondary?: string;
  address?: string;
  pincode?: string;
  google_maps_embed: string;
}

export interface Testimonial {
  name: string;
  text: string;
}

export interface PublicSiteData {
  events_covered: string;
  years_of_craft: string;
  happy_families: string;
  owner_portrait_url: string;
  logo_url?: string;
  testimonials: Testimonial[];
  contact: ContactInfo;
}
