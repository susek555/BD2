import { Suspense } from "react";
import SideBar from "@/app/ui/(home)/sidebar";
import { OffersFoundSkeleton, OffersTableSkeleton } from "../ui/skeletons";
import Pagination from "../ui/(home)/pagination";
import OffersFoundInfo from "../ui/(home)/offers-found-info";
import OffersTable from "../ui/(home)/table";
import { fetchTotalPages } from "../lib/data";
import { SearchParams } from "../lib/definitions";

//TODO fix bug with filters not syncing when Home button clicked
//TODO handle loading errors

export default async function Home(props: {
  searchParams?: Promise<{
    query?: string;
    page?: string;
    sortKey?: string;
    isSortDesc?: string;
    offerType?: string;
    Producer?: string[];
    Gearbox?: string[];
    Fuel?: string[];
    Price_min?: string;
    Price_max?: string;
    Mileage_min?: string;
    Mileage_max?: string;
    Productionyear_min?: string;
    Productionyear_max?: string;
  }>;
  }) {

  const searchParams = await props.searchParams;

  const params: SearchParams = searchParams
  ? {
    query: searchParams.query || null,
    page: searchParams.page ? parseInt(searchParams.page || "1", 10) : null,
    orderKey: searchParams.sortKey || null,
    isOrderDesc: searchParams.isSortDesc === "true" ? true : null,
    offerType: searchParams.offerType || null,
    producers: searchParams.Producer || null,
    gearboxes: searchParams.Gearbox || null,
    fuelTypes: searchParams.Fuel || null,
    price: {
      min: searchParams.Price_min ? parseInt(searchParams.Price_min || "0", 10) : null,
      max: searchParams.Price_max ? parseInt(searchParams.Price_max || "0", 10) : null,
    },
    mileage: {
      min: searchParams.Mileage_min ? parseInt(searchParams.Mileage_min || "0", 10) : null,
      max: searchParams.Mileage_max ? parseInt(searchParams.Mileage_max || "0", 10) : null,
    },
    year: {
      min: searchParams.Productionyear_min ? parseInt(searchParams.Productionyear_min || "0", 10) : null,
      max: searchParams.Productionyear_max ? parseInt(searchParams.Productionyear_max || "0", 10) : null,
    },
    }
  : {
    query: null,
    page: null,
    orderKey: null,
    isOrderDesc: null,
    offerType: null,
    producers: null,
    gearboxes: null,
    fuelTypes: null,
    price: { min: null, max: null },
    mileage: { min: null, max: null },
    year: { min: null, max: null },
    };

  console.log("Search Params:", params);

  const totalPages = await fetchTotalPages(params);

  return (
  <main>
    <div className="flex flex-col md:flex-row flex-grow">
      <div className="w-full md:w-80 py-4 h-full flex-none">
        <Suspense >
          <SideBar />
        </Suspense>
      </div>
      <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
        <Suspense fallback={<OffersFoundSkeleton />}>
          <OffersFoundInfo params={params} />
        </Suspense>
        <div className="my-4 " />
        <Suspense fallback={<OffersTableSkeleton />}>
          <OffersTable params={params} />
        </Suspense>
        <div className="mt-5 flex w-full justify-center pr-20">
          <Suspense>
            <Pagination totalPages={totalPages} />
          </Suspense>
        </div>
      </div>
    </div>
  </main>
  );
}
