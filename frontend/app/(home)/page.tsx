import { fetchHomePageData } from '@/app/lib/data/(home)/data';
import SideBar from '@/app/ui/(filters-sidebar)/sidebar';
import { Suspense } from 'react';
import { parseFiltersParams } from '../lib/definitions/SearchParams';
import FoundInfo from '../ui/(common)/found-info';
import OffersTable from '../ui/(offers-table)/offers-table';
import Pagination from '../ui/(offers-table)/pagination';
import {
  OffersFoundSkeleton,
  OffersTableSkeleton,
  SideBarSkeleton,
} from '../ui/skeletons';

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

  const params = parseFiltersParams(searchParams);

  console.log('Search Params:', params);

  const { totalPages, totalOffers, offers } = await fetchHomePageData(params);

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
            <FoundInfo title={'Offers found'} totalOffers={totalOffers} />
          </Suspense>
          <div className='my-4' />
          <Suspense fallback={<OffersTableSkeleton />}>
            <OffersTable offers={offers} />
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
    </main>
  );
}
