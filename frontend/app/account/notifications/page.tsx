import Pagination from "@/app/ui/(offers-table)/pagination";
import NotificationsFilterBox from "@/app/ui/account/notifications/notifications-filter-box";
import { Suspense } from "react";

export default function Page() {
  const totalPages = 10; // TODO get page count from API

  return (
      <div className='flex flex-grow flex-row'>
        <div className='h-full w-80 flex-none py-4'>
          <NotificationsFilterBox />
        </div>

        <div className='flex flex-1 flex-col'>
          <div className='flex-1, p-6 px-12'>
            {/* <Suspense fallback={<ReviewGridSkeleton />}> */}
              {/* <ReviewGrid
                variant={variant}
                userId={userId}
                searchParams={reviewSearchParams}
              /> */}
            {/* </Suspense> */}
          </div>

          <div className='mt-5 flex w-full justify-center pb-4'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        </div>
      </div>
    );
}