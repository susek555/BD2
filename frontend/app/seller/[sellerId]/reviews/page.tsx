import {
  fetchAverageRating,
  fetchRatingDistribution,
  fetchReviewsByReviewee,
} from '@/app/lib/data/reviews/data';
import { ReviewSearchParams } from '@/app/lib/definitions/reviews';
import Pagination from '@/app/ui/(offers-table)/pagination';
import { RatingDistribution } from '@/app/ui/review/rating-distribution';
import { ReviewFilterBox } from '@/app/ui/review/review-filters';
import { ReviewGrid } from '@/app/ui/review/review-grid';
import { ReviewGridSkeleton } from '@/app/ui/review/skeletons';
import { ChevronLeftIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import { Suspense } from 'react';

interface PageProps {
  params: Promise<{
    sellerId: string;
  }>;
  searchParams?: Promise<{
    orderKey?: 'rating' | 'review_date';
    isOrderDesc?: string;
    ratings?: string;
    page?: string;
    callbackUrl?: string;
  }>;
}

export default async function ReviewsPage({ params, searchParams }: PageProps) {
  const { sellerId } = await params;
  const serachParams = (await searchParams) || {};

  const callbackUrl = serachParams.callbackUrl || '/';
  const orderKey = serachParams.orderKey || 'review_date';
  const isOrderDesc = serachParams.isOrderDesc !== 'false';
  const page = Number(serachParams.page || 1);

  const ratingsParam = serachParams.ratings;
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

  const [sellerAverageRating, ratingDistribution, reviewPage] =
    await Promise.all([
      fetchAverageRating(Number(sellerId)),
      fetchRatingDistribution(Number(sellerId)),
      fetchReviewsByReviewee(Number(sellerId), reviewSearchParams),
    ]);

  return (
    <div className='my-3'>
      <div className='relative border-b bg-white shadow-sm'>
        <div className='mx-auto max-w-full px-6 py-6 lg:px-12'>
          <div className='flex flex-col space-y-2 sm:flex-row sm:items-center sm:justify-between sm:space-y-0'>
            <div className='flex items-center space-x-4'>
              <Link
                href={callbackUrl}
                className='flex items-center text-gray-500 transition-colors hover:text-gray-700'
              >
                <ChevronLeftIcon className='h-6 w-6 pr-1.5' />

                <span className='hidden sm:inline'>Back to offer</span>
                <span className='sm:hidden'>Back</span>
              </Link>
              <div className='border-l border-gray-300 pl-4'>
                <h1 className='text-xl font-semibold text-gray-900 sm:text-2xl'>
                  Reviews
                </h1>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className='mx-auto max-w-full px-6 py-8'>
        <div className='grid grid-cols-1 gap-4 lg:grid-cols-12'>
          <div className='lg:col-span-3'>
            <div className='lg:sticky lg:top-24'>
              <ReviewFilterBox includeReviewToggle={false} />
            </div>
          </div>

          <div className='flex flex-1 flex-col lg:col-span-6'>
            <div className='flex-1, p-6 px-12'>
              <Suspense fallback={<ReviewGridSkeleton />}>
                <ReviewGrid variant='for' reviewPage={reviewPage} />
              </Suspense>
            </div>

            {reviewPage.pagination.totalRecords !== 0 && (
              <div className='mt-5 flex w-full justify-center'>
                <Suspense>
                  <Pagination totalPages={reviewPage.pagination.totalPages} />
                </Suspense>
              </div>
            )}
          </div>

          <div className='lg:col-span-3'>
            <div className='lg:sticky lg:top-24'>
              <RatingDistribution
                averageRating={sellerAverageRating}
                distribution={ratingDistribution}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
