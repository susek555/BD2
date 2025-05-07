import SideBar from "@/app/ui/(home)/sidebar";
import { Suspense } from "react";
import { OffersFoundSkeleton, OffersTableSkeleton } from "../ui/skeletons";


export default async function Home() {
  return (
  <main>
    <div className="flex flex-col md:flex-row flex-grow">
      <div className="w-full md:w-80 py-4 h-full flex-none">
        <Suspense>
          <SideBar />
        </Suspense>
      </div>
      <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
        <Suspense>
          <OffersFoundSkeleton />
        </Suspense>
        <div className="my-4" />
        <Suspense>
          <OffersTableSkeleton />
        </Suspense>
      </div>
    </div>
  </main>
  );
}
