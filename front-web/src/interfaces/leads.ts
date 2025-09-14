// interfaces/leadas.ts
export interface LeadAPI {
  id: string;
  businessName: string;
  registeredName: string;
  foundationDate?: string; 
  address: string;
  city: string;
  state: string;
  country: string;
  zipCode: string;
  owner: string;
  source: string;
  phone: string;
  whatsapp: string;
  website: string;
  email: string;
  instagram: string;
  facebook: string;
  tikTok: string;
  companyRegistrationID: string;
  categories: string;
  rating: number;
  priceLevel: number;
  userRatingsTotal: number;
  vicinity: string;
  permanentlyClosed: boolean;
  companySize: string;
  revenue: number;
  employeesCount: number;
  description: string;
  primaryActivity: string;
  secondaryActivities: string;
  types: string;
  equityCapital: number;
  businessStatus: string;
  quality: string;
  searchTerm: string;
  fieldsFilled: number;
  googleId: string;
  category: string;
  radius: number;
  createdAt: string;
  updatedAt: string;
}


export interface LeadFront {
  id: string;
  businessName: string;
  email: string;
  phone: string;
 
}


export const mapLeadAPIToFront = (lead: LeadAPI): LeadFront => ({
  id: lead.id,
  businessName: lead.businessName,
  email: lead.email,
  phone: lead.phone,
});