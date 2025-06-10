import { z } from "zod";
import { SaleOfferDetails } from "./sale-offer-details";

// Add / Edit Offer

export const offerActionEnum = {
    ADD_OFFER: true,
    EDIT_OFFER: false
}

export const OfferDetailsFormSchema = z.object({
  manufacturer: z.string().min(1, { message: 'Producer is required' }),
  model: z.string().min(1, { message: 'Model is required' }),
  color: z.string().min(1, { message: 'Color is required' }),
  fuel_type: z.string().min(1, { message: 'Fuel type is required' }),
  transmission: z.string().min(1, { message: 'Gearbox is required' }),
  drive: z.string().min(1, { message: 'Drive type is required' }),
  production_year: z
    .number()
    .min(1900, { message: 'Year must be greater than 1900' })
    .max(new Date().getFullYear(), { message: 'Year must be less than or equal to the current year' }),
  mileage: z
    .number()
    .min(0, { message: 'Mileage must be greater than or equal to 0' })
    .max(1_000_000, { message: 'Mileage must be less than or equal to 1,000,000' }),
  number_of_doors: z
    .number()
    .min(1, { message: 'Number of doors must be greater than or equal to 0' })
    .max(6, { message: 'Number of doors must be less than or equal to 6' }),
  number_of_seats: z
    .number()
    .min(2, { message: 'Number of seats must be greater than or equal to 1' })
    .max(100, { message: 'Number of seats must be less than or equal to 100' }),
  number_of_gears: z
    .number()
    .min(1, { message: 'Number of gears must be greater than or equal to 0' })
    .max(10, { message: 'Number of gears must be less than or equal to 10' }),
  engine_power: z
    .number()
    .min(0, { message: 'Power must be greater than or equal to 0' })
    .max(9_999, { message: 'Power must be less than or equal to 9,999' }),
  registration_date: z
    .string()
    .refine((date) => {
      const parsedDate = new Date(date);
      return parsedDate.getFullYear() >= 1900;
    }, { message: 'Date must be greater than 1900' })
    .refine((date) => {
      const parsedDate = new Date(date);
      return parsedDate <= new Date();
    }, { message: 'Date must be less than or equal to the current date' }),
  registration_number: z.string().min(1, { message: 'Plate number is required' }),
  engine_capacity: z
    .number()
    .min(0, { message: 'Engine displacement must be greater than or equal to 0' })
    .max(9_000, { message: 'Engine displacement must be less than or equal to 9,000' }),
  vin: z
    .string()
    .min(17, { message: 'VIN must be 17 characters long' })
    .max(17, { message: 'VIN must be 17 characters long' })
    .regex(/^[A-HJ-NPR-Z0-9]+$/, { message: 'VIN must contain only valid characters' }),
  description: z.string().min(1, { message: 'Description is required' }),
}).superRefine((data, ctx) => {
  if (data.registration_date && data.production_year) {
    const registrationDate = new Date(data.registration_date);
    if (registrationDate.getFullYear() < data.production_year) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Registration date cannot be earlier than production year',
        path: ['registration_date'],
      });
    }
  }
})

export const OfferPricingFormSchema = z.object({
  price: z
    .number()
    .min(0, { message: 'Price must be greater than or equal to 0' })
    .max(10_000_000, { message: 'Price must be less than or equal to 10,000,000' }),
  margin: z
    .number()
    .min(3, { message: 'Margin must be greater than or equal to 3' })
    .max(10,{ message: 'Margin must be less than or equal to 10' }),
  is_auction: z.boolean(),
  date_end: z
    .string()
    .optional()
    .refine((date) => {
      if (!date) return true;

      const regex = /^([01]\d|2[0-3]):([0-5]\d) (\d{4})-(\d{2})-(\d{2})$/;
      if (!regex.test(date)) return false;

      const [time, dateStr] = date.split(' ');
      const [hours, minutes] = time.split(':').map(Number);
      const [year, month, day] = dateStr.split('-').map(Number);

      const parsedDate = new Date(year, month - 1, day, hours, minutes);
      const now = new Date();

      // Check if the date is in the future
      if (parsedDate <= now) return false;

      // Check if the date is within a week from now
      const oneWeekFromNow = new Date();
      oneWeekFromNow.setDate(oneWeekFromNow.getDate() + 7);

      return parsedDate <= oneWeekFromNow;
    }, { message: 'Date must be in the future but not more than a week from now' }),
  buy_now_price: z
    .number()
    .nullish()
    .transform(val => val === 0 || val === null || val === undefined ? undefined : val)
    .refine(
      (price) => price === undefined || price > 0,
      { message: 'Price must be greater than 0' }
    )
    .refine(
      (price) => price === undefined || price < 10_000_000,
      { message: 'Price must be less than 10,000,000' }
    ),
}).superRefine((data, ctx) => {
  if (data.is_auction) {
    if (!data.date_end) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Auction end date is required when auction is enabled',
        path: ['date_end'],
      });
    }
  }
  if (
    data.buy_now_price !== undefined &&
    data.buy_now_price <= data.price
  ) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: 'Buy now auction price must be greater than the regular price',
      path: ['buy_now_price'],
    });
  }
});

