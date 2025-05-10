import Filters from '@/app/ui/(home)/filters';
import Sorting from './sorting';
import OfferType from './offer-type';
import { fetchFilterFields, fetchSortingOptions, prepareRangeFields } from '@/app/lib/data';


export default async function SideBar() {
    const filters = await fetchFilterFields();
    const ranges = await prepareRangeFields();
    const sortingOptions = await fetchSortingOptions();

    return (
      <div className="flex h-full flex-col px-3 py-2 md:px-2 rounded-r-lg border-black border-[2px]">
        <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
          <OfferType />
          <Sorting sortingOptions={sortingOptions}/>
          <Filters filtersData={filters} rangesData={ranges}/>
          <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
        </div>
      </div>
    );
}