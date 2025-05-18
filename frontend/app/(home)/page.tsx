import SideBar from '@/app/ui/sidebar';
import { Suspense } from 'react';
import { fetchOffers, fetchTotalPages } from '../lib/data';
import {
  parseArrayOrUndefined,
  parseIntOrUndefined,
  SearchParams,
} from '../lib/definitions';
import OffersFoundInfo from '../ui/offers-found-info';
import OffersTable from '../ui/offers-table';
import Pagination from '../ui/pagination';
import { OffersFoundSkeleton, OffersTableSkeleton, SideBarSkeleton } from '../ui/skeletons';

//TODO fix bug with filters not syncing when Home button clicked
//TODO handle loading errors

export default async function Home(props: {
  searchParams?: Promise<{
    query?: string;
    page?: string;
    sortKey?: string;
    isSortDesc?: string;
    offerType?: string;
    Producers?: string[];
    Models?: string[];
    Colors?: string[];
    Drivetypes?: string[];
    Gearboxes?: string[];
    Fueltypes?: string[];
    Price_min?: string;
    Price_max?: string;
    Mileage_min?: string;
    Mileage_max?: string;
    Productionyear_min?: string;
    Productionyear_max?: string;
    Enginecapacity_min?: string;
    Enginecapacity_max?: string;
    Enginepower_min?: string;
    Enginepower_max?: string;
  }>;
}) {
  const searchParams = await props.searchParams;

  const params: SearchParams = searchParams
  ? {
      query: searchParams.query,
      pagination: {
        page: searchParams.page ? parseInt(searchParams.page || "1", 10) : null,
        page_size: 6,
      },
      order_key: searchParams.sortKey || null,
      is_order_desc: searchParams.isSortDesc === "true" ? true : null,
      offer_type: searchParams.offerType || null,
      manufacturers: parseArrayOrUndefined(searchParams.Producers),
      models: parseArrayOrUndefined(searchParams.Models),
      colors: parseArrayOrUndefined(searchParams.Colors),
      transmissions: parseArrayOrUndefined(searchParams.Gearboxes),
      fuel_types: parseArrayOrUndefined(searchParams.Fueltypes),
      drives: parseArrayOrUndefined(searchParams.Drivetypes),
      price_range: {
        min: parseIntOrUndefined(searchParams.Price_min),
        max: parseIntOrUndefined(searchParams.Price_max),
      },
      mileage_range: {
        min: parseIntOrUndefined(searchParams.Mileage_min),
        max: parseIntOrUndefined(searchParams.Mileage_max),
      },
      year_range: {
        min: parseIntOrUndefined(searchParams.Productionyear_min),
        max: parseIntOrUndefined(searchParams.Productionyear_max),
      },
      engine_capacity_range: {
        min: parseIntOrUndefined(searchParams.Enginecapacity_min),
        max: parseIntOrUndefined(searchParams.Enginecapacity_max),
      },
      engine_power_range: {
        min: parseIntOrUndefined(searchParams.Enginepower_min),
        max: parseIntOrUndefined(searchParams.Enginepower_max),
      },
    }
  : {
    query: null,
    pagination: {
      page: 1,
      page_size: 6,
    },
    order_key: null,
    is_order_desc: null,
    offer_type: null,
    manufacturers: null,
    models: null,
    colors: null,
    drives: null,
    transmissions: null,
    fuel_types: null,
    price_range: { min: null, max: null },
    mileage_range: { min: null, max: null },
    year_range: { min: null, max: null },
    engine_capacity_range: { min: null, max: null },
    engine_power_range: { min: null, max: null },
    };

  console.log('Search Params:', params);

  const totalPages = await fetchTotalPages(params);

  return (
    <main>
      <div className='flex flex-grow flex-col md:flex-row'>
        <div className='h-full w-full flex-none py-4 md:w-80'>
          <Suspense fallback={<SideBarSkeleton />}>
            <SideBar />
          </Suspense>
        </div>
        <div className='flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8'>
          <Suspense fallback={<OffersFoundSkeleton />}>
            <OffersFoundInfo params={params} />
          </Suspense>
          <div className='my-4' />
          <Suspense fallback={<OffersTableSkeleton />}>
            <OffersTable params={params} fetchFunction={fetchOffers} />
          </Suspense>
          <div className='mt-5 flex w-full justify-center pr-20'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        </div>
      </div>
    </main>
  );
}
