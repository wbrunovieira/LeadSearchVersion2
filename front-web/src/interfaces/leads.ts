// interfaces/leads.ts
export interface LeadAPI {
  id: string;
  BusinessName: string;
  RegisteredName: string;
  FoundationDate?: { Time: string; Valid: boolean };
  Address: string;
  City: string;
  State: string;
  Country: string;
  ZIPCode: string;
  Owner: string;
  Source: string;
  Phone: string;
  Whatsapp: string;
  Website: string;
  Email: string;
  Instagram: string;
  Facebook: string;
  TikTok: string;
  CompanyRegistrationID: string;
  Categories: string;
  Rating: number;
  PriceLevel: number;
  UserRatingsTotal: number;
  Vicinity: string;
  PermanentlyClosed: boolean;
  CompanySize: string;
  Revenue: number;
  EmployeesCount: number;
  Description: string;
  PrimaryActivity: string;
  SecondaryActivities: string;
  Types: string;
  EquityCapital: number;
  BusinessStatus: string;
  Quality: string;
  SearchTerm: string;
  FieldsFilled: number;
  GoogleId: string;
  Category: string;
  Radius: number;
  CreatedAt: string;
  UpdatedAt: string;
}


export interface LeadFront {
  id: string;
  businessName: string;
  registeredName: string;
  email: string;
  phone: string;
  whatsapp: string;
  website: string;
  instagram: string;
  facebook: string;
  tiktok: string;
  cnpj: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  owner: string;
  category: string;
  rating: number;
  description: string;
  foundationDate: string;
}


export const mapLeadAPIToFront = (lead: LeadAPI): LeadFront => ({
  id: lead.id,
  businessName: lead.BusinessName || '',
  registeredName: lead.RegisteredName || '',
  email: lead.Email || '',
  phone: lead.Phone || '',
  whatsapp: lead.Whatsapp || '',
  website: lead.Website || '',
  instagram: lead.Instagram || '',
  facebook: lead.Facebook || '',
  tiktok: lead.TikTok || '',
  cnpj: lead.CompanyRegistrationID || '',
  address: lead.Address || '',
  city: lead.City || '',
  state: lead.State || '',
  zipCode: lead.ZIPCode || '',
  owner: lead.Owner || '',
  category: lead.Category || '',
  rating: lead.Rating || 0,
  description: lead.Description || '',
  foundationDate: lead.FoundationDate?.Time || ''
});