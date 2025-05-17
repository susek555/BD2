import { StringValidation, z } from 'zod';

export const SignupFormSchema = z
  .object({
    selector: z.enum(['P', 'C']),
    username: z
      .string()
      .min(2, { message: 'Username must be at least 2 characters long.' })
      .trim(),
    email: z.string().email({ message: 'Please enter a valid email.' }).trim(),
    password: z
      .string()
      .min(2, { message: 'Be at least 2 characters long' })
      // .regex(/[a-zA-Z]/, { message: 'Contain at least one letter.' })
      // .regex(/[0-9]/, { message: 'Contain at least one number.' })
      // .regex(/[^a-zA-Z0-9]/, {
      //   message: 'Contain at least one special character.',
      // })
      .trim(),
    confirm_password: z.string().trim(),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords don't match",
    path: ['confirm_password'],
  })
  .and(
    z.discriminatedUnion('selector', [
      z.object({
        selector: z.literal('P'),
        person_name: z.string().min(1, { message: 'Name is required' }),
        person_surname: z.string().min(1, { message: 'Surname is required' }),
      }),
      z.object({
        selector: z.literal('C'),
        company_name: z
          .string()
          .min(1, { message: 'Business name is required' }),
        company_nip: z
          .string()
          .min(10, { message: 'Be 10 characters long' })
          .max(10, { message: 'Be 10 characters long' })
          .regex(/^\d+$/, { message: 'Contain only digits' }),
      }),
    ]),
  )
  .transform((data) => {
    const { confirm_password: _, ...rest } = data;
    return rest;
  });

export type SignupForm = z.input<typeof SignupFormSchema>;

export type SignupFormState = {
  errors?: {
    username?: string[];
    email?: string[];
    password?: string[];
    confirm_password?: string[];
    selector?: string[];

    person_name?: string[];
    person_surname?: string[];

    company_name?: string[];
    company_nip?: string[];
  };
  values?: {
    selector?: string;
    username?: string;
    email?: string;

    person_name?: string;
    person_surname?: string;

    company_name?: string;
    company_nip?: string;
  };
};

export type LoginError = {
  credentials?: string[];
  server?: string[];
};

// Filters

export type FilterFieldData = {
  fieldName: string;
  options: string[];
  selected?: string[];
};

// Ranges

export type RangeFieldData = {
  fieldName: string;
  range: {
    min: number | null;
    max: number | null;
  };
};

// SearchParams

export type SearchParams = {
  query?: string;
  page?: number;
  orderKey?: string;
  isOrderDesc?: boolean;
  offerType?: string;
  producers?: string[];
  gearboxes?: string[];
  fuelTypes?: string[];
  price?: {
    min?: number;
    max?: number;
  };
  mileage?: {
    min?: number;
    max?: number;
  };
  year?: {
    min?: number;
    max?: number;
  };
};

export const parseIntOrUndefined = (
  value: string | undefined,
): number | undefined => {
  if (!value) return undefined;
  const parsed = parseInt(value, 10);
  return isNaN(parsed) ? undefined : parsed;
};

export const parseArrayOrUndefined = (
  value: string | string[] | undefined,
): string[] | undefined => {
  if (!value) return undefined;
  return Array.isArray(value) ? value : [value];
};

// SaleOffer

export interface BaseOffer {
  id: string;
  // image: type to  be determined;
  name: string; // producer and model
  productionYear: number;
  mileage: number;
  color: string;
  price: number;
  isAuction: boolean;
}

export interface SaleOffer extends BaseOffer {
  isFavorite: boolean;
}

export interface HistoryOffer extends BaseOffer {
  dateEnd: string;
}

// Account page definitions

export interface Tab {
  id: string;
  label: string;
  href: string;
}

export const profileTabs: Tab[] = [
  { id: 'activity', label: 'Activity', href: '/account/activity' },
  { id: 'listings', label: 'Listings', href: '/account/listings' },
  { id: 'favorites', label: 'Favorites', href: '/account/favorites' },
  { id: 'reviews', label: 'Reviews', href: '/account/reviews' },
  { id: 'settings', label: 'Settings', href: '/account/settings' },
];

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  personName?: string;
  personSurname?: string;
  companyName?: string;
  companyNip?: string;
}

export type SaleOfferDetails = {
  name: string; // producer and model
  price?: number;
  isAuction: boolean;
  auctionData?: AuctionData;
  isActive: boolean;
  imagesURLs: string[]; // URLs to images
  details: OfferDetails[];
  description: string;
  // location: string;
  sellerName: string;
}

export type AuctionData = {
  currentBid: number;
  endDate: Date;
}

export type OfferDetails = {
  name: string;
  value: string;
}




// Reviews

interface ReviewUser {
  id: number;
  username: string;
}

export interface Review {
  id: number;
  description: string;
  rating: number;
  date: string;
  reviewer: ReviewUser;
  reviewee: ReviewUser;
}
