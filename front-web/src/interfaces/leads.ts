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
  email: string;
  phone: string;
 
}


export const mapLeadAPIToFront = (lead: LeadAPI): LeadFront => ({
  id: lead.id,
  businessName: lead.BusinessName || '',
  email: lead.Email || '',
  phone: lead.Phone || '',
});