import { getColors, getDrives, getFuelTypes, getTransmissions } from "@/app/lib/api/filters";
import { getOfferTypes } from "@/app/lib//api/offerType";
import { getOrderKeys } from "@/app/lib//api/orderKeys";
import { FilterFieldData, ModelFieldData, RangeFieldData,AddOfferFormData } from "@/app/lib//definitions";

// Offer types

export async function fetchOfferTypes() : Promise<string[]> {
    let data = await getOfferTypes();

    data = ["All", ...data];
    console.log("Offer types data: ", data);

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

    console.log("Gearboxes data: ", data);

    return data;
}

async function fetchFuelTypes() : Promise<string[]> {
    const data = await getFuelTypes();

    console.log("Fuel types data: ", data);

  return data;
}

async function fetchColors() : Promise<string[]> {
    const data = await getColors();

    console.log("Colors data: ", data);

    return data;
}

async function fetchDriveTypes() : Promise<string[]> {
    const data = await getDrives();

    console.log("Drive types data: ", data);

    return data;
}

async function fetchCountries() : Promise<string[]> {
    // TODO connect API

    const data = ["Germany", "France", "Italy", "Spain", "USA"]

    return data;
}



export async function fetchProducersAndModels() : Promise<ModelFieldData> {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    type ProducersAndModelsResponse = {
        producers: string[];
        models: string[][];
    }

    const response: ProducersAndModelsResponse = {
        producers: ["Audi", "Volkswagen", "Porsche"],
        models: [
            ["A4", "A5", "A6"],
            ["Golf", "Passat", "Tiguan"],
            ["911", "Cayenne"]
        ]
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
            (async () => ({
                fieldName: "Countries",
                options: await fetchCountries(),
            }))()
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

export async function fetchAddOfferFormData() : Promise<AddOfferFormData> {

    // TODO remove delay
    await new Promise(resolve => setTimeout(resolve, 1000));

    const producersAndModels = fetchProducersAndModels();
    const colors = fetchColors();
    const fuelTypes = fetchFuelTypes();
    const gearboxes = fetchGearboxes();
    const driveTypes = fetchDriveTypes();
    const countries = fetchCountries();

    const [ producersAndModelsResult,
            colorsResult,
            fuelTypesResult,
            gearboxesResult,
            driveTypesResult,
            countriesResult
        ] = await Promise.all([
            producersAndModels,
            colors,
            fuelTypes,
            gearboxes,
            driveTypes,
            countries
        ]);

    return {
        producers: producersAndModelsResult.producers.options,
        models: producersAndModelsResult.models,
        colors: colorsResult,
        fuelTypes: fuelTypesResult,
        gearboxes: gearboxesResult,
        driveTypes: driveTypesResult,
        countries: countriesResult
    };
}