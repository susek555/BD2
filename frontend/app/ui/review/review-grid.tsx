import { ReviewSearchParams } from '@/app/lib/definitions';
import {
  fetchReviewsByReviewee,
  fetchReviewsByReviewer,
} from '@/app/lib/data/reviews/data';
import { ReviewCard } from './review-card';

interface ReviewGridProps {
  variant: 'for' | 'by';
  userId: number;
  className?: string;
  searchParams: ReviewSearchParams;
}

export async function ReviewGrid({
  variant,
  userId,
  className = '',
  searchParams: serachParams,
}: ReviewGridProps) {
  let fetchFunc;

  if (variant === 'for') {
    fetchFunc = fetchReviewsByReviewee;
  } else {
    fetchFunc = fetchReviewsByReviewer;
  }

  const reviews = await fetchFunc(userId, serachParams);

  if (reviews.length === 0) {
    return (
      <div className='p-4 text-center text-sm text-gray-500'>
        No reviews found.
      </div>
    );
  }

  return (
    <div className={`grid grid-cols-1 gap-4 md:grid-cols-2 ${className}`}>
      {reviews.map((review) => (
        <ReviewCard key={review.id} review={review} variant={variant} />
      ))}
    </div>
  );
}
