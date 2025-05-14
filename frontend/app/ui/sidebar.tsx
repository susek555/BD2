import {
  fetchFilterFields,
  fetchSortingOptions,
  prepareRangeFields,
} from '@/app/lib/data';
import Filters from '@/app/ui/(home)/filters';
import OfferType from './offer-type';
import Sorting from './sorting';

export default async function SideBar() {
  const filters = await fetchFilterFields();
  const ranges = await prepareRangeFields();
  const sortingOptions = await fetchSortingOptions();

  return (
    <div className='flex h-full flex-col rounded-lg border-2 border-gray-300 bg-white p-4'>
      <div className='flex flex-col space-y-4'>
        <OfferType />
        <Sorting sortingOptions={sortingOptions} />
        <Filters filtersData={filters} rangesData={ranges} />
      </div>
    </div>
  );
}
