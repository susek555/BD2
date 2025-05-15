import { fetchTotalPages } from '@/app/lib/data';
import {
  parseArrayOrUndefined,
  parseIntOrUndefined,
  SearchParams,
} from '@/app/lib/definitions';
import UsersListings from '@/app/ui/account/listings/user-listings-table';
import OffersFoundInfo from '@/app/ui/offers-found-info';
import Pagination from '@/app/ui/pagination';
import SideBar from '@/app/ui/sidebar';
import { OffersFoundSkeleton, OffersTableSkeleton } from '@/app/ui/skeletons';
import { PlusIcon } from '@heroicons/react/20/solid';
import Link from 'next/link';
import { Suspense } from 'react';

export default async function ListingsPage(props: {
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
        <div className='flex w-full items-center justify-between'>
          <Suspense fallback={<OffersFoundSkeleton />}>
            <OffersFoundInfo params={params} />
          </Suspense>
          <Link
            href='/listing/create' // TODO add actual creeate redirect
            className='flex items-center rounded-md bg-green-500 px-4 py-2 text-white hover:bg-green-600'
          >
            Add listing <PlusIcon className='ml-2 h-5 w-5' />
          </Link>
        </div>

        <div className='my-4' />
        <Suspense fallback={<OffersTableSkeleton />}>
          <UsersListings params={params} />
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
