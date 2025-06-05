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

  // Fetch all necessary data
  const [sellerAverageRating, ratingDistribution] = await Promise.all([
    fetchAverageRating(Number(sellerId)),
    fetchRatingDistribution(Number(sellerId)),
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
                <svg
                  className='mr-2 h-5 w-5'
                  fill='none'
                  stroke='currentColor'
                  viewBox='0 0 24 24'
                >
                  <path
                    strokeLinecap='round'
                    strokeLinejoin='round'
                    strokeWidth={2}
                    d='M15 19l-7-7 7-7'
                  />
                </svg>
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

          <div className='lg:col-span-6'>
            <Suspense fallback={<ReviewGridSkeleton />}>
              <ReviewGrid
                variant='for'
                userId={Number(sellerId)}
                searchParams={reviewSearchParams}
              />
            </Suspense>
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
