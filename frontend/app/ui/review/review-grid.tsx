import { ReviewPage } from '@/app/lib/definitions/reviews';
import { ReviewCard } from './review-card';

interface ReviewGridProps {
  variant: 'for' | 'by';
  className?: string;
  reviewPage: ReviewPage;
}

export async function ReviewGrid({
  variant,
  className = '',
  reviewPage,
}: ReviewGridProps) {
  if (reviewPage.reviews.length === 0) {
    return (
      <div className='p-4 text-center text-sm text-gray-500'>
        No reviews found.
      </div>
    );
  }

  return (
    <div className={`grid grid-cols-1 gap-4 md:grid-cols-2 ${className}`}>
      {reviewPage.reviews.map((review) => (
        <ReviewCard key={review.id} review={review} variant={variant} />
      ))}
    </div>
  );
}
