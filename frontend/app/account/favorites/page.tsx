import { fetchFavorites } from '@/app/lib/account/data';
import { fetchTotalPages } from '@/app/lib/data';
import {
  parseArrayOrUndefined,
  parseIntOrUndefined,
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
        query: searchParams.query,
        page: parseIntOrUndefined(searchParams.page),
        orderKey: searchParams.sortKey,
        isOrderDesc: searchParams.isSortDesc === 'true',
        offerType: searchParams.offerType,
        producers: parseArrayOrUndefined(searchParams.Producer),
        gearboxes: parseArrayOrUndefined(searchParams.Gearbox),
        fuelTypes: parseArrayOrUndefined(searchParams.Fuel),
        price: {
          min: parseIntOrUndefined(searchParams.Price_min),
          max: parseIntOrUndefined(searchParams.Price_max),
        },
        mileage: {
          min: parseIntOrUndefined(searchParams.Mileage_min),
          max: parseIntOrUndefined(searchParams.Mileage_max),
        },
        year: {
          min: parseIntOrUndefined(searchParams.Productionyear_min),
          max: parseIntOrUndefined(searchParams.Productionyear_max),
        },
      }
    : {};

  console.log('Search Params:', params);

  const totalPages = await fetchTotalPages(params);

  return (
    <div className='flex flex-grow flex-col md:flex-row'>
      <div className='h-full w-full flex-none py-4 md:w-80'>
        <Suspense>
          <SideBar />
        </Suspense>
      </div>
      <div className='flex-grow p-6 md:px-12 md:py-8'>
        <Suspense fallback={<OffersFoundSkeleton />}>
          <OffersFoundInfo params={params} />
        </Suspense>
        <div className='my-4' />
        <Suspense fallback={<OffersTableSkeleton />}>
          <OffersTable params={params} fetchFunction={fetchFavorites} />
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
