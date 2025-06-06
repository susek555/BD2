import { fetchListings } from '@/app/lib/data/account/data';
import { parseFiltersParams } from '@/app/lib/definitions/SearchParams';
import FoundInfo from '@/app/ui/(common)/found-info';
import SideBar from '@/app/ui/(filters-sidebar)/sidebar';
import Pagination from '@/app/ui/(offers-table)/pagination';
import UsersListingsTable from '@/app/ui/account/listings/user-listings-table';
import {
  OffersFoundSkeleton,
  OffersTableSkeleton,
  SideBarSkeleton,
} from '@/app/ui/skeletons';
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

  const params = parseFiltersParams(searchParams);

  console.log('Search Params:', params);

  const { totalPages, totalOffers, offers } = await fetchListings(params);

  return (
    <div className='flex flex-grow flex-col md:flex-row'>
      <div className='h-full w-full flex-none py-4 md:w-80'>
        <Suspense fallback={<SideBarSkeleton />}>
          <SideBar />
        </Suspense>
      </div>
      <div className='flex-grow p-6 md:px-12 md:py-8'>
        <div className='flex w-full items-center justify-between'>
          <Suspense fallback={<OffersFoundSkeleton />}>
            <FoundInfo title={'Offers found'} totalOffers={totalOffers} />
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
          <UsersListingsTable offers={offers} />
        </Suspense>
        {totalOffers !== 0 && (
          <div className='mt-5 flex w-full justify-center'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        )}
      </div>
    </div>
  );
}
