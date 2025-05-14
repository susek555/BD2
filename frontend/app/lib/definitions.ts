import { number, StringValidation, z } from 'zod';

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

export type ModelFieldData = {
  producers: FilterFieldData;
  models: string[][];
}

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

// Add Offer

export const AddOfferFormSchema = z.object({
  producer: z.string().min(1, { message: 'Producer is required' }),
  model: z.string().min(1, { message: 'Model is required' }),
  color: z.string().min(1, { message: 'Color is required' }),
  fuelType: z.string().min(1, { message: 'Fuel type is required' }),
  gearbox: z.string().min(1, { message: 'Gearbox is required' }),
  driveType: z.string().min(1, { message: 'Drive type is required' }),
  country: z.string().min(1, { message: 'Country is required' }),
  productionYear: z
    .number()
    .min(1900, { message: 'Year must be greater than 1900' })
    .max(new Date().getFullYear(), { message: 'Year must be less than or equal to the current year' }),
  mileage: z
    .number()
    .min(0, { message: 'Mileage must be greater than or equal to 0' })
    .max(1_000_000, { message: 'Mileage must be less than or equal to 1,000,000' }),
  numberOfOwners: z
    .number()
    .min(0, { message: 'Number of owners must be greater than or equal to 0' })
    .max(100, { message: 'Number of owners must be less than or equal to 100' }),
  numberOfSeats: z
    .number()
    .min(1, { message: 'Number of seats must be greater than or equal to 1' })
    .max(100, { message: 'Number of seats must be less than or equal to 100' }),
  power: z
    .number()
    .min(0, { message: 'Power must be greater than or equal to 0' })
    .max(2_000, { message: 'Power must be less than or equal to 1,000' }),
  dateOfFirstRegistration: z
    .date()
    .refine((date) => date.getFullYear() >= 1900, { message: 'Date must be greater than 1900' })
    .refine((date) => date <= new Date(), { message: 'Date must be less than or equal to the current date' }),
  plateNumber: z.string().min(1, { message: 'Plate number is required' }),
  engineDisplacement: z
    .number()
    .min(0, { message: 'Engine displacement must be greater than or equal to 0' })
    .max(10_000, { message: 'Engine displacement must be less than or equal to 10,000' }),
  location: z.string().min(1, { message: 'Location is required' }),
  price: z
    .number()
    .min(0, { message: 'Price must be greater than or equal to 0' })
    .max(10_000_000, { message: 'Price must be less than or equal to 10,000,000' }),
  isAuction: z.boolean(),
  auctionEndDate: z
    .date()
    .optional()
    .refine((date) => date && date > new Date(), { message: 'Date must be greater than the current date' }),
  buyNowAuctionPrice: z
    .number()
    .optional()
    .refine((price) => price && price > 0, { message: 'Price must be greater than 0' })
    .refine((price) => price && price< 10_000_000, { message: 'Price must be less than 10,000,000' }),
  description: z.string().min(1, { message: 'Description is required' }),
  images: z
    .array(z.instanceof(File))
    .min(1, { message: 'At least one image is required' })
    .max(10, { message: 'A maximum of 10 images is allowed' })
    .refine((files) => files.every(file => file.size <= 5 * 1024 * 1024), {
      message: 'Each image must be less than 5MB',
    }),
}).refine((data) => {
  if (data.isAuction && !data.auctionEndDate) {
    return false;
  }
  if (data.isAuction && !data.buyNowAuctionPrice) {
    return false;
  }
  if (data.isAuction && data.buyNowAuctionPrice && data.buyNowAuctionPrice <= data.price) {
    return false;
  }
  return true;
}, {
  message: 'Buy now auction price must be greater than the regular price',
  path: ['buyNowAuctionPrice'],
});

export type AddOfferFormState = {
  errors?: {
    producer?: string[];
    model?: string[];
    color?: string[];
    fuelType?: string[];
    gearbox?: string[];
    driveType?: string[];
    country?: string[];
    productionYear?: number[];
    mileage?: number[];
    numberOfOwners?: number[];
    numberOfSeats?: number[];
    power?: number[];
    dateOfFirstRegistration?: Date[];
    plateNumber?: string[];
    engineDisplacement?: number[];
    location?: string[];
    price?: number[];
    isAuction?: boolean[];
    auctionEndDate?: Date[];
    auctionStartPrice?: number[];
    description?: string[];
    images?: string[];
  }
  values?: {
    producer?: string;
    model?: string;
    color?: string;
    fuelType?: string;
    gearbox?: string;
    driveType?: string;
    country?: string;
    productionYear?: number;
    mileage?: number;
    numberOfOwners?: number;
    numberOfSeats?: number;
    power?: number;
    dateOfFirstRegistration?: Date;
    plateNumber?: string;
    engineDisplacement?: number;
    location?: string;
    price?: number;
    isAuction?: boolean;
    auctionEndDate?: Date;
    buyNowAuctionPrice?: number;
    description?: string;
    images?: File[]; // Array of image files
  }
}

export type AddOfferFormData = {
  producers: string[];
  models: string[][];
  colors: string[];
  fuelTypes: string[];
  gearboxes: string[];
  driveTypes: string[];
  countries: string[];
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

export interface ReviewSearchParams {
  is_order_desc: boolean;
  order_key: 'rating' | 'date';
  pagination: {
    page: number;
    page_size: number;
  };
  ratings?: number[];
}

export interface RatingPercentages {
  '1': number;
  '2': number;
  '3': number;
  '4': number;
  '5': number;
}
