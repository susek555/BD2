import { fetchSessionData } from '@/app/lib/data/account/data';
import {
  fetchReviewsByReviewee,
  fetchReviewsByReviewer,
} from '@/app/lib/data/reviews/data';
import { ReviewPage, ReviewSearchParams } from '@/app/lib/definitions/reviews';
import Pagination from '@/app/ui/(offers-table)/pagination';
import { ReviewFilterBox } from '@/app/ui/review/review-filters';
import { ReviewGrid } from '@/app/ui/review/review-grid';
import { ReviewGridSkeleton } from '@/app/ui/review/skeletons';
import { Suspense } from 'react';

export default async function ReviewsPage({
  searchParams,
}: {
  searchParams?: Promise<{
    variant?: 'by' | 'for';
    orderKey?: 'rating' | 'review_date';
    isOrderDesc?: string;
    ratings?: string;
    page?: string;
  }>;
}) {
  const { userId } = await fetchSessionData();

  const params = (await searchParams) || {};

  const variant = (params.variant === 'by' ? 'by' : 'for') as 'for' | 'by';

  const orderKey = params.orderKey || 'review_date';
  const isOrderDesc = params.isOrderDesc !== 'false';
  const page = Number(params.page || 1);

  const ratingsParam = params.ratings;
  const selectedRatings = ratingsParam
    ? ratingsParam.split(',').map(Number)
    : [];

  const reviewSearchParams: ReviewSearchParams = {
    is_order_desc: isOrderDesc,
    order_key: orderKey,
    pagination: {
      page: page,
      page_size: 20,
    },
    ratings: selectedRatings.length > 0 ? selectedRatings : undefined,
  };

  let fetchFunc: (
    id: number,
    searchParams: ReviewSearchParams,
  ) => Promise<ReviewPage>;

  if (variant === 'for') {
    fetchFunc = fetchReviewsByReviewee;
  } else {
    fetchFunc = fetchReviewsByReviewer;
  }

  const reviewPage = await fetchFunc(userId, reviewSearchParams);
  const totalPages = reviewPage.pagination.totalPages;

  return (
    <div className='flex flex-grow flex-row'>
      <div className='h-full w-80 flex-none py-4'>
        <ReviewFilterBox includeReviewToggle={true} />
      </div>

      <div className='flex flex-1 flex-col'>
        <div className='flex-1, p-6 px-12'>
          <Suspense fallback={<ReviewGridSkeleton />}>
            <ReviewGrid variant={variant} reviewPage={reviewPage} />
          </Suspense>
        </div>
        {totalPages === 0 && (
          <div className='mt-5 flex w-full justify-center pb-4'>
            <Suspense>
              <Pagination totalPages={totalPages} />
            </Suspense>
          </div>
        )}
      </div>
    </div>
  );
}
