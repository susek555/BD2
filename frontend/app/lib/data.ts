import { FilterFieldData, RangeFieldData } from "./definitions";

// Filters

export async function fetchProducers() {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Producer",
        options: ["Audi", "Volksvagen", "Porsche", "Toyota", "Honda", "Mercedes", "Renault"]
    };

    return data;
}

export async function fetchGearboxes() {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Gearbox",
        options: ["Manual", "Sequential Manual", "Automatic"]
    };

    return data;
}

export async function fetchFuelTypes() {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Fuel",
        options: ["Diesel", "Electric", "Gasoline"]
    };

    return data;
}

export async function fetchFilterFields() {
    try{
        const producersData = fetchProducers();
        const gearboxesData = fetchGearboxes();
        const fuelTypesData = fetchFuelTypes();

        const data = await Promise.all([
            producersData,
            gearboxesData,
            fuelTypesData
        ]);

        return data;
    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch filters data.');
    }
}

// Ranges

export function prepareRangeFields() {
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
        }
    ];

    return data;
}