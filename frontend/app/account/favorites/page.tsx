import { fetchFavorites } from '@/app/lib/data/account/data';
import {
  parseFiltersParams,
} from '@/app/lib/definitions/SearchParams';
import FoundInfo from '@/app/ui/(common)/found-info';
import OffersTable from '@/app/ui/(offers-table)/offers-table';
import Pagination from '@/app/ui/(offers-table)/pagination';
import SideBar from '@/app/ui/(filters-sidebar)/sidebar';
import { OffersFoundSkeleton, OffersTableSkeleton, SideBarSkeleton } from '@/app/ui/skeletons';
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

  const params = parseFiltersParams(searchParams);

  console.log('Search Params:', params);

  const { totalPages, totalOffers, offers } = await fetchFavorites(params);

  return (
    <div className='flex flex-grow flex-col md:flex-row'>
      <div className='h-full w-full flex-none py-4 md:w-80'>
        <Suspense fallback={<SideBarSkeleton />}>
          <SideBar />
        </Suspense>
      </div>
      <div className='flex-grow p-6 md:px-12 md:py-8'>
        <Suspense fallback={<OffersFoundSkeleton />}>
          <FoundInfo title={"Offers found"} totalOffers={totalOffers} />
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
