import Filters from '@/app/ui/(filters-sidebar)/filters';
import Sorting from '@/app/ui/(filters-sidebar)/sorting';
import OfferType from '@/app/ui/(filters-sidebar)/offer-type';
import { fetchFilterFields, fetchOfferTypes, fetchProducersAndModels, fetchSortingOptions, prepareRangeFields } from '@/app/lib/data/filters-sidebar/data';
import ProducersAndModels from '@/app/ui/(filters-sidebar)/producers-and-models';
import ApplyButton from '@/app/ui/(filters-sidebar)/apply-button';


export default async function SideBar() {
    const [offerTypes, filters, ranges, sortingOptions, producersAndModels] = await Promise.all([
      fetchOfferTypes(),
      fetchFilterFields(),
      prepareRangeFields(),
      fetchSortingOptions(),
      fetchProducersAndModels()
    ]);

    return (
      <div className="flex h-full flex-col px-3 py-2 md:px-2 rounded-r-lg border-black border-[2px]">
        <div className="flex grow flex-col justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
          <OfferType offerTypes={offerTypes}/>
          <Sorting sortingOptions={sortingOptions}/>
          <ProducersAndModels producersAndModels={producersAndModels}/>
          <Filters filtersData={filters} rangesData={ranges}/>
          <ApplyButton />
          <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
        </div>
      </div>
    );
}