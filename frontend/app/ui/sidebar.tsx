import Filters from '@/app/ui/(home)/filters';
import Sorting from './sorting';
import OfferType from './offer-type';
import { fetchFilterFields, fetchProducersAndModels, fetchSortingOptions, prepareRangeFields } from '@/app/lib/data';
import ProducersAndModels from '@/app/ui/producers-and-models';


export default async function SideBar() {
    const [filters, ranges, sortingOptions, producersAndModels] = await Promise.all([
      fetchFilterFields(),
      prepareRangeFields(),
      fetchSortingOptions(),
      fetchProducersAndModels()
    ]);

    return (
      <div className="flex h-full flex-col px-3 py-2 md:px-2 rounded-r-lg border-black border-[2px]">
        <div className="flex grow flex-col justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
          <OfferType />
          <Sorting sortingOptions={sortingOptions}/>
          <ProducersAndModels producersAndModels={producersAndModels}/>
          <Filters filtersData={filters} rangesData={ranges}/>
          <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
        </div>
      </div>
    );
}