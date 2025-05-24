import { SaleOffer } from "@/app/lib/definitions/SaleOffer";
import { SearchParams } from "@/app/lib/definitions/SearchParams";

const API_URL = process.env.API_URL;

export async function getHomePageData(params : SearchParams) : Promise<

    {
        pagination: {
            total_pages: number;
            total_records: number;
        };
        offers: SaleOffer[];
    }
> {
  const response = await fetch(`${API_URL}/sale-offer/filtered`, {
    method: "POST",
    headers: {
    'Content-Type': 'application/json',
    },
    body: JSON.stringify(params),
  });

  if (!response.ok) {
    throw new Error("Failed to fetch home page data");
  }

  const data = await response.json();

  console.log("Home page data: ", data);

  return data;
}