import { FilterFieldData, RangeFieldData, SaleOffer, SaleOfferDetails, SearchParams } from "./definitions";

// Sorting

export async function fetchSortingOptions() : Promise<string[]> {


    await new Promise(resolve => setTimeout(resolve, 1000));

    // TODO connect API and add Suspense
    const data: string[] = ["Base", "name", "date", "price"];

    return data;
}

// Filters

async function fetchProducers() : Promise<FilterFieldData> {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Producer",
        options: ["Audi", "Volksvagen", "Porsche", "Toyota", "Honda", "Mercedes", "Renault"]
    };

    return data;
}

async function fetchGearboxes() : Promise<FilterFieldData> {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Gearbox",
        options: ["Manual", "Sequential Manual", "Automatic"]
    };

    return data;
}

async function fetchFuelTypes() : Promise<FilterFieldData> {
    // TODO connect API

    const data: FilterFieldData = {
        fieldName: "Fuel",
        options: ["Diesel", "Electric", "Gasoline"]
    };

    return data;
}

export async function fetchFilterFields() : Promise<FilterFieldData[]> {

    await new Promise(resolve => setTimeout(resolve, 1000));

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
        }
    ];

    return data;
}

// Home page
export async function fetchTotalPages(params: SearchParams) : Promise<number> {

    // TODO connect API
    return 10;
}

export async function fetchTotalOffers(params: SearchParams) : Promise<number> {


    await new Promise(resolve => setTimeout(resolve, 1000));

    // TODO connect API
    return 100;
}

export async function fetchOffers(params: SearchParams) : Promise<SaleOffer[]> {

    await new Promise(resolve => setTimeout(resolve, 1000));

    // TODO connect API

    const data: SaleOffer[] = [
        {
            id: "1",
            name: "Audi A4",
            productionYear: 2000,
            mileage: 150000,
            color: "Green",
            price: 10000,
            isAuction: true,
        },
        {
            id: "2",
            name: "Volkswagen Golf",
            productionYear: 2005,
            mileage: 120000,
            color: "Blue",
            price: 15000,
            isAuction: false,
        },
        {
            id: "3",
            name: "Porsche 911",
            productionYear: 2010,
            mileage: 80000,
            color: "Red",
            price: 50000,
            isAuction: true,
        }
    ];
    return data;
}

export async function fetchHomePageData(params: SearchParams) : Promise<{totalPages: number, totalOffers: number, offers: SaleOffer[]}> {


    try{

        const totalPages = await fetchTotalPages(params);
        const totalOffers = await fetchTotalOffers(params);
        const offers = await fetchOffers(params);

        const [totalPagesResult, totalOffersResult, offersResult] = await Promise.all([
            totalPages,
            totalOffers,
            offers
        ]);
        return {
            totalPages: totalPagesResult,
            totalOffers: totalOffersResult,
            offers: offersResult
        };

    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch home page data.');
    }
}

// Offer page
export async function fetchOfferDetails(id: string) : Promise<SaleOfferDetails> {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    const data: SaleOfferDetails = {
        name: "Volkswagen Golf",
        imagesURLs: [
            "http://localhost:8081/test_car_image_1.webp",
            "http://localhost:8081/test_car_image_2.webp",
            "http://localhost:808/test_car_image_3.webp",
        ] // temp mock
    };
    return data;
}