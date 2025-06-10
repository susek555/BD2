import { SaleOffer } from "@/app/lib/definitions/SaleOffer";
import { SearchParams } from "@/app/lib/definitions/SearchParams";
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

// export async function fetchOffers(params: SearchParams): Promise<SaleOffer[]> {

//     await new Promise(resolve => setTimeout(resolve, 1000));

//   // TODO connect API

//   const data: SaleOffer[] = [
//     {
//       id: '1',
//       name: 'Audi A4',
//       production_year: 2000,
//       mileage: 150000,
//       color: 'Green',
//       price: 10000,
//       is_auction: true,
//       is_liked: true,
//     },
//     {
//       id: '2',
//       name: 'Volkswagen Golf',
//       production_year: 2005,
//       mileage: 120000,
//       color: 'Blue',
//       price: 15000,
//       is_auction: false,
//       is_liked: false,
//     },
//     {
//       id: '3',
//       name: 'Porsche 911',
//       production_year: 2010,
//       mileage: 80000,
//       color: 'Red',
//       price: 50000,
//       is_auction: true,
//       is_liked: true,
//     },
//   ];
//   return data;
// }

export async function fetchHomePageData(params: SearchParams) : Promise<{totalPages: number, totalOffers: number, offers: SaleOffer[]}> {
    try{
        const data = await getHomePageData(params);

        console.log("Fetched offers", data.offers);

        return {
            totalPages: data.pagination.total_pages,
            totalOffers: data.pagination.total_records,
            offers: data.offers
        };

    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch home page data.');
    }
}