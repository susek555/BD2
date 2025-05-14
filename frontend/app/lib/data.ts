import { FilterFieldData, RangeFieldData, SaleOffer, SaleOfferDetails, SearchParams } from "./definitions";

// Sorting

export async function fetchSortingOptions(): Promise<string[]> {


    await new Promise(resolve => setTimeout(resolve, 1000));

  // TODO connect API and add Suspense
  const data: string[] = ['Base', 'name', 'date', 'price'];

  return data;
}

// Filters

async function fetchProducers() : Promise<string[]> {
    // TODO connect API

    const data = ["Audi", "Volksvagen", "Porsche", "Toyota", "Honda", "Mercedes", "Renault"]

  return data;
}

async function fetchGearboxes() : Promise<string[]> {
    // TODO connect API

    const data = ["Manual", "Sequential Manual", "Automatic"]

  return data;
}

async function fetchFuelTypes() : Promise<string[]> {
    // TODO connect API

    const data = ["Diesel", "Electric", "Gasoline"]

  return data;
}

export async function fetchFilterFields(): Promise<FilterFieldData[]> {

    await new Promise(resolve => setTimeout(resolve, 1000));

    try{
        const data: FilterFieldData[] = await Promise.all([
            (async () => ({
                fieldName: "Producers",
                options: await fetchProducers(),
            }))(),
            (async () => ({
                fieldName: "Gearboxes",
                options: await fetchGearboxes(),
            }))(),
            (async () => ({
                fieldName: "Fuel types",
                options: await fetchFuelTypes(),
            }))(),
        ]);


        return data;
    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch filters data.');
    }
}

// Ranges

export function prepareRangeFields(): RangeFieldData[] {
  const data: RangeFieldData[] = [
    {
      fieldName: 'Production year',
      range: {
        min: null,
        max: null,
      },
    },
    {
      fieldName: 'Mileage',
      range: {
        min: null,
        max: null,
      },
    },
    {
      fieldName: 'Price',
      range: {
        min: null,
        max: null,
      },
    },
  ];

  return data;
}

// Home page
export async function fetchTotalPages(params: SearchParams): Promise<number> {

  // TODO connect API
  return 10;
}

export async function fetchTotalOffers(params: SearchParams): Promise<number> {


    await new Promise(resolve => setTimeout(resolve, 1000));

  // TODO connect API
  return 100;
}

export async function fetchOffers(params: SearchParams): Promise<SaleOffer[]> {

    await new Promise(resolve => setTimeout(resolve, 1000));

  // TODO connect API

  const data: SaleOffer[] = [
    {
      id: '1',
      name: 'Audi A4',
      productionYear: 2000,
      mileage: 150000,
      color: 'Green',
      price: 10000,
      isAuction: true,
      isFavorite: true,
    },
    {
      id: '2',
      name: 'Volkswagen Golf',
      productionYear: 2005,
      mileage: 120000,
      color: 'Blue',
      price: 15000,
      isAuction: false,
      isFavorite: false,
    },
    {
      id: '3',
      name: 'Porsche 911',
      productionYear: 2010,
      mileage: 80000,
      color: 'Red',
      price: 50000,
      isAuction: true,
      isFavorite: true,
    },
  ];
  return data;
}

export async function fetchHomePageData(
  params: SearchParams,
): Promise<{ totalPages: number; totalOffers: number; offers: SaleOffer[] }> {
  try {
    const totalPages = await fetchTotalPages(params);
    const totalOffers = await fetchTotalOffers(params);
    const offers = await fetchOffers(params);

    const [totalPagesResult, totalOffersResult, offersResult] =
      await Promise.all([totalPages, totalOffers, offers]);
    return {
      totalPages: totalPagesResult,
      totalOffers: totalOffersResult,
      offers: offersResult,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error('Failed to fetch home page data.');
  }
}


// Offer page
export async function fetchOfferDetails(id: string) : Promise<SaleOfferDetails | null> {
    // TODO connect API

    await new Promise(resolve => setTimeout(resolve, 1000));

    let data: SaleOfferDetails | null = null;

    if (id === "2") {
        data = {
            name: "Volkswagen Golf",
            price: 15000,
            isAuction: true,
            isActive: true,
            auctionData: {
            endDate: new Date("2025-07-07T12:00:00Z"),
            currentBid: 12000
            },
            imagesURLs: [
            "http://localhost:8081/test_car_image_1.webp",
            "http://localhost:8081/test_car_image_2.webp",
            "http://localhost:808/test_car_image_3.webp",
            ], // temp mock
            description: "The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology,The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide.The Volkswagen Golf is a compact car known for its reliability, practicality, and efficient performance. It features a comfortable interior, advanced technology, and a timeless design, making it a popular choice among drivers worldwide. and a timeless design, making it a popular choice among drivers worldwide.",
            details: [
                {
                    name: "Engine",
                    value: "2.0 TDI"
                },
                {
                    name: "Power",
                    value: "150 HP"
                },
                {
                    name: "Color",
                    value: "Blue"
                }
            ],
            sellerName: "Zwinny Ambrozy"
        };
    }

    return data;
}
