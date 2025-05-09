import Link from 'next/link';
import Filters from '@/app/ui/(home)/filters';
import Sorting from './sorting';
import OfferType from './offer-type';


export default function SideBar() {
    return (
      <div className="flex h-full flex-col px-3 py-2 md:px-2 rounded-r-lg border-black border-[2px]">
        <div className="flex grow flex-row justify-between space-x-2 md:flex-col md:space-x-0 md:space-y-2">
          <OfferType  />
          <Sorting />
          <Filters />
          <div className="hidden h-auto w-full grow rounded-md bg-gray-50 md:block"></div>
        </div>
      </div>
    );
}