import { getColors, getDrives, getFuelTypes, getProducersAndModels, getTransmissions } from "@/app/lib/api/filters";
import { getOfferTypes } from "@/app/lib//api/offerType";
import { getOrderKeys } from "@/app/lib//api/orderKeys";
import { FilterFieldData, ModelFieldData, RangeFieldData } from "@/app/lib//definitions/filters-sidebar";
import { OfferFormData } from "@/app/lib/definitions/offer-form";

// Offer types

export async function fetchOfferTypes() : Promise<string[]> {
    let data = await getOfferTypes();

    data = ["All", ...data];
    // console.log("Offer types data: ", data);

    return data;
}

// Sorting

export async function fetchSortingOptions() : Promise<string[]> {
    let data = await getOrderKeys();

    data = ["Base", ...data];

  return data;
}

// Filters

async function fetchGearboxes() : Promise<string[]> {
    const data = await getTransmissions();

    // console.log("Gearboxes data: ", data);

    return data;
}

async function fetchFuelTypes() : Promise<string[]> {
    const data = await getFuelTypes();

    // console.log("Fuel types data: ", data);

  return data;
}

async function fetchColors() : Promise<string[]> {
    const data = await getColors();

    // console.log("Colors data: ", data);

    return data;
}

async function fetchDriveTypes() : Promise<string[]> {
    const data = await getDrives();

    // console.log("Drive types data: ", data);

    return data;
}



export async function fetchProducersAndModels() : Promise<ModelFieldData> {
    const receivedData = await getProducersAndModels();

    // console.log("Producers and models data: ", receivedData);

    type ProducersAndModelsResponse = {
        producers: string[];
        models: string[][];
    }

    const response: ProducersAndModelsResponse = {
        producers: receivedData.producers,
        models: receivedData.models
    };

    const data: ModelFieldData = {
        producers: {
            fieldName: "Producers",
            options: response.producers,
        },
        models: response.models
    };

    return data;
}

export async function fetchFilterFields() : Promise<FilterFieldData[]> {
    try{
        const data: FilterFieldData[] = await Promise.all([
            (async () => ({
                fieldName: "Colors",
                options: await fetchColors(),
            }))(),
            (async () => ({
                fieldName: "Gearboxes",
                options: await fetchGearboxes(),
            }))(),
            (async () => ({
                fieldName: "Fuel types",
                options: await fetchFuelTypes(),
            }))(),
            (async () => ({
                fieldName: "Drive types",
                options: await fetchDriveTypes(),
            }))(),
        ]);


        return data;
    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch filters data.');
    }
}

// Ranges

export function prepareRangeFields() : RangeFieldData[] {
    const data: RangeFieldData[] =
    [
        {
            fieldName: "Production year",
            range: {
                min: null,
                max: null
            }
        },
        {
            fieldName: "Mileage",
            range: {
                min: null,
                max: null
            }
        },
        {
            fieldName: "Price",
            range: {
                min: null,
                max: null
            }
        },
        {
            fieldName: "Engine capacity",
            range: {
                min: null,
                max: null
            }
        },
        {
            fieldName: "Engine power",
            range: {
                min: null,
                max: null
            }
        }
    ];

  return data;
}

// Add offer

export async function fetchOfferFormData() : Promise<OfferFormData> {
    const producersAndModels = fetchProducersAndModels();
    const colors = fetchColors();
    const fuelTypes = fetchFuelTypes();
    const gearboxes = fetchGearboxes();
    const driveTypes = fetchDriveTypes();

    const [ producersAndModelsResult,
            colorsResult,
            fuelTypesResult,
            gearboxesResult,
            driveTypesResult,
        ] = await Promise.all([
            producersAndModels,
            colors,
            fuelTypes,
            gearboxes,
            driveTypes,
        ]);

    return {
        producers: producersAndModelsResult.producers.options,
        models: producersAndModelsResult.models,
        colors: colorsResult,
        fuelTypes: fuelTypesResult,
        gearboxes: gearboxesResult,
        driveTypes: driveTypesResult,
    };
}