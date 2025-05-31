import { OfferFormState } from "./definitions/offer-form";
import { OfferDetailsFormSchema, CombinedOfferFormSchema, CombinedImagesOfferFormSchema } from "./definitions/offer-form";


export const OfferFormEnum = {
  initialState: 0,
  pricingPart: 1,
  imagesPart: 2,
  readyToApi: 3,
} as const;

export const parseOfferForm = (formData: FormData, progressState: number) : {progressState: number, offerFormState: OfferFormState} => {

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
      buy_now_price: formDataObj.buy_now_price
      ? parseInt(formDataObj.buy_now_price as string) || undefined
      : undefined,
      images: formData.getAll('images').filter(file => file instanceof File)
    }

    let validatedFields;
    switch (progressState) {
      case OfferFormEnum.initialState:
        validatedFields = OfferDetailsFormSchema.safeParse(normalizedData);
      break;
      case OfferFormEnum.pricingPart:
        validatedFields = CombinedOfferFormSchema.safeParse(normalizedData);
      break;
      case OfferFormEnum.imagesPart:
        validatedFields = CombinedImagesOfferFormSchema.safeParse(normalizedData);
      break;
      default:
        throw new Error(`Invalid progress state: ${progressState}`);
    }



    if (!validatedFields.success) {
      console.log("Add Offer validation errors:", validatedFields.error.flatten().fieldErrors);
      return {
        progressState: progressState,
        offerFormState: {
        errors: validatedFields.error.flatten().fieldErrors,
        values: normalizedData as OfferFormState['values']
        }
      };
    }

    if (progressState === OfferFormEnum.imagesPart) {
      // convert model from "producer model" to "model"
      if ('model' in validatedFields.data) {
          const firstSpaceIndex = validatedFields.data.model!.indexOf(' ');
          if (firstSpaceIndex !== -1) {
              validatedFields.data.model! = validatedFields.data.model!.substring(firstSpaceIndex + 1);
          }
      console.log("Sanitized model:", validatedFields.data.model);
      }
    }

    return {
      progressState: progressState + 1,
      offerFormState: {
      errors: {},
      values: validatedFields.data as OfferFormState['values']
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

