import { fetchFavorites } from '@/app/lib/account/data';
import {
  parseArrayOrUndefined,
  parseIntOrUndefined,
  trimAllAfterFirstSpace,
  SearchParams,
} from '@/app/lib/definitions';
import OffersFoundInfo from '@/app/ui/offers-found-info';
import OffersTable from '@/app/ui/offers-table';
import Pagination from '@/app/ui/pagination';
import SideBar from '@/app/ui/sidebar';
import { OffersFoundSkeleton, OffersTableSkeleton } from '@/app/ui/skeletons';
import { Suspense } from 'react';

export default async function ActivityPage(props: {
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

  const params: SearchParams = {
    query: searchParams?.query,
    pagination: {
      page: searchParams?.page ? parseInt(searchParams.page, 10) : 1,
      page_size: 6,
    },
    order_key: searchParams?.sortKey,
    is_order_desc: searchParams?.isSortDesc === "true" ? true : undefined,
    offer_type: searchParams?.offerType,
    manufacturers: parseArrayOrUndefined(searchParams?.Producers),
    models: trimAllAfterFirstSpace(parseArrayOrUndefined(searchParams?.Models)),
    colors: parseArrayOrUndefined(searchParams?.Colors),
    transmissions: parseArrayOrUndefined(searchParams?.Gearboxes),
    fuel_types: parseArrayOrUndefined(searchParams?.Fueltypes),
    drives: parseArrayOrUndefined(searchParams?.Drivetypes),
    price_range: parseIntOrUndefined(searchParams?.Price_min) || parseIntOrUndefined(searchParams?.Price_max) ? {
      min: parseIntOrUndefined(searchParams?.Price_min),
      max: parseIntOrUndefined(searchParams?.Price_max),
    } : undefined,
    mileage_range: parseIntOrUndefined(searchParams?.Mileage_min) || parseIntOrUndefined(searchParams?.Mileage_max) ? {
      min: parseIntOrUndefined(searchParams?.Mileage_min),
      max: parseIntOrUndefined(searchParams?.Mileage_max),
    } : undefined,
    year_range: parseIntOrUndefined(searchParams?.Productionyear_min) || parseIntOrUndefined(searchParams?.Productionyear_max) ? {
      min: parseIntOrUndefined(searchParams?.Productionyear_min),
      max: parseIntOrUndefined(searchParams?.Productionyear_max),
    } : undefined,
    engine_capacity_range: parseIntOrUndefined(searchParams?.Enginecapacity_min) || parseIntOrUndefined(searchParams?.Enginecapacity_max) ? {
      min: parseIntOrUndefined(searchParams?.Enginecapacity_min),
      max: parseIntOrUndefined(searchParams?.Enginecapacity_max),
    } : undefined,
    engine_power_range: parseIntOrUndefined(searchParams?.Enginepower_min) || parseIntOrUndefined(searchParams?.Enginepower_max) ? {
      min: parseIntOrUndefined(searchParams?.Enginepower_min),
      max: parseIntOrUndefined(searchParams?.Enginepower_max),
    } : undefined,
  };

  console.log('Search Params:', params);

  const { totalPages, totalOffers, offers } = await fetchFavorites(params);

  return (
    <div className='flex flex-grow flex-col md:flex-row'>
      <div className='h-full w-full flex-none py-4 md:w-80'>
        <Suspense>
          <SideBar />
        </Suspense>
      </div>
      <div className='flex-grow p-6 md:px-12 md:py-8'>
        <Suspense fallback={<OffersFoundSkeleton />}>
          <OffersFoundInfo totalOffers={totalOffers} />
        </Suspense>
        <div className='my-4' />
        <Suspense fallback={<OffersTableSkeleton />}>
          <OffersTable offers={offers} />
        </Suspense>
        <div className='mt-5 flex w-full justify-center'>
          <Suspense>
            <Pagination totalPages={totalPages} />
          </Suspense>
        </div>
      </div>
    </div>
  );
}
