import { OfferFormState } from "./definitions/offer-form";
import { OfferDetailsFormSchema, CombinedOfferFormSchema } from "./definitions/offer-form";


export const OfferFormEnum = {
  readyToApi: true,
  pricingPartLeft: false
} as const;

export const parseOfferForm = (formData: FormData, detailsPart: boolean) : {boolean: boolean, offerFormState: OfferFormState} => {

    const formDataObj = Object.fromEntries(formData.entries());

    const normalizedData = {
        ...formDataObj,
        production_year: formDataObj.production_year ? parseInt(formDataObj.production_year as string) || undefined : undefined,
        mileage: formDataObj.mileage ? parseInt(formDataObj.mileage as string) || undefined : undefined,
        number_of_doors: formDataObj.number_of_doors ? parseInt(formDataObj.number_of_doors as string) || undefined : undefined,
        number_of_seats: formDataObj.number_of_seats ? parseInt(formDataObj.number_of_seats as string) || undefined : undefined,
        number_of_gears: formDataObj.number_of_gears ? parseInt(formDataObj.number_of_gears as string) || undefined : undefined,
        engine_power: formDataObj.engine_power ? parseInt(formDataObj.engine_power as string) || undefined : undefined,
        engine_capacity: formDataObj.engine_capacity ? parseInt(formDataObj.engine_capacity as string) || undefined : undefined,
        price: formDataObj.price ? parseInt(formDataObj.price as string) || undefined : undefined,
        margin: formDataObj.margin
            ? (parseInt((formDataObj.margin as string).replace('%', '')) || undefined)
            : undefined,
        is_auction: formDataObj.is_auction === "true" ? true : false,
        buy_now_auction_price: formDataObj.buy_now_auction_price
            ? parseInt(formDataObj.buy_now_auction_price as string) || undefined
            : undefined
    }

    const validatedFields = detailsPart
        ? OfferDetailsFormSchema.safeParse(normalizedData)
        : CombinedOfferFormSchema.safeParse(normalizedData);;



    if (!validatedFields.success) {
      console.log("Add Offer validation errors:", validatedFields.error.flatten().fieldErrors);
      return {
        boolean: OfferFormEnum.pricingPartLeft,
        offerFormState: {
        errors: validatedFields.error.flatten().fieldErrors,
        values: normalizedData as OfferFormState['values']
        }
      };
    }

    if (!detailsPart) {
      // convert model from "producer model" to "model"
      if ('model' in validatedFields) {
          const firstSpaceIndex = validatedFields.data.model!.indexOf(' ');
          if (firstSpaceIndex !== -1) {
              validatedFields.model! = validatedFields.data.model!.substring(firstSpaceIndex + 1);
          }
      }
    }

    return {
      boolean: detailsPart ? OfferFormEnum.pricingPartLeft : OfferFormEnum.readyToApi,
      offerFormState: {
      errors: {},
      values: normalizedData as OfferFormState['values']
      }
    };
}



export const generatePagination = (currentPage: number, totalPages: number) => {
    // If the total number of pages is 7 or less,
    // display all pages without any ellipsis.
    if (totalPages <= 7) {
      return Array.from({ length: totalPages }, (_, i) => i + 1);
    }

    // If the current page is among the first 3 pages,
    // show the first 3, an ellipsis, and the last 2 pages.
    if (currentPage <= 3) {
      return [1, 2, 3, '...', totalPages - 1, totalPages];
    }

    // If the current page is among the last 3 pages,
    // show the first 2, an ellipsis, and the last 3 pages.
    if (currentPage >= totalPages - 2) {
      return [1, 2, '...', totalPages - 2, totalPages - 1, totalPages];
    }

    // If the current page is somewhere in the middle,
    // show the first page, an ellipsis, the current page and its neighbors,
    // another ellipsis, and the last page.
    return [
      1,
      '...',
      currentPage - 1,
      currentPage,
      currentPage + 1,
      '...',
      totalPages,
    ];
  };