export const OfferImagesFormSchema = (minImages: number = 3) => z.object({
  images: z
    .array(z.instanceof(File))
    // .min(minImages, { message: `At least ${minImages} image${minImages === 1 ? '' : 's'} required` })
    .max(10, { message: 'A maximum of 10 images is allowed' })
    .refine((files) => files.every(file => file.size <= 5 * 1024 * 1024), {
      message: 'Each image must be less than 5MB',
    })
    .refine((files) => {
      // Filter out ghost files with empty names
      const validFiles = files.filter(file => file.name.trim() !== "");
      return minImages === 0 || validFiles.length >= minImages;
    }, {
      message: `At least ${minImages} image${minImages === 1 ? '' : 's'} required`,
    }),
});

export const CombinedOfferFormSchema = z
  .object({
    ...OfferDetailsFormSchema._def.schema.shape,
    ...OfferPricingFormSchema._def.schema.shape,
  })
  .superRefine((data, ctx) => {
    if (data.registration_date && data.production_year) {
      const registrationDate = new Date(data.registration_date);
      if (registrationDate.getFullYear() < data.production_year) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: 'Registration date cannot be earlier than production year',
          path: ['registration_date'],
        });
      }
    }

    if (data.is_auction) {
      if (!data.date_end) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: 'Auction end date is required when auction is enabled',
          path: ['date_end'],
        });
      }
    }

    if (
      data.buy_now_price !== undefined &&
      data.buy_now_price <= data.price
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Buy now auction price must be greater than the regular price',
        path: ['buy_now_price'],
      });
    }
  })
  .transform((data) => {
    if (data.buy_now_price === undefined) {
      const { buy_now_price: buy_now_price, ...rest } = data;
      return rest;
    }
    return data;
  });

export const CombinedImagesOfferFormSchema = (minImages: number = 3) => z
  .object({
    ...OfferDetailsFormSchema._def.schema.shape,
    ...OfferPricingFormSchema._def.schema.shape,
    ...OfferImagesFormSchema(minImages).shape,
  })
  .superRefine((data, ctx) => {
    if (data.registration_date && data.production_year) {
      const registrationDate = new Date(data.registration_date);
      if (registrationDate.getFullYear() < data.production_year) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: 'Registration date cannot be earlier than production year',
          path: ['registration_date'],
        });
      }
    }

    if (data.is_auction) {
      if (!data.date_end) {
        ctx.addIssue({
          code: z.ZodIssueCode.custom,
          message: 'Auction end date is required when auction is enabled',
          path: ['date_end'],
        });
      }
    }

    if (
      data.buy_now_price !== undefined &&
      data.buy_now_price <= data.price
    ) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Buy now auction price must be greater than the regular price',
        path: ['buy_now_price'],
      });
    }
  })
  .transform((data) => {
    if (data.buy_now_price === undefined) {
      const { buy_now_price: buy_now_price, ...rest } = data;
      return rest;
    }
    return data;
  });


