import { SaleOffer, SearchParams } from "@/app/lib/definitions";
import { getHomePageData } from "@/app/lib/api/(home)/homePageData";

// Home page

// export async function fetchTotalPages(params: SearchParams): Promise<number> {

//   // TODO connect API
//   return 10;
// }

// export async function fetchTotalOffers(params: SearchParams): Promise<number> {


//     await new Promise(resolve => setTimeout(resolve, 1000));

//   // TODO connect API
//   return 100;
// }

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

export async function fetchHomePageData(params: SearchParams) : Promise<{totalPages: number, totalOffers: number, offers: SaleOffer[]}> {
    try{
        const data = await getHomePageData(params);

        return {
            totalPages: data.pagination.total_pages,
            totalOffers: data.pagination.total_records,
            // offers: data.offers
            //TODO - remove mock data
            offers: await fetchOffers(params)
        };

    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch home page data.');
    }
}