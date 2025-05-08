'use client';

import SideBar from "@/app/ui/(home)/sidebar";
import { Suspense, useEffect, useState } from "react";
import { OffersFoundSkeleton, OffersTableSkeleton } from "../ui/skeletons";
import Pagination from "../ui/(home)/pagination";
import OffersFoundInfo from "../ui/(home)/offers-found-info";
import OffersTable from "../ui/(home)/table";
import { fetchHomePageData } from "../lib/data";
import { useSearchParams } from "next/navigation";
import { SearchParams } from "../lib/definitions";
import { SaleOffer } from "../lib/definitions";


export default function Home() {
  const searchParams = useSearchParams();
  const [totalPages, setTotalPages] = useState(0);
  const [totalOffers, setTotalOffers] = useState(0);
  const [offers, setOffers] = useState<SaleOffer[]>([]);

  const params: SearchParams = {
    query: searchParams.get("query") || null,
    page: searchParams.get("page") ? parseInt(searchParams.get("page") || "1", 10) : null,
    producers: searchParams.getAll("producers") || null,
    gearboxes: searchParams.getAll("gearboxes") || null,
    fuelTypes: searchParams.getAll("fuelTypes") || null,
    price: {
      min: searchParams.get("price[min]") ? parseInt(searchParams.get("price[min]") || "0", 10) : null,
      max: searchParams.get("price[max]") ? parseInt(searchParams.get("price[max]") || "0", 10) : null,
    },
    mileage: {
      min: searchParams.get("mileage[min]") ? parseInt(searchParams.get("mileage[min]") || "0", 10) : null,
      max: searchParams.get("mileage[max]") ? parseInt(searchParams.get("mileage[max]") || "0", 10) : null,
    },
    year: {
      min: searchParams.get("year[min]") ? parseInt(searchParams.get("year[min]") || "0", 10) : null,
      max: searchParams.get("year[max]") ? parseInt(searchParams.get("year[max]") || "0", 10) : null,
    },
  };

  useEffect(() => {
    async function fetchData() {
      try
      {
        const { totalPages, totalOffers, offers } = await fetchHomePageData(params);
        setTotalPages(totalPages);
        setTotalOffers(totalOffers);
        setOffers(offers);
      }
      catch (error) {
        console.error("Error fetching data:", error);
      }
    }
    fetchData();
  }, [searchParams]);

  return (
    <main>
      <div className="flex flex-col md:flex-row flex-grow">
        <div className="w-full md:w-80 py-4 h-full flex-none">
          <Suspense>
            <SideBar />
          </Suspense>
        </div>
        <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
          <Suspense key={totalOffers} fallback={<OffersFoundSkeleton />}>
            <OffersFoundInfo totalOffers={totalOffers} />
          </Suspense>
          <div className="my-4 " />
          <div className="flex flex-col gap-4">
            {offers.map((offer) => (
            <Suspense key={offer.name} fallback={<OffersTableSkeleton />}>
              <OffersTable {...offer} />
            </Suspense>
            ))}
          </div>
          <div className="mt-5 flex w-full justify-center pr-20">
            <Pagination totalPages={totalPages} />
          </div>
        </div>
      </div>
    </main>
  );
}