export function editFormWrapper(offer: SaleOfferDetails) : Partial<OfferFormState['values']> {
  return {
    manufacturer: offer.manufacturer,
    model: offer.name.replace(offer.manufacturer + ' ', ''),
    color: offer.details.find(detail => detail.name === 'Color')?.value || '',
    fuel_type: offer.details.find(detail => detail.name === 'Fuel Type')?.value || '',
    transmission: offer.details.find(detail => detail.name === 'Transmission')?.value || '',
    drive: offer.details.find(detail => detail.name === 'Drive')?.value || '',
    production_year: parseInt(offer.details.find(detail => detail.name === 'Production Year')?.value || '0'),
    mileage: parseInt(offer.details.find(detail => detail.name === 'Mileage')?.value.replace(' km', '') || '0'),
    number_of_doors: parseInt(offer.details.find(detail => detail.name === 'Number of doors')?.value || '0'),
    number_of_seats: parseInt(offer.details.find(detail => detail.name === 'Number of seats')?.value || '0'),
    number_of_gears: parseInt(offer.details.find(detail => detail.name === 'Number of gears')?.value || '0'),
    engine_power: parseInt(offer.details.find(detail => detail.name === 'Engine Power')?.value.replace(' HP', '') || '0'),
    engine_capacity: parseInt(offer.details.find(detail => detail.name === 'Engine Capacity')?.value.replace(' L', '') || '0'),
    registration_date: offer.details.find(detail => detail.name === 'Date of first registration')?.value || '',
    registration_number: offer.details.find(detail => detail.name === 'Registration number')?.value || '',
    vin: offer.details.find(detail => detail.name === 'VIN')?.value || '',
    description: offer.description || '',
    price: offer.isAuction ? offer.auctionData?.currentBid || 0 : offer.price || 0,
    is_auction: offer.isAuction || false,
    margin: offer.margin || 0,
    date_end: offer.isAuction ?
      offer.auctionData?.endDate.toLocaleTimeString('en-GB', {hour: '2-digit', minute:'2-digit', hour12: false}) + ' ' +
      offer.auctionData?.endDate.toISOString().split('T')[0] :
      undefined,
    buy_now_price: offer.isAuction ? offer.price : undefined,
  }
}




export type OfferFormState = {
  errors?: {
    [key: string]: string[] | undefined;
    manufacturer?: string[];
    model?: string[];
    color?: string[];
    fuel_type?: string[];
    transmission?: string[];
    drive?: string[];
    production_year?: string[];
    mileage?: string[];
    number_of_doors?: string[];
    number_of_seats?: string[];
    number_if_gears?: string[];
    engine_power?: string[];
    registration_date?: string[];
    registration_number?: string[];
    vin?: string[];
    engine_capacity?: string[];
    price?: string[];
    is_auction?: string[];
    date_end?: string[];
    description?: string[];
    images?: string[];
    buy_now_price?: string[];
    upload_offer?: string[];
  }
  values?: {
    [key: string]: string | number | boolean | File[] | undefined;
    manufacturer?: string;
    model?: string;
    color?: string;
    fuel_type?: string;
    transmission?: string;
    drive?: string;
    production_year?: number;
    mileage?: number;
    number_of_doors?: number;
    number_of_seats?: number;
    number_of_gears?: number;
    engine_power?: number;
    engine_capacity?: number;
    registration_date?: string;
    registration_number?: string;
    vin?: string;
    description?: string;
    images?: File[];
    price?: number;
    margin?: number;
    is_auction?: boolean;
    date_end?: string;
    buy_now_price?: number;
  }
}

export type OfferFormData = {
  producers: string[];
  models: string[][];
  colors: string[];
  fuelTypes: string[];
  gearboxes: string[];
  driveTypes: string[];
}

export type RegularOfferData = {
  manufacturer: string;
  model: string;
  color: string;
  fuel_type: string;
  transmission: string;
  drive: string;
  production_year: number;
  mileage: number;
  number_of_doors: number;
  number_of_seats: number;
  number_of_gears: number;
  engine_power: number;
  registration_date: string;
  registration_number: string;
  vin: string;
  engine_capacity: number;
  description: string;
  price: number;
  margin: number;
}

export type AuctionOfferData = RegularOfferData & {
  buy_now_price?: number;
  date_end: string;
}