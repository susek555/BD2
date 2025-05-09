'use client';

import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import SideBar from "@/app/ui/(home)/sidebar";
import { OffersFoundSkeleton, OffersTableSkeleton } from "../ui/skeletons";
import Pagination from "../ui/(home)/pagination";
import OffersFoundInfo from "../ui/(home)/offers-found-info";
import OffersTable from "../ui/(home)/table";
import { fetchHomePageData } from "../lib/data";
import { SearchParams } from "../lib/definitions";
import { SaleOffer } from "@/app/lib/definitions";

export default function Home() {
  const searchParams = useSearchParams();
  const [totalPages, setTotalPages] = useState(0);
  const [totalOffers, setTotalOffers] = useState(0);
  const [offers, setOffers] = useState<SaleOffer[]>([]);
  const [loading, setLoading] = useState(true);  // Loading state

  const params: SearchParams = {
    query: searchParams.get("query") || null,
    page: searchParams.get("page") ? parseInt(searchParams.get("page") || "1", 10) : null,
    orderKey: searchParams.get("sortKey") || null,
    isOrderDesc: searchParams.get("isSortDesc") === "true" ? true : null,
    producers: searchParams.getAll("Producer") || null,
    gearboxes: searchParams.getAll("Gearbox") || null,
    fuelTypes: searchParams.getAll("Fuel") || null,
    price: {
      min: searchParams.get("Price_min") ? parseInt(searchParams.get("Price_min") || "0", 10) : null,
      max: searchParams.get("Price_max") ? parseInt(searchParams.get("Price_max") || "0", 10) : null,
    },
    mileage: {
      min: searchParams.get("Mileage_min") ? parseInt(searchParams.get("Mileage_min") || "0", 10) : null,
      max: searchParams.get("Mileage_max") ? parseInt(searchParams.get("Mileage_max") || "0", 10) : null,
    },
    year: {
      min: searchParams.get("Production year_min") ? parseInt(searchParams.get("Production year_min") || "0", 10) : null,
      max: searchParams.get("Production year_max") ? parseInt(searchParams.get("Production year_max") || "0", 10) : null,
    },
  };

  console.log("Search Params:", params);

  useEffect(() => {
    async function fetchData() {
      setLoading(true);  // Start loading
      try {
        const { totalPages, totalOffers, offers } = await fetchHomePageData(params);
        setTotalPages(totalPages);
        setTotalOffers(totalOffers);
        setOffers(offers);
      } catch (error) {
        console.error("Error fetching data:", error);
      } finally {
        setLoading(false);  // End loading
      }
    }
    fetchData();
  }, [searchParams]);

  return (
    <main>
      <div className="flex flex-col md:flex-row flex-grow">
        <div className="w-full md:w-80 py-4 h-full flex-none">
          <Suspense >
            <SideBar />
          </Suspense>
        </div>
        <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
          {loading ? (
            <OffersFoundSkeleton />
          ) : (
            <OffersFoundInfo totalOffers={totalOffers} />
          )}
          <div className="my-4 " />
          <div className="flex flex-col gap-4">
            {loading ? (
                <OffersTableSkeleton />
            ) : (
              offers.map((offer) => (
                <OffersTable key={offer.name} {...offer} />
              ))
            )}
          </div>
          <div className="mt-5 flex w-full justify-center pr-20">
            <Pagination totalPages={totalPages} />
          </div>
        </div>
      </div>
    </main>
  );
}
