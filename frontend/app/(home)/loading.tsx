import { OffersTableSkeleton, OffersFoundSkeleton, SideBarSkeleton } from "../ui/skeletons";

export default function Loading() {
    return (
        <main>
          <div className="flex flex-col md:flex-row flex-grow">
            <div className="w-full md:w-80 py-4 h-full flex-none">
              <SideBarSkeleton />
            </div>
            <div className="flex-grow p-6 md:overflow-y-auto md:px-12 md:py-8">
              <OffersFoundSkeleton />
              <div className="my-4 " />
              <OffersTableSkeleton />
            </div>
          </div>
        </main>
      );
}