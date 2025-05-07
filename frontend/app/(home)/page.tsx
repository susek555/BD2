import SideBar from "@/app/ui/(home)/sidebar";
import { Suspense } from "react";
import { OffersFoundSkeleton, OffersTableSkeleton } from "../ui/skeletons";
import Pagination from "../ui/(home)/pagination";
import OffersFoundInfo from "../ui/(home)/offers-found-info";
import OffersTable from "../ui/(home)/table";


export default async function Home() {
  const totalPages = 10;
  const totalOffers = 100;
  // TODO fetch available offers here and pass them to the table

  return (
  <main>
    <div className="flex flex-col md:flex-row flex-grow">
      <div className="w-full md:w-80 py-4 h-full flex-none">
        <Suspense>
          <SideBar />
        </Suspense>
      </div>
      <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
        <Suspense key={totalOffers} fallback={<OffersFoundSkeleton />}>
          <OffersFoundInfo totalOffers={totalOffers} />
        </Suspense>
        <div className="my-4" />
        <Suspense fallback={<OffersTableSkeleton />}>
          <OffersTable />
        </Suspense>
        <div className="mt-5 flex w-full justify-center pr-20">
          <Pagination totalPages={totalPages} />
        </div>
      </div>
    </div>
  </main>
  );
}
