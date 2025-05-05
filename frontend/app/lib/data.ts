import { FilterFieldData } from "./definitions";

export async function fetchProducers() {
    // TODO connect API

    const data: FilterFieldData = {
        name: "Producer",
        options: ["Audi", "Volksvagen", "Porsche", "Toyota", "Honda", "Mercedes", "Renault"]
    };

    return data;
}

export async function fetchGearboxes() {
    // TODO connect API

    const data: FilterFieldData = {
        name: "Gearbox",
        options: ["Manual", "Sequential Manual", "Automatic"]
    };

    return data;
}

export async function fetchFuelTypes() {
    // TODO connect API

    const data: FilterFieldData = {
        name: "Fuel type",
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